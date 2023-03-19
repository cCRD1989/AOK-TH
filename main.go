package main

import (
	"ccrd/db"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
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

	// Create and Connect DB
	db.ConnectDB()
	db.Migrate()

	// set cors all port
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	//create Folder
	os.MkdirAll("public/download", 0755)

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

func checkAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func createViews() multitemplate.Render {
	var fn = template.FuncMap{
		"getPlayersOnlineCount": func() string {
			return "199 คน"
		},
	}
	var r = multitemplate.New()
	var vtpath = filepath.Join("views", "templates")
	var dirs, err = ioutil.ReadDir("views/layouts/")
	checkAndPanic(err)
	for _, dir := range dirs {
		var dirName = dir.Name()
		layouts, err := filepath.Glob(fmt.Sprintf("views/layouts/%s/*.html", dirName))
		checkAndPanic(err)

		var templates = []string{}
		err = filepath.Walk(fmt.Sprintf("views/templates/%s/", dirName), func(path string, f os.FileInfo, err error) error {
			if filepath.Ext(path) == ".html" {
				templates = append(templates, path)
			}
			return nil
		})
		checkAndPanic(err)
		for _, tmpl := range templates {
			var tname = strings.Replace(tmpl, vtpath, "", 1)  // ลบพาทออก
			tname = strings.Replace(tname, "\\", "/", -1)[1:] //เปลี่ยนให้เป็นรูท
			log.Printf("[GIN-debug] %-6s %-25s --> %s\n", "VIEW", dirName, tname)
			r.AddFromFilesFuncs(tname, fn, append(layouts, tmpl)...)
			//r.AddFromFiles(tname, append(layouts, tmpl)...)
		}
	}
	return r
}
