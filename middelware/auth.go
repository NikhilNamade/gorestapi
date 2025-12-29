package middelware

import (
	"net/http"

	"example.com/REST-API/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	authtoken := context.Request.Header.Get("Token")

	if authtoken == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Token is invalid1"})
		return
	}
	userId, err := utils.AuthenticateUser(authtoken)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "Token is invalid2"})
		return
	}
	context.Set("userId",userId)
	context.Next()
}
