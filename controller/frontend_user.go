package controller

import (
	"ccrd/db"
	"ccrd/model"
	"net/http"

	"github.com/gin-gonic/gin"
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
