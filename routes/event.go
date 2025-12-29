package routes

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"example.com/REST-API/models"
	"example.com/REST-API/utils"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	token := context.Request.Header.Get("Token")
	if token == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "User is not authenticated"})
		return
	}
	userId, err := utils.AuthenticateUser(token)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "User is not authenticated"})
		return
	}
	result, err := models.GetAllEvents(userId)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Something wrong happened at server while fecthing user", "err": err})
		return
	}
	context.JSON(http.StatusOK, result) // here whatever the data is it will transform to the json by gin
}

func getEventsById(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Id not found"})
		return
	}

	result, err := models.GetIDEvent(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "Id not found"})
		return
	}
	context.JSON(http.StatusOK, result)
}

func uploadToCloudinary(file any, context *gin.Context, fileHeader *multipart.FileHeader) (string, error) {
	cld, err := cloudinary.NewFromParams("dusrnmnvs", "514497943162127", "qxpBZ0JXdzrt-kyYv5yxOoc59SE")
	if err != nil {
		fmt.Println(err, "4")
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Something wrong happened at server"})
		return "", err
	}

	uploadResult, err := cld.Upload.Upload(context, file, uploader.UploadParams{
		Folder:   "Events",
		PublicID: fileHeader.Filename,
	})
	return uploadResult.SecureURL, err
}

func createEvents(context *gin.Context) {
	var events models.Event
	err := context.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println(err, "2")
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Something wrong happened at server"})
		return
	}

	layout := time.RFC3339

	fmt.Println(context.Request.FormValue("date_time"))

	dateTime, err := time.Parse(layout, context.Request.FormValue("date_time"))

	if err != nil {
		fmt.Println("Bind Error:", err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Failed to bind form-data"})
		return
	}
	fmt.Println(dateTime)
	events.Name = context.PostForm("name")
	events.Description = context.PostForm("description")
	events.Location = context.PostForm("location")
	events.Category = context.PostForm("Category")
	events.Datetime = dateTime
	feeStr := context.PostForm("Fees")
	fees, err := strconv.Atoi(feeStr)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Something wrong happened at server"})
		return
	}
	events.Fees = fees

	file, fileHeader, err := context.Request.FormFile("event_image")
	if err != nil {
		fmt.Println(err, "3")
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Something wrong happened at server"})
		return
	}

	defer file.Close()

	uploadResult, err := uploadToCloudinary(file, context, fileHeader)

	if err != nil {
		fmt.Println(err, "5")
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Something wrong happened at server"})
		return
	}
	userId := context.GetInt64("userId")
	fmt.Println(userId)
	events.UserId = userId
	events.Profile = uploadResult
	err = events.Save()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "unable to store data", "Error": err})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"msg": "Event is created succesfully", "event": events})
}

func updateEvent(context *gin.Context) {

	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Id not found"})
		return
	}

	// var event models.Event
	event, err := models.GetIDEvent(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Event not found"})
		return
	}
	// authenticate the user login
	userId := context.GetInt64("userId")
	if userId != event.UserId {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "user not valid "})
		return
	}

	err = context.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println(err, "2")
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Something wrong happened at server"})
		return
	}

	name := context.PostForm("name")
	description := context.PostForm("description")
	location := context.PostForm("location")
	category := context.PostForm("Category")
	feeStr := context.PostForm("Fees")
	fees, err := strconv.Atoi(feeStr)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Something wrong happened at server"})
		return
	}
	Fees := fees
	file, header, err := context.Request.FormFile("profile_image")
	var profileUrl string
	if err == nil && header != nil {
		uploadResult, err := uploadToCloudinary(file, context, header)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"msg": "Unable to upload file"})
			return
		}
		profileUrl = uploadResult
	} else {
		profileUrl = context.Request.FormValue("profile_image")
	}
	dateTime, err := time.Parse(time.RFC3339, context.Request.FormValue("date_time"))
	updateEvent := models.Event{
		ID:          id,
		Name:        name,
		Description: description,
		Location:    location,
		Datetime:    dateTime,
		UserId:      userId,
		Profile:     profileUrl,
		Category:    category,
		Fees:        Fees,
	}
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Data not found"})
		return
	}

	err = updateEvent.UpdateEvent()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Update not done"})
		return
	}
	context.JSON(http.StatusOK, "")
}

func deleteEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "Id not found"})
		return
	}

	event, err := models.GetIDEvent(id)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "Event not found"})
		return
	}
	// authenticate the user login
	userId := context.GetInt64("userId")
	if userId != event.UserId {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "user not valid "})
		return
	}
	err = event.DeleteEvent()
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "Event not deleted"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"msg": "Event Deleted"})
}

func getEventsByUserId(context *gin.Context) {
	userId := context.GetInt64("userId")
	if userId == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "UserId not found"})
		return
	}
	result, err := models.GetEventsByUser(userId)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "UserId not found"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"events": result})
}

func getEventsUserId(context *gin.Context) {
	userId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}

	result, err := models.GetEventsByUser(userId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"msg": err});
		return;
	}
	context.JSON(http.StatusOK, gin.H{"events":result});
}
