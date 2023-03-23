package controller

import (
	"ccrd/db"
	"ccrd/model/aokmodel"

	"ccrd/server/khanscr"
	"fmt"
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
	test()
	ctx.HTML(http.StatusOK, "admin/items/items.html", gin.H{
		"title": "Age Of Khagan Thailand | Dashboard",
		"items": items,
	})

}

func test() {

	var user_id []aokmodel.Userlogin
	db.AOK_DB.Find(&user_id)
	fmt.Println(user_id)
}
