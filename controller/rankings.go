package controller

import (
	"ccrd/model/aokmodel"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Rankings struct{}

func (r Rankings) Ranking(ctx *gin.Context) {
	var usr, _ = ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	class := ctx.Param("class")

	fmt.Println("class: ", class)

	ctx.HTML(http.StatusOK, "frontend/rankings.html", gin.H{
		"title": "Age Of Khagan | Rankings.",
		"user":  user,
	})

}
