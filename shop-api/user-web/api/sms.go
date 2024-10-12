package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"shop-api/user-web/forms"
	"shop-api/user-web/global"
	"strings"
	"time"
)

func GenerateSmsCode(width int) string {
	//生成width长度的短信验证码
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendSms(ctx *gin.Context) {
	//client, err := dysmsapi.NewClientWithAccessKey(
	//	"cn-beijing", "", "")
	//if err != nil {
	//	panic(err)
	//}
	//smsCode := GenerateSmsCode(6)
	//request := requests.NewCommonRequest()
	//request.Method = "POST"
	//request.Scheme = "https" // https | http
	//request.Domain = "dysmsapi.aliyuncs.com"
	//request.Version = "2017-05-25"
	//request.ApiName = "SendSms"
	//request.QueryParams["RegionId"] = "cn-beijing"
	//request.QueryParams["PhoneNumbers"] = "18888887777"                 //手机号
	//request.QueryParams["SignName"] = "慕学在线"                        //阿里云验证过的项目名 自己设置
	//request.QueryParams["TemplateCode"] = "SMS_181850725"               //阿里云的短信模板号 自己设置
	//request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	//response, err := client.ProcessCommonRequest(request)
	//fmt.Print(client.DoAction(request, response))
	//if err != nil {
	//	fmt.Print(err.Error())
	//}

	// 模拟发送短信，直接将生成的验证码存入Redis
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	smsCode := GenerateSmsCode(6)
	redisKey := fmt.Sprintf("sms_code_%s_%d", sendSmsForm.Mobile, sendSmsForm.Type)
	global.RDB.Set(context.Background(), redisKey, smsCode, time.Duration(300)*time.Second)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
