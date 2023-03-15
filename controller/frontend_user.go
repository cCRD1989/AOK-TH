package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Frontend struct{}

func (u *Frontend) UserGetHome(c *gin.Context) {

	c.HTML(http.StatusOK, "frontend/index.html", gin.H{
		"title": "Age Of Khagan Thailand",
	})
}
