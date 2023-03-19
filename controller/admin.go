package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Admin struct{}

func (a *Admin) UserGetAdmin(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "admin/index.html", gin.H{
		"title": "Age Of Khagan Thailand | Dashboard",
	})
}
