package controller

import (
	"ccrd/server/khanscr"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Admin struct{}

type ItemsAllShow struct {
	ID string
}

func (a *Admin) UserGetAdmin(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "admin/index.html", gin.H{
		"title": "Age Of Khagan Thailand | Dashboard",
	})
}

// GetItems displays admin home page
func (a *Admin) GetItemsAll(ctx *gin.Context) {

	items := khanscr.GetAllItems()
	ctx.HTML(http.StatusOK, "admin/items/items.html", gin.H{
		"title": "Age Of Khagan Thailand | Dashboard",
		"items": items,
	})

}
