package main

import (
	"fmt"
	"time"

	"example.com/REST-API/db"
	"example.com/REST-API/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found")
	}
	db.InitDB()
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // or specify "http://localhost:59489" for security
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	routes.RegisterRountes(server)

	port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "3000"
	// }

	err = server.Run(":" + port)
	if err != nil {
		fmt.Println("Err:", err)
	}
}
