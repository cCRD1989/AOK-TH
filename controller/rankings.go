package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Rankings struct{}

func (r Rankings) Ranking(ctx *gin.Context) {

	class := ctx.Param("class")

	fmt.Println("class: ", class)

	ctx.HTML(http.StatusOK, "frontend/rankings.html", gin.H{
		"title": "Age Of Khagan | Rankings.",
	})

}
