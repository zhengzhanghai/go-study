package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop-api/user-web/forms"
	"shop-api/user-web/global"
	"shop-api/user-web/global/response"
	"shop-api/user-web/middlewares"
	"shop-api/user-web/models"
	"shop-api/user-web/proto"
	"strconv"
	"strings"
	"time"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			}
		}
	}
}

func GenUserLoginToken(userInfo *proto.UserInfoResponse) (string, error) {
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(userInfo.Id),
		Nickname:    userInfo.Nickname,
		AuthorityId: uint(userInfo.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, // 30天过期
			Issuer:    "zzh",
		},
	}
	token, err := j.CreateToken(claims)
	return token, err
}

func HandleValidatorError(ctx *gin.Context, err error) {
	errors, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errors.Translate(global.Trans)),
	})
}

func GetUserList(ctx *gin.Context) {
	zap.S().Debug("获取用户列表")
	pn := ctx.DefaultQuery("pn", "1")
	pnInt, _ := strconv.Atoi(pn)
	psize := ctx.DefaultQuery("psize", "10")
	psizeInt, _ := strconv.Atoi(psize)

	resp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Page:     uint32(pnInt),
		PageSize: uint32(psizeInt),
	})
	if err != nil {
		zap.S().Error("[GetUserList] 查询 【用户列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range resp.Users {
		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.Nickname,
			Birthday: time.Time(time.Unix(int64(value.Birthday), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for filed, err := range fileds {
		rsp[filed[strings.Index(filed, ".")+1:]] = err
	}
	return rsp
}

func LoginByPassword(ctx *gin.Context) {
	zap.S().Debug("密码登录")

	passwordForm := forms.PasswordLoginForm{}
	if err := ctx.ShouldBind(&passwordForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	if rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}
			return
		}
	} else {
		// 通过用户查询到了用户，并没有检查密码
		if passRsp, passErr := global.UserSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordForm.Password,
			EncryptedPassword: rsp.Password,
		}); passErr != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{
				"msg": "登录错误",
			})
		} else {
			if passRsp.Success {
				token, err := GenUserLoginToken(rsp)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成Token失败",
					})
					return
				}
				ctx.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.Nickname,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "登录失败",
				})
			}
		}
	}
}

func Register(ctx *gin.Context) {
	form := forms.RegisterForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	smsCodeRedisCodeKey := fmt.Sprintf("sms_code_%s_%d", form.Mobile, 1)
	smsCode, err := global.RDB.Get(context.Background(), smsCodeRedisCodeKey).Result()
	if errors.Is(err, redis.Nil) || smsCode != form.Code {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Mobile:   form.Mobile,
		Password: form.Password,
		Nickname: form.Mobile,
	})

	if err != nil {
		zap.S().Errorf("[Register] 查新新建用户失败: %s", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 注册成功后，登录生成token
	token, err := GenUserLoginToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.Nickname,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}
