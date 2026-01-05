package routes

import (
	"example.com/REST-API/middelware"
	"github.com/gin-gonic/gin"
)

func RegisterRountes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.Use(middelware.Authenticate)
	authenticated.GET("/user/:id",getUserById)
	authenticated.GET("/events", getEvents)
	authenticated.GET("/eventsbyuserId/:id", getEventsUserId)
	authenticated.POST("/events", createEvents)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)
	authenticated.POST("/register/:id", registerEvent)
	authenticated.DELETE("/register/:id", registerDELETE)
	authenticated.GET("/user",getUserByToken)
	authenticated.GET("/eventbyuser",getEventsByUserId)
	authenticated.GET("/registerUserId",getRegisterEventByUserId)
	authenticated.GET("/detailEvents",getDetiailEvents)
	authenticated.POST("/addStory",addNewStory);
	authenticated.GET("/allStory",getAllStory);
	authenticated.GET("/userStory",getUserStory)
	authenticated.PUT("/updateviewStory/:id",updateStoryView)
	authenticated.POST("/follow/:id",followToConnections)
	authenticated.DELETE("/unfollow/:id",disConnect)
	authenticated.GET("/getConnections",getAllConnectionByUser)
	authenticated.GET("/getAllUsers", getallUser)
	server.POST("/getUserByEmail",getUserByEmail)
	server.POST("/resetPassword",resetPassword)
	server.GET("/events/:id", getEventsById)
	server.GET("/register", registerGET)
	server.POST("/signup", signupUser)
	server.POST("/login", loginUser)

	
}
