package main

import (
	"ccrd/db"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("APP_ENV") != "ReleaseMode" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error Loading File")
		}
		fmt.Println("Env Loading File")

	} else {
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("SetMode = ReleaseMode")
	}

	db.ConnectDB()
	db.Migrate()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	os.MkdirAll("uploads/products", 0755)

	r := gin.Default()

	r.Use(cors.New(corsConfig))
	r.Static("/uploads", "./uploads")

	serveRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	r.Run(":" + port)
}
