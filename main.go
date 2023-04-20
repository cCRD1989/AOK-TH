package main

import (
	"ccrd/db"
	"ccrd/middleware"
	"ccrd/server/khanscr"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zalando/gin-oauth2/google"
)

var redirectURL, credFile string

// init web google api
func init() {
	bin := path.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s
================
`, bin)
		flag.PrintDefaults()
	}

	// auth https://console.cloud.google.com/
	// http://127.0.0.1:80/auth/google/
	// https://ageofkhaganth.com/auth/google/
	flag.StringVar(&redirectURL, "redirect", "https://ageofkhaganth.com/auth/google/", "URL to be redirected to after authorization.")
	flag.StringVar(&credFile, "cred-file", "./test-clientid.google.json", "Credential JSON file")

}

func main() {
	khanscr.Init()

	// load env
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

	// google
	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}

	secret := []byte("secret")
	sessionName := "GOCSPX"

	// init settings for google auth
	google.Setup(redirectURL, credFile, scopes, secret)
	r.Use(google.Session(sessionName))

	r.Use(middleware.UserSession())
	serveRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	r.Run(":" + port)
}
