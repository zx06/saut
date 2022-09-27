package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Auth) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo:
		var user = "xuzhuo"
		token, err := a.GrantToken(user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "验证失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"access_token": token,
		})

	}
}
