package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"example.com/REST-API/models"
	"github.com/gin-gonic/gin"
)

func addNewStory(context *gin.Context) {
	var story models.Story

	err := context.Request.ParseMultipartForm(10 >> 20)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Something wrong happened at server"})
		return
	}
	story.UserID = int(context.GetInt64("userId"))

	file, fileHeader, err := context.Request.FormFile("story-file")

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Something wrong happened at server"})
		return
	}

	defer file.Close()

	uploadResult, err := uploadToCloudinary(file, context, fileHeader)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Something wrong happened at server"})
		return
	}
	story.File = uploadResult

	err = story.Save()
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Something wrong happened at server"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"msg": "Story is created succesfully", "Story": story})
}

func getAllStory(context *gin.Context) {
	userId := context.GetInt64("userId")
	if userId == 0 {
		context.JSON(http.StatusBadGateway, gin.H{"msg": "User is not authenticated"})
		return
	}

	result, err := models.GetAllStory()
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Internal server error"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"result": result})
}

func getUserStory(context *gin.Context) {
	userId := context.GetInt64("userId")
	if userId == 0 {
		context.JSON(http.StatusBadGateway, gin.H{"msg": "User is not authenticated"})
		return
	}
	result, err := models.GetuserStory(int(userId))
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Internal server error"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"result": result})
}

func updateStoryView(context *gin.Context) {
	userId := context.GetInt64("userId")
	storyId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if userId == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "User is not authenticated"})
		return
	}

	if storyId == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "User is not authenticated"})
		return
	}
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": "User is not authenticated"})
		return
	}
	err = models.UpdateStoryView(int(storyId),time.Now());
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Something wrong happened at server"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"msg": "Updated Succesfully"})
}
