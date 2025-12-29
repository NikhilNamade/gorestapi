package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/REST-API/models"
	"github.com/gin-gonic/gin"
)

func signupUser(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		fmt.Println(err, "1")
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Data not found"})
		return
	}

	err = user.Save()

	if err != nil {
		fmt.Println(err, "2")
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Data not found"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"msg": "User signup"})
}

func getallUser(context *gin.Context) {
	users, err := models.GetAllUsers()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Users not found"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"users": users})
}

func loginUser(context *gin.Context) {
	var user models.User

	context.ShouldBindJSON(&user)

	token, err := models.LoginUser(user)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "Users not valid"})
		return
	}

	context.JSON(http.StatusAccepted, gin.H{"msg": "Users valid", "token": token})
}

func getUserByToken(context *gin.Context) {
	token := context.Request.Header.Get("Token")
	if token == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "Authenticate User!"})
		return
	}

	user, err := models.GetUserByToken(token)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": "User not found"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"user": user})
}

func getUserById(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}

	user, err := models.GetUserById(id)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	context.JSON(http.StatusOK,gin.H{"User":user});
}
