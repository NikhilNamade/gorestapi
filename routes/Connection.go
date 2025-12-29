package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/REST-API/models"
	"github.com/gin-gonic/gin"
)

func followToConnections(context *gin.Context) {
	var connect models.Connection
	idStr := context.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	userId := context.GetInt64("userId") // already int; no need to convert
	fmt.Println(id)
	fmt.Println(userId)
	connect.FollowBy = int(userId)
	connect.FollowTo = id
	if connect.FollowBy == connect.FollowTo{
		context.JSON(http.StatusBadRequest, gin.H{"msg":"Connection not allowed"})
		return
	}
	err = connect.Save()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	// correct message
	msg := fmt.Sprintf("%d followed %d", userId, id)

	context.JSON(http.StatusOK, gin.H{"msg": msg})
}

func getAllConnectionByUser(context *gin.Context) {
	userId := context.GetInt64("userId")
	result, err := models.GetAllConnectionByUser(int(userId))
	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"result": result})
}

func disConnect(context *gin.Context){
	idStr := context.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	userId := context.GetInt64("userId") // already int; no need to convert
	err = models.DisConnect(id,int(userId));
	if err != nil || id == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	// correct message
	msg := fmt.Sprintf("%d unfollowed %d", userId, id)

	context.JSON(http.StatusOK, gin.H{"msg": msg})
}