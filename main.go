package main

import (
	"ccrd/db"
	"ccrd/server/khanscr"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	khanscr.Init()

	if os.Getenv("APP_ENV") != "ReleaseMode" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error Loading File")
		}
		fmt.Println("Env Loading File")

	} else {
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("SetMode = ReleaseMode")
	}

	// Create and Connect DB
	db.ConnectDB()
	db.Migrate()

	// set cors all port
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	//create Folder
	os.MkdirAll("/public/download", 0755)

	r := gin.Default()
	r.HTMLRender = createViews()
	r.Static("/public", "./public")
	r.Static("/download/public", "./public")

	r.StaticFile("/favicon.ico", "./public/favicon.ico")

	r.Use(cors.New(corsConfig))

	serveRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	r.Run(":" + port)
}
