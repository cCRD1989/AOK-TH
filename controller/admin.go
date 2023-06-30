package controller

import (
	"ccrd/db"
	"ccrd/model"
	"ccrd/model/aokmodel"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Admin struct{}

func (a *Admin) UserGetHome(ctx *gin.Context) {

	visit := model.LogWeb{
		DataType:  "Admin",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	ctx.HTML(http.StatusOK, "admin/index.html", gin.H{
		"title": "Age Of Khagan Thailand | ADMIN",
		"user":  user,
	})
}

// POST สร้างโพส
func (f *Admin) CreateNew(ctx *gin.Context) {

	form := model.LogNewsRequest{}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error1": err.Error()})
		return
	}

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	if user.Userlevel != 1 {
		fmt.Println("err. Userlevel")
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	imagePath := "/public/data/img/img_news/" + uuid.NewString() + "." + strings.Split(form.Image.Filename, ".")[1]

	ctx.SaveUploadedFile(&form.Image, imagePath)

	addNew := model.LogNews{
		Datatype:     form.Datatype,
		Author:       form.Author,
		Subject:      form.Subject,
		Data:         form.Data,
		Image:        imagePath,
		Externallink: form.Externallink,
	}
	if err := db.Conn.Create(&addNew).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error3": err.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/newall")
}

// // GetItems displays admin home page
// func (a *Admin) GetItemsAll(ctx *gin.Context) {

// 	items := khanscr.GetAllItems()
// 	ctx.HTML(http.StatusOK, "admin/items/items.html", gin.H{
// 		"title": "Age Of Khagan Thailand | Dashboard",
// 		"items": items,
// 	})

// }

// // Log Topup
// func (a *Admin) Logtopup(ctx *gin.Context) {

// 	selectN := ctx.Query("select")
// 	inputtxt := ctx.Query("inputtxt")
// 	from := ctx.Query("from")
// 	to := ctx.Query("to")

// 	//ค้นหาทั้งหมด
// 	if selectN == "" && inputtxt == "" && from == "" && to == "" {

// 		logall := []model.LogTopup{}
// 		db.Conn.Where("data_type", "NotificationTopup").Find(&logall)

// 		PriceAll := 0
// 		for _, v := range logall {
// 			al := strings.Replace(v.Price, "THB", "", 1)
// 			sti, _ := strconv.Atoi(al)
// 			PriceAll += sti
// 		}

// 		ctx.HTML(http.StatusOK, "admin/logtopup/logtopup.html", gin.H{
// 			"title":    "Age Of Khagan Thailand | Log Topup",
// 			"logall":   logall,
// 			"priceall": PriceAll,
// 		})
// 		return
// 	}

// 	// Username
// 	if selectN == "Username" {
// 		logall := []model.LogTopup{}
// 		db.Conn.Where("data_type", "NotificationTopup").Where("user_id = ?", inputtxt).Find(&logall)

// 		PriceAll := 0
// 		for _, v := range logall {
// 			al := strings.Replace(v.Price, "THB", "", 1)
// 			sti, _ := strconv.Atoi(al)
// 			PriceAll += sti
// 		}

// 		ctx.HTML(http.StatusOK, "admin/logtopup/logtopup.html", gin.H{
// 			"title":    "Age Of Khagan Thailand | Log Topup",
// 			"logall":   logall,
// 			"priceall": PriceAll,
// 		})
// 		return
// 	}

// 	// Channel
// 	if selectN == "Channel" {
// 		logall := []model.LogTopup{}
// 		db.Conn.Where("data_type", "NotificationTopup").Where("channel = ?", inputtxt).Find(&logall)

// 		PriceAll := 0
// 		for _, v := range logall {
// 			al := strings.Replace(v.Price, "THB", "", 1)
// 			sti, _ := strconv.Atoi(al)
// 			PriceAll += sti
// 		}

// 		ctx.HTML(http.StatusOK, "admin/logtopup/logtopup.html", gin.H{
// 			"title":    "Age Of Khagan Thailand | Log Topup",
// 			"logall":   logall,
// 			"priceall": PriceAll,
// 		})
// 		return
// 	}

// 	// Select Date
// 	if from != "" && to != "" {
// 		logall := []model.LogTopup{}
// 		db.Conn.Where("data_type", "NotificationTopup").Where("created_at BETWEEN ? AND ?", from, to).Find(&logall)

// 		PriceAll := 0
// 		for _, v := range logall {
// 			al := strings.Replace(v.Price, "THB", "", 1)
// 			sti, _ := strconv.Atoi(al)
// 			PriceAll += sti
// 		}

// 		ctx.HTML(http.StatusOK, "admin/logtopup/logtopup.html", gin.H{
// 			"title":    "Age Of Khagan Thailand | Log Topup",
// 			"logall":   logall,
// 			"priceall": PriceAll,
// 		})
// 		return
// 	}

// 	//

// 	// ไม่เข้าพวก
// 	ctx.HTML(http.StatusOK, "admin/logtopup/logtopup.html", gin.H{
// 		"title": "Age Of Khagan Thailand | Log Topup",
// 	})
// }
