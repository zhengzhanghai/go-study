package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-api/user-web/models"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		curUser := claims.(*models.CustomClaims)
		if curUser.AuthorityId != 2 {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
