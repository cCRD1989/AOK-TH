package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Guides struct{}

func (g *Guides) WroclawFortress(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "frontend/guides.html", gin.H{
		"title": "Age Of Khagan | Wroclaw Fortress.",
	})

}
