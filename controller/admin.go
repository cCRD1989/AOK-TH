package controller

import (
	"ccrd/db"
	"ccrd/model"
	"ccrd/server/khanscr"
	"net/http"
	"strconv"
	"strings"

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

// Log Topup
func (a *Admin) Logtopup(ctx *gin.Context) {

	selectN := ctx.Query("select")
	inputtxt := ctx.Query("inputtxt")
	from := ctx.Query("from")
	to := ctx.Query("to")

	//ค้นหาทั้งหมด
	if selectN == "" && inputtxt == "" && from == "" && to == "" {

		logall := []model.LogTopup{}
		db.Conn.Where("data_type", "NotificationTopup").Find(&logall)

		PriceAll := 0
		for _, v := range logall {
			al := strings.Replace(v.Price, "THB", "", 1)
			sti, _ := strconv.Atoi(al)
			PriceAll += sti
		}

		ctx.HTML(http.StatusOK, "admin/logtopup/logtopup.html", gin.H{
			"title":    "Age Of Khagan Thailand | Log Topup",
			"logall":   logall,
			"priceall": PriceAll,
		})
		return
	}

	// Username
	if selectN == "Username" {
		logall := []model.LogTopup{}
		db.Conn.Where("data_type", "NotificationTopup").Where("user_id = ?", inputtxt).Find(&logall)

		PriceAll := 0
		for _, v := range logall {
			al := strings.Replace(v.Price, "THB", "", 1)
			sti, _ := strconv.Atoi(al)
			PriceAll += sti
		}

		ctx.HTML(http.StatusOK, "admin/logtopup/logtopup.html", gin.H{
			"title":    "Age Of Khagan Thailand | Log Topup",
			"logall":   logall,
			"priceall": PriceAll,
		})
		return
	}

	// Channel
	if selectN == "Channel" {
		logall := []model.LogTopup{}
		db.Conn.Where("data_type", "NotificationTopup").Where("channel = ?", inputtxt).Find(&logall)

		PriceAll := 0
		for _, v := range logall {
			al := strings.Replace(v.Price, "THB", "", 1)
			sti, _ := strconv.Atoi(al)
			PriceAll += sti
		}

		ctx.HTML(http.StatusOK, "admin/logtopup/logtopup.html", gin.H{
			"title":    "Age Of Khagan Thailand | Log Topup",
			"logall":   logall,
			"priceall": PriceAll,
		})
		return
	}

	// Select Date
	if from != "" && to != "" {
		logall := []model.LogTopup{}
		db.Conn.Where("data_type", "NotificationTopup").Where("created_at BETWEEN ? AND ?", from, to).Find(&logall)

		PriceAll := 0
		for _, v := range logall {
			al := strings.Replace(v.Price, "THB", "", 1)
			sti, _ := strconv.Atoi(al)
			PriceAll += sti
		}

		ctx.HTML(http.StatusOK, "admin/logtopup/logtopup.html", gin.H{
			"title":    "Age Of Khagan Thailand | Log Topup",
			"logall":   logall,
			"priceall": PriceAll,
		})
		return
	}

	//

	// ไม่เข้าพวก
	ctx.HTML(http.StatusOK, "admin/logtopup/logtopup.html", gin.H{
		"title": "Age Of Khagan Thailand | Log Topup",
	})
}
