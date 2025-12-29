package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/REST-API/models"
	"example.com/REST-API/utils"
	"github.com/gin-gonic/gin"
)

func registerEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "EventId is not found"})
		return
	}
	userId := context.GetInt64("userId")
	var res models.Register
	res.EventId = id
	res.UserId = userId
	err = res.Save()
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "Not able to store registration","err":err})
		return
	}
	context.JSON(http.StatusOK,gin.H{"msg":"User Register For Event"})
}

func registerGET(context *gin.Context) {
	resi, err := models.Getall()
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Evenst not found"})
		return
	}
	context.JSON(http.StatusAccepted, gin.H{"regi": resi})
}

func registerDELETE(context *gin.Context) {
	token := context.Request.Header.Get("Token")
	if token == ""{
		context.JSON(http.StatusUnauthorized,gin.H{"msg":"User is not authorized"})
		return
	}

	userId,err := utils.AuthenticateUser(token)

	if err != nil{
		fmt.Println(err)
		context.JSON(http.StatusUnauthorized,gin.H{"msg":"User is not authorized"})
		return
	}
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "EventId is not found"})
		return
	}

	err = models.Deleteresi(id,userId)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "Event is not delete"})
		return
	}
	context.JSON(http.StatusOK,"")
}

func getRegisterEventByUserId(context *gin.Context){
	token := context.Request.Header.Get("Token")

	if token == ""{
		context.JSON(http.StatusUnauthorized,gin.H{"msg":"User is not authorized"})
		return
	}

	userId,err := utils.AuthenticateUser(token)

	if err != nil{
		fmt.Println(err)
		context.JSON(http.StatusUnauthorized,gin.H{"msg":"User is not authorized"})
		return
	}
	result,err := models.GetAllRegisterByUserId(userId)

	if err != nil{
		fmt.Println(err)
		context.JSON(http.StatusUnauthorized,gin.H{"msg":"Events not found"})
		return
	}
	context.JSON(http.StatusOK,result)
}

func getDetiailEvents(context *gin.Context){
	token := context.Request.Header.Get("Token")

	if token == ""{
		context.JSON(http.StatusUnauthorized,gin.H{"msg":"User is not authorized"})
		return
	}

	userId,err := utils.AuthenticateUser(token)

	if err != nil{
		fmt.Println(err)
		context.JSON(http.StatusUnauthorized,gin.H{"msg":"User is not authorized"})
		return
	}
	result,err := models.GetDetailEvents(userId)

	if err != nil{
		fmt.Println(err)
		context.JSON(http.StatusUnauthorized,gin.H{"msg":"Events not found"})
		return
	}
	context.JSON(http.StatusOK,result)
}