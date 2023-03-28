package controller

import (
	"ccrd/db"
	"ccrd/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zalando/gin-oauth2/google"
)

type Frontend struct{}

func (f *Frontend) UserGetHome(ctx *gin.Context) {

	visit := model.LogWeb{
		DataType:  "visit",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	ctx.HTML(http.StatusOK, "frontend/index.html", gin.H{
		"title": "Age Of Khagan Thailand",
	})
}

func (f *Frontend) UserGetDownload(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "Window" {
		visit := model.LogWeb{
			DataType:  "download_window",
			IPAddress: ctx.ClientIP(),
		}
		db.Conn.Save(&visit)
		ctx.Redirect(http.StatusTemporaryRedirect, "public/download/window.rar")
	} else if id == "Android" {
		visit := model.LogWeb{
			DataType:  "download_android",
			IPAddress: ctx.ClientIP(),
		}
		db.Conn.Save(&visit)
		ctx.Redirect(http.StatusTemporaryRedirect, "public/download/android.rar")
	} else if id == "App" {
		visit := model.LogWeb{
			DataType:  "download_App",
			IPAddress: ctx.ClientIP(),
		}
		db.Conn.Save(&visit)
		ctx.Redirect(http.StatusTemporaryRedirect, "public/download/app.rar")
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func (f *Frontend) Auth_google(ctx *gin.Context) {

	name := ctx.MustGet("user").(google.User)
	//ctx.JSON(http.StatusOK, gin.H{"Hello": "from private", "user": name})

	ctx.HTML(http.StatusOK, "frontend/auth.html", gin.H{
		"title":  "Age Of Khagan Thailand | Account",
		"email":  name.Email,
		"name":   name.Name,
		"imgsrc": name.Picture,
	})
}

func (f *Frontend) Auth_google_Regis(ctx *gin.Context) {
	email := ctx.DefaultQuery("email", "-")
	pass := ctx.DefaultQuery("password", "-")
	rpass := ctx.DefaultQuery("repassword", "-")

	fmt.Println("email", email)
	fmt.Println("pass", pass)
	fmt.Println("rpass", rpass)

}
