package controller

import (
	"ccrd/db"
	"ccrd/model"
	"ccrd/model/aokmodel"
	"ccrd/unit"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ItemCode struct{}

func (i *ItemCode) PItemcode(ctx *gin.Context) {
	visit := model.LogWeb{
		DataType:  "PItemcode",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	//
	ctx.HTML(http.StatusOK, "frontend/itemcode.html", gin.H{
		"title": "Age Of Khagan Thailand | ItemCode",
		"user":  user,
		"bg":    "/public/data/img/LOGIN-BG.png",
	})
}

func (i *ItemCode) Itemcode(ctx *gin.Context) {

	idcode := ctx.DefaultPostForm("idcode", "")
	visit := model.LogWeb{
		DataType:  "Itemcode",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	errs := unit.Validate(map[string]interface{}{
		"idcode": idcode,
	}, map[string]string{

		"idcode": "required|min:32|max:32|alphanum",
	})
	if errs != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": "error", "msg": "ไม่สามารถทำรายการได้ ข้อมูลไม่ถูกต้อง"})
		return
	}

	tdata := model.LogTokenregister{}
	if err := db.Conn.Where("tokenid = ?", idcode).First(&tdata).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": "error", "msg": "ไม่สามารถทำรายการได้ ItemCode ไม่ถูกต้อง"})
		return
	}
	if tdata.Status != 0 {
		ctx.JSON(http.StatusOK, gin.H{"status": "error", "msg": "ไม่สามารถทำรายการได้ ItemCode ถูกเปิดใช้งานไปแล้ว !"})
		return
	}

	username := aokmodel.Userlogin{}
	if err := db.AOK_DB.Where("Username = ?", tdata.Username).First(&username).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": "error", "msg": "ไม่สามารถทำรายการได้ ไม่มีบัญชีในระบบ"})
		return
	}

	sendMail := model.LogMailTopup{
		Eventid:    "10",
		Senderid:   username.Id,
		Sendername: "GM",
		Receiverid: username.Id,
		Title:      "ItemCode",
		Content:    "กด Claim All หรือ Claim เพื่อรับของ",
		Gold:       0,
		Cash:       0,
		Currencies: " ",
		Items:      "",
	}
	if err := db.AOK_DB.Save(&sendMail).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": "error", "msg": "ไม่สามารถทำรายการได้ บันทึกไม่สำเร็จ !"})
		return
	}
	// ระบบได้ส่งของเข้าในเกมให้แล้ว
	tdata.Status = 1

	if err := db.Conn.Save(&tdata).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status": "error", "msg": "ไม่สามารถทำรายการได้ บันทึกlogไม่สำเร็จ !"})
		return
	}

	fmt.Println("Tdata:", tdata)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "ไอเทมโค้ดเปิดใช้งาน รับไอเทมสำเร็จ",
	})
}
