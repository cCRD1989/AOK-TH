package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Frontend struct{}

func (u *Frontend) UserGetHome(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "frontend/index.html", gin.H{
		"title": "Age Of Khagan Thailand",
	})
}
