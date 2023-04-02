package controller

import (
	"ccrd/db"
	"ccrd/model"
	"ccrd/model/aokmodel"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zalando/gin-oauth2/google"
)

type Frontend struct{}

func (f *Frontend) UserGetHome(ctx *gin.Context) {

	visit := model.LogWeb{
		DataType:  "visit",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	ctx.HTML(http.StatusOK, "frontend/index.html", gin.H{
		"title": "Age Of Khagan Thailand",
	})
}

func (f *Frontend) UserGetDownload(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "Window" {
		visit := model.LogWeb{
			DataType:  "download_window",
			IPAddress: ctx.ClientIP(),
		}
		db.Conn.Save(&visit)
		ctx.Redirect(http.StatusTemporaryRedirect, "public/download/window.rar")
	} else if id == "Android" {
		visit := model.LogWeb{
			DataType:  "download_android",
			IPAddress: ctx.ClientIP(),
		}
		db.Conn.Save(&visit)
		ctx.Redirect(http.StatusTemporaryRedirect, "public/download/android.rar")
	} else if id == "App" {
		visit := model.LogWeb{
			DataType:  "download_App",
			IPAddress: ctx.ClientIP(),
		}
		db.Conn.Save(&visit)
		ctx.Redirect(http.StatusTemporaryRedirect, "public/download/app.rar")
	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
}

func (f *Frontend) Auth_google(ctx *gin.Context) {

	name := ctx.MustGet("user").(google.User)

	//ctx.JSON(http.StatusOK, gin.H{"Hello": "from private", "user": name})
	//fmt.Println("Auth_google: ", name)

	//บันทึก Log  LogRegister Wait
	db.Conn.Save(&model.LogRegister{
		Sub:      name.Sub,
		Email:    name.Email,
		Name:     name.Name,
		Img:      name.Picture,
		Username: "",
		Password: "",
		Status:   "Wait",
	})

	ctx.HTML(http.StatusOK, "frontend/auth.html", gin.H{
		"title":  "Age Of Khagan Thailand | Account",
		"email":  name.Email,
		"name":   name.Name,
		"imgsrc": name.Picture,
		"status": "",
		"sub":    name.Sub,
	})
}

func (f *Frontend) Auth_google_Regis(ctx *gin.Context) {
	idcode := ctx.DefaultQuery("idcode", "-")
	email := ctx.DefaultQuery("email", "-")
	pass := ctx.DefaultQuery("password", "-")
	repass := ctx.DefaultQuery("repassword", "-")

	fmt.Println("email", email)
	fmt.Println("pass", pass)
	fmt.Println("rpass", repass)

	username1 := strings.Split(email, "@")[0]

	//ตรวจสอบไอดีในระบบ
	if err := db.AOK_DB.First(&aokmodel.Userlogin{}, "username = ?", username1).Error; err == nil {
		fmt.Println("ไอดีซ้ำ")
		ctx.HTML(http.StatusOK, "frontend/auth.html", gin.H{
			"title":  "Age Of Khagan Thailand",
			"tirle1": "ไอดีนี มีอยู่ในระบบ ไม่สารถใช้ไอดีนี้ได้",
			"status": "true",
		})
		return
	}

	//ตรวจสอบพาสตรงกัน
	if pass != repass {
		fmt.Println("ไอดีซ้ำ")
		ctx.HTML(http.StatusOK, "frontend/auth.html", gin.H{
			"title":  "Age Of Khagan Thailand",
			"tirle1": "พาสไม่ตรงกัน",
			"status": "true",
		})
		return
	}

	// บันทึก
	h := md5.New()
	io.WriteString(h, pass)
	passSig := hex.EncodeToString(h.Sum(nil))
	logid := aokmodel.Userlogin{
		Id:       idcode,
		Username: username1,
		Password: passSig,
		Email:    email,
	}
	if err := db.AOK_DB.Save(&logid).Error; err != nil {
		fmt.Println("บันทึกไอดี ไม่สำเร็จ")
		ctx.HTML(http.StatusOK, "frontend/auth.html", gin.H{
			"title":  "Age Of Khagan Thailand",
			"tirle1": "ระบบไม่สามารถ บันทึกข้อมูลของท่านได้",
			"status": "true",
		})
		return
	}

	//บันทึก Log  LogRegister Success
	db.Conn.Model(&model.LogRegister{}).Where("sub = ?", idcode).Updates(model.LogRegister{Status: "Success", Username: username1, Password: passSig})

	ctx.HTML(http.StatusOK, "frontend/auth.html", gin.H{
		"title":  "Age Of Khagan Thailand | Sign up successfully",
		"tirle1": "Sign up Successfully",
		"email":  username1,
		"pass":   pass,
		"name":   "name.Name",
		"imgsrc": "/public/data/รวมไฟล์ 2D by มีน/Standy Rol/cleric.png",
		"status": "true",
	})
}

func (f *Frontend) Auth_custom(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "frontend/customregis.html", gin.H{
		"title":  "Age Of Khagan | Custom Registration",
		"imgsrc": "/public/data/รวมไฟล์ 2D by มีน/Standy Rol/knight.png",
		"name":   "Custom Registration",
		"status": "false",
	})
}

func (f *Frontend) Auth_custom_regis(ctx *gin.Context) {

	userID := ctx.DefaultQuery("username", "-")
	email := ctx.DefaultQuery("email", "-")
	pass := ctx.DefaultQuery("password", "-")
	repass := ctx.DefaultQuery("repassword", "-")

	if userID == "-" || email == "-" || pass == "-" || repass == "-" {
		ctx.HTML(http.StatusOK, "frontend/customregis.html", gin.H{
			"title":  "Age Of Khagan | Custom Registration",
			"imgsrc": "/public/data/รวมไฟล์ 2D by มีน/Standy Rol/knight.png",
			"name":   "กรอกข้อมูลให้ครบ",
			"status": "false",
		})
		return
	}

	//ตรวจสอบพาสตรงกัน
	if pass != repass {
		ctx.HTML(http.StatusOK, "frontend/customregis.html", gin.H{
			"title":  "Age Of Khagan | Custom Registration",
			"imgsrc": "/public/data/รวมไฟล์ 2D by มีน/Standy Rol/knight.png",
			"name":   "Password ไม่ตรงกัน",
			"status": "false",
		})
		return
	}

	//ตรวจสอบไอดีในระบบ
	if err := db.AOK_DB.First(&aokmodel.Userlogin{}, "username = ?", userID).Error; err == nil {

		ctx.HTML(http.StatusOK, "frontend/customregis.html", gin.H{
			"title":  "Age Of Khagan | Custom Registration",
			"imgsrc": "/public/data/รวมไฟล์ 2D by มีน/Standy Rol/knight.png",
			"name":   "Username มีอยู่ในระบบแล้ว โปรดลองใหม่",
			"status": "false",
		})
		return
	}

	// สุ่ม IDCode
	h := md5.New()
	io.WriteString(h, strconv.Itoa(rand.Int()))
	idcode := hex.EncodeToString(h.Sum(nil))

	// เข้ารหัส พาสเวด
	h = md5.New()
	io.WriteString(h, pass)
	passSig := hex.EncodeToString(h.Sum(nil))

	//บันทึกลงฐานข้อมูล
	logid := aokmodel.Userlogin{
		Id:       idcode,
		Username: userID,
		Password: passSig,
		Email:    email,
	}
	if err := db.AOK_DB.Save(&logid).Error; err != nil {
		ctx.HTML(http.StatusOK, "frontend/customregis.html", gin.H{
			"title":  "Age Of Khagan | Custom Registration",
			"imgsrc": "/public/data/รวมไฟล์ 2D by มีน/Standy Rol/knight.png",
			"name":   "บันทึกลงฐานข้อมูลไม่สำเร็จ Error",
			"status": "false",
		})
		return
	}

	//บันทึก Log  LogRegister
	db.Conn.Save(&model.LogRegister{
		Sub:      idcode,
		Email:    email,
		Name:     "",
		Img:      "",
		Username: userID,
		Password: passSig,
		Status:   "Custom Registration",
	})

	ctx.HTML(http.StatusOK, "frontend/customregis.html", gin.H{
		"title":  "Age Of Khagan | Custom Registration",
		"imgsrc": "/public/data/รวมไฟล์ 2D by มีน/Standy Rol/knight.png",
		"name":   "บันทึกลงฐานข้อมูลสำเร็จ",
		"status": "true",
	})

}
