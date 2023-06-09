package controller

import (
	"ccrd/db"
	"ccrd/model"
	"ccrd/model/aokmodel"
	"ccrd/unit"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/zalando/gin-oauth2/google"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"gorm.io/gorm"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Frontend struct{}

// FaceBook
type Message struct {
	Name     string
	Id       string
	Likes    string
	Gender   string
	Birthday string
}

type cRegister struct {
	Username   string
	Email      string
	Password   string
	Repassword string
}

var (
	// AuthURL:  "https://www.facebook.com/v16.0/dialog/oauth",
	// TokenURL: "https://graph.facebook.com/v16.0/oauth/access_token",
	OauthConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "https://www.ageofkhaganth.com/auth/facebookCall",
		//RedirectURL: "https://localhost/",
		Scopes:   []string{"public_profile", "email"},
		Endpoint: facebook.Endpoint,
	}
	OauthStateString = "thisshouldberandom"
)

// ////////////////////////////////////////////////////////////////
// M map type string of interfaces
type M map[string]interface{}

// Model struct
type Model struct {
	gin          *gin.Context
	QuerySearch  string //อาชีพ
	QueryKeyword string //ค่าต่างๆ

	Errors M
}

func (model *Model) addError(i string, v interface{}) {
	if model.Errors == nil {
		model.Errors = make(M)
	}
	model.Errors[i] = v
}

// NewModel model
func NewModel(ctx *gin.Context) *Model {
	var model Model

	model.gin = ctx
	model.QuerySearch = "allclass"
	model.QueryKeyword = "level"

	return &model
}

// FindAll models
func (model *Model) FindAll(x interface{}) *Model {
	var err error
	var dbt = db.AOK_DB

	var c = model.gin
	var qKeyword = c.DefaultQuery("qkeyword", model.QueryKeyword)

	switch v := x.(type) {

	case *[]aokmodel.Character:

		if qKeyword == "level" || qKeyword == "durlukin" || qKeyword == "nurin" {
			err = model.buildSQL(dbt.Model(v)).Find(v).Error
		} else {
			err = model.buildSQL(dbt.Model(v)).Joins("JOIN characterattribute ON characters.id = characterattribute.characterId").Find(v).Error
		}

	default:
		fmt.Println("default", reflect.TypeOf(x))

	}

	if err != nil {
		fmt.Println("database")
		model.addError("database", err.Error())
	}

	return model

}

// BuildSQL
func (model *Model) buildSQL(db *gorm.DB) *gorm.DB {

	var c = model.gin

	// Get
	var qSearch = c.DefaultQuery("jobclass", model.QuerySearch)
	var qKeyword = c.DefaultQuery("qkeyword", model.QueryKeyword)

	model.QuerySearch = qSearch
	model.QueryKeyword = qKeyword

	Job := map[string]int{
		"assassin":    970178100,  //
		"necromancer": 479184257,  //
		"micko":       512936679,  //
		"knight":      1817826663, //
		"sorcerer":    607677489,  //
		"cleric":      -859687870,
		"allclass":    0,
	}
	attribute := map[string]int{
		"str": 376173797,
		"con": -1554483258,
		"dex": 1881519984,
		"int": -486441535,
		"wis": 459649607,
		"cha": 1746293732,
	}

	ClassId := Job[qSearch]
	attributeindex := attribute[qKeyword]

	fmt.Println("qSearch", qSearch)
	fmt.Println("ClassId", ClassId)
	fmt.Println("qKeyword", qKeyword)

	if qSearch != "" && qKeyword != "" {

		if qKeyword == "level" {
			if ClassId == 0 {

				db.Select("characters.Id, Userid, characters.Dataid, Charactername, Level, Factionid, Currenthp, Currentmp, Guildid").Limit(10).Order("LEVEL DESC")
				//db.Select("characters.Id, Userid, characters.Dataid, Charactername, Level, Factionid, Currenthp, Currentmp, Guildid").Limit(10).Order("LEVEL DESC")
			} else {
				db.Select("characters.Id, Userid, characters.Dataid, Charactername, Level, Factionid, Currenthp, Currentmp, Guildid").Where("characters.Dataid = ?", ClassId).Limit(10).Order("LEVEL DESC")
			}
		} else if qKeyword == "hp" {
			if ClassId == 0 {
				db.Select("characters.Id, Userid, characters.Dataid, Charactername, Level, Factionid, Currenthp, Currentmp, Guildid").Limit(10).Order("Currenthp DESC")

			} else {
				db.Select("characters.Id, Userid, characters.Dataid, Charactername, Level, Factionid, Currenthp, Currentmp, Guildid").Where("characters.Dataid = ?", ClassId).Limit(10).Order("characters.Currenthp DESC")

			}
		} else if qKeyword == "durlukin" {
			//-1216203983
			if ClassId == 0 {
				db.Select("characters.Id, Userid, characters.Dataid, Charactername, Level, Factionid, Currenthp, Currentmp, Guildid").Where("Factionid = ?", -1216203983).Limit(10).Order("Level DESC")
			} else {
				db.Select("characters.Id, Userid, characters.Dataid, Charactername, Level, Factionid, Currenthp, Currentmp, Guildid").Where("characters.Dataid = ?", ClassId).Where("characters.Factionid = ?", -1216203983).Limit(10).Order("Level DESC")
			}
		} else if qKeyword == "nurin" {
			//1256175201
			if ClassId == 0 {
				db.Select("characters.Id, Userid, characters.Dataid, Charactername, Level, Factionid, Currenthp, Currentmp, Guildid").Where("Factionid = ?", 1256175201).Limit(10).Order("Level DESC")
			} else {
				db.Select("characters.Id, Userid, characters.Dataid, Charactername, Level, Factionid, Currenthp, Currentmp, Guildid").Where("characters.Dataid = ?", ClassId).Where("characters.Factionid = ?", -1216203983).Limit(10).Order("Level DESC")
			}
		} else if qKeyword == "str" ||
			qKeyword == "dex" ||
			qKeyword == "con" ||
			qKeyword == "int" ||
			qKeyword == "wis" ||
			qKeyword == "cha" {
			fmt.Println("attributeindex", attributeindex)
			if ClassId == 0 {
				db.Select("characters.Id, Userid, characters.Dataid, Charactername, Level, Factionid, Currenthp, Currentmp, Guildid, characterattribute.amount AS Attribute").Where("characterattribute.dataId = ?", attributeindex).Limit(10).Order("characterattribute.amount DESC")
			} else {
				db.Select("characters.Id, Userid, characters.Dataid, Charactername, Level, Factionid, Currenthp, Currentmp, Guildid, characterattribute.amount AS Attribute").Where("characterattribute.dataId = ?", attributeindex).Where("characters.Dataid = ?", ClassId).Limit(10).Order("characterattribute.amount DESC")
			}

		} else {
			db.Select("Id, Userid, Dataid, Charactername, Level, Factionid, Currenthp, Currentmp, Guildid").Limit(10).Order("LEVEL DESC")
		}

	} else {

		db.Select("Id, Userid, Dataid, Charactername, Level, Factionid, Currenthp, Currentmp, Guildid").Limit(10).Order("LEVEL DESC")
	}

	return db
}

//////////////////////////////////////////////////////////////////

func (f *Frontend) UserGetTest(ctx *gin.Context) {

	logs := []aokmodel.Character{}

	var logsModel = NewModel(ctx).FindAll(&logs)

	ctx.HTML(http.StatusOK, "frontend/test.html", gin.H{
		"title":     "Age Of Khagan Thailand",
		"bg":        "/public/data/img/main-bg.png",
		"logs":      logs,
		"logsModel": logsModel,
	})
}

func (f *Frontend) UserGetHome(ctx *gin.Context) {

	visit := model.LogWeb{
		DataType:  "visit",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	//
	//
	//
	//

	logs := []aokmodel.Character{}
	var logsModel = NewModel(ctx).FindAll(&logs)

	guild := aokmodel.Guild{}

	//attribute := aokmodel.Characterattribute{}

	for i := 0; i < len(logs); i++ {

		db.AOK_DB.Select("Guildname").Where("id = ?", logs[i].Guildid).Find(&guild)

		if guild.Guildname == "" {
			logs[i].Guildids = "-"
		} else {
			logs[i].Guildids = guild.Guildname
		}
	}

	iconuser := "/public/data/img/user" + strconv.Itoa(rand.Intn(4-1)+1) + ".png"

	//
	//

	//
	//
	new_all := []model.LogNews{}
	new_event := []model.LogNews{}
	new_other := []model.LogNews{}

	db.Conn.Order("created_at DESC").Limit(5).Find(&new_all)
	db.Conn.Where("datatype = ?", "Event").Order("created_at DESC").Limit(20).Find(&new_event)
	db.Conn.Where("datatype = ?", "Other").Order("created_at DESC").Limit(20).Find(&new_other)

	ctx.HTML(http.StatusOK, "frontend/index.html", gin.H{
		"title":     "Age Of Khagan Thailand",
		"user":      user,
		"logsModel": logsModel,
		"logs":      logs,
		"bg":        "/public/data/img/main-bg-top.png",
		"iconuser":  iconuser,
		"new_all":   new_all,
		"new_event": new_event,
		"new_other": new_other,
	})
}

func (f *Frontend) UserGetSingin(ctx *gin.Context) {

	visit := model.LogWeb{
		DataType:  "Singin",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	//
	ctx.HTML(http.StatusOK, "frontend/login.html", gin.H{
		"title": "Age Of Khagan Thailand | Login",
		"user":  user,
		"bg":    "/public/data/img/LOGIN-BG.png",
	})
}

func (f *Frontend) UserGetLogin(ctx *gin.Context) {

	Form := aokmodel.Userlogin{Username: ctx.PostForm("username"), Password: ctx.PostForm("password")}
	fmt.Println("Form: ", Form)

	errs := unit.Validate(map[string]interface{}{
		"username": Form.Username,
	}, map[string]string{

		"username": "required|min:6|max:15|alphanum",
	})

	if errs != nil {
		fmt.Println("FindUserByName Validate Error.")
		ctx.Redirect(http.StatusFound, "/")
		fmt.Println("UserGetLogin Check ID Validate")
		return
	}

	user := aokmodel.Userlogin{}
	user = user.FindUserByName(Form.Username)
	fmt.Println("user: ", user)

	//Chaeck ID index
	if user.Isemailverified == 99 {

		ctx.HTML(http.StatusOK, "frontend/login.html", gin.H{
			"title":   "Age Of Khagan Thailand | Login",
			"user":    aokmodel.Userlogin{},
			"bg":      "/public/data/img/LOGIN-BG.png",
			"message": "ไอดีนี้จะดำเนินการลบในอีก 30 วัน",
		})
	}

	// Check ID
	if user.Id == "" {
		ctx.HTML(http.StatusOK, "frontend/login.html", gin.H{
			"title":   "Age Of Khagan Thailand | Login",
			"user":    aokmodel.Userlogin{},
			"bg":      "/public/data/img/LOGIN-BG.png",
			"message": "ไอดีนี้ไม่มีในระบบ",
		})
	}

	// CompareHashAndPassword MD5
	if unit.HashMD5(Form.Password) != user.Password {
		fmt.Println("Check Pass")
		ctx.HTML(http.StatusOK, "frontend/login.html", gin.H{
			"title":   "Age Of Khagan Thailand | Login",
			"user":    user,
			"bg":      "/public/data/img/LOGIN-BG.png",
			"message": "พาสเวิร์ดไม่ถูกต้อง โปรดลองใหม่",
		})
		return
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("MY_SECRET_KEY")))
	if err != nil {
		fmt.Println("Sign Token")
		ctx.HTML(http.StatusOK, "frontend/login.html", gin.H{
			"title":   "Age Of Khagan Thailand | Login",
			"user":    user,
			"bg":      "/public/data/img/LOGIN-BG.png",
			"message": "Token โปรดลองใหม่",
		})
		return
	}

	// SetCookie
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	// Redirect
	fmt.Println("บันทึก Token สำเร็จ")
	ctx.Redirect(http.StatusFound, "/profile")
}

func (f *Frontend) UserGetRegister(ctx *gin.Context) {
	visit := model.LogWeb{
		DataType:  "Register",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	//
	ctx.HTML(http.StatusOK, "frontend/register.html", gin.H{
		"title":   "Age Of Khagan Thailand | Register",
		"user":    user,
		"bg":      "/public/data/img/REGISTER-BG.png",
		"message": "",
	})
}

func (f *Frontend) UserGetClass(ctx *gin.Context) {
	visit := model.LogWeb{
		DataType:  "Class",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	//
	ctx.HTML(http.StatusOK, "frontend/classjob.html", gin.H{
		"title": "Age Of Khagan Thailand | Class",
		"user":  user,
		"bg":    "/public/data/img/CLASS_BG.png",
	})
}

func (f *Frontend) UserGetMaps(ctx *gin.Context) {
	visit := model.LogWeb{
		DataType:  "Maps",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	//
	ctx.HTML(http.StatusOK, "frontend/mapselect.html", gin.H{
		"title": "Age Of Khagan Thailand | Maps",
		"user":  user,
		"bg":    "/public/data/img/LOGIN-BG.png",
	})
}

func (f *Frontend) UserGetProfile(ctx *gin.Context) {

	visit := model.LogWeb{
		DataType:  "Profile",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	// เช็ค ว่ามีการยืนยันตัวตนหรือไหม CheckVerifyEmail
	// verify := aokmodel.Userlogin{}
	// verify = verify.FindUserByName(user.Username)
	// if verify.Isemailverified == 0 {
	// 	ctx.HTML(http.StatusOK, "frontend/profileverify.html", gin.H{
	// 		"title": "Age Of Khagan Thailand | Login",
	// 		"user":  user,
	// 		"bg":    "/public/data/img/LOGIN-BG.png",
	// 	})
	// 	return
	// }

	//
	//
	//log เติมเงิน

	logtopup := []model.LogTopup{}
	db.Conn.Where("user_id", user.Username).Where("data_type = ?", "NotificationTopup").Order("created_at DESC").Limit(7).Find(&logtopup)

	//
	ctx.HTML(http.StatusOK, "frontend/profile.html", gin.H{
		"title":    "Age Of Khagan Thailand | Login",
		"user":     user,
		"bg":       "/public/data/img/LOGIN-BG.png",
		"logtopup": logtopup,
		// "submitregis": "",
	})
}

func (f *Frontend) UserGetProfileVerifySendNew(ctx *gin.Context) {
	visit := model.LogWeb{
		DataType:  "ProfileVerify",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	//CheckVerifyEmail
	verify := aokmodel.Userlogin{}
	verify = verify.FindUserByName(user.Username)
	if verify.Isemailverified != 1 {

		f.UserEmailVerifySend(verify.Username, verify.Id, verify.Email)

		ctx.HTML(http.StatusOK, "frontend/profileverify.html", gin.H{
			"title": "Age Of Khagan Thailand | Login",
			"user":  user,
			"bg":    "/public/data/img/LOGIN-BG.png",
			"msg":   "on",
		})
		return
	}

	ctx.HTML(http.StatusOK, "frontend/profile.html", gin.H{
		"title": "Age Of Khagan Thailand | Login",
		"user":  user,
		"bg":    "/public/data/img/LOGIN-BG.png",
	})
}

func (f *Frontend) UserGetChangPass(ctx *gin.Context) {
	visit := model.LogWeb{
		DataType:  "ChangPassword",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	Username := user.Username
	Password_old := ctx.PostForm("old_password")
	Password_new := ctx.PostForm("new_password")
	RePassword_new := ctx.PostForm("new_repassword")

	fmt.Println("Username", Username)
	fmt.Println("Password_old", Password_old)
	fmt.Println("Password_new", Password_new)
	fmt.Println("RePassword_new", RePassword_new)

	Form := aokmodel.Userlogin{
		Username: Username,
		Password: Password_old,
	}

	errs := unit.Validate(map[string]interface{}{
		"username": Form.Username,
	}, map[string]string{

		"username": "required|min:6|max:15|alphanum",
	})

	if errs != nil {
		fmt.Println("UserGetChangPass Validate Error.")
		ctx.Redirect(http.StatusFound, "/")
		fmt.Println("UserGetLogin Check ID Validate")
		return
	}

	userDB := aokmodel.Userlogin{}
	userDB = userDB.FindUserByName(Form.Username)

	// Check ID
	if userDB.Id == "" {
		//ctx.Redirect(http.StatusFound, "/")
		fmt.Println("Check ID หาไอดีไม่เจอ")
		return
	}

	// CompareHashAndPassword MD5
	if unit.HashMD5(Form.Password) != userDB.Password {
		//ctx.Redirect(http.StatusFound, "/")
		fmt.Println("Check Pass พาสเดิม ไม่ตรง")
		return
	}

	if Password_new != RePassword_new {
		fmt.Println("พาสใหม่ไม่ตรง")
		return
	}

	// เข้ารหัส พาสเวด
	newPass := unit.HashMD5(RePassword_new)

	db.AOK_DB.Model(&userDB).Update("password", newPass)

	////log เติมเงิน
	logtopup := []model.LogTopup{}
	db.Conn.Where("user_id", user.Username).Where("data_type = ?", "NotificationTopup").Order("created_at DESC").Limit(7).Find(&logtopup)

	ctx.HTML(http.StatusOK, "frontend/profile.html", gin.H{
		"title":       "Age Of Khagan Thailand | Login",
		"user":        user,
		"bg":          "/public/data/img/LOGIN-BG.png",
		"logtopup":    logtopup,
		"submitregis": "",
	})

	ctx.Redirect(http.StatusFound, "/profile")
}

func (f *Frontend) UserGetPDelete(ctx *gin.Context) {

	visit := model.LogWeb{
		DataType:  "Profile",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	//
	//
	//log เติมเงิน

	logtopup := []model.LogTopup{}
	db.Conn.Where("user_id", user.Username).Where("data_type = ?", "NotificationTopup").Order("created_at DESC").Limit(7).Find(&logtopup)

	//
	ctx.HTML(http.StatusOK, "frontend/profiledelete.html", gin.H{
		"title":    "Age Of Khagan Thailand | Login",
		"user":     user,
		"bg":       "/public/data/img/LOGIN-BG.png",
		"logtopup": logtopup,
	})
}

func (f *Frontend) UserGetDelete(ctx *gin.Context) {

	checkbox := ctx.PostForm("checkbox")

	if checkbox == "on" {
		// ตรวจสอบ User Cookie
		usr, _ := ctx.Get("user")
		user, _ := usr.(aokmodel.Userlogin)

		aokuser := aokmodel.Userlogin{}

		db.AOK_DB.Where("username = ?", user.Username).Delete(&aokuser)

		// aokuser.Isemailverified = 99
		// db.AOK_DB.Save(&aokuser)

		fmt.Println("aokuser", aokuser, checkbox)
		fmt.Println("aokuser", aokuser.Isemailverified)

		// //บันทึก Log  LogRegister
		// db.Conn.Save(&model.LogRegister{
		// 	Sub:      idcode,
		// 	Email:    email,
		// 	Name:     "",
		// 	Img:      "",
		// 	Username: user.Username,
		// 	Password: "",
		// 	Status:   "Delete",
		// })

		ctx.HTML(http.StatusOK, "frontend/profiledeletedone.html", gin.H{
			"title": "Age Of Khagan Thailand | Login",
			"user":  user,
			"bg":    "/public/data/img/LOGIN-BG.png",
		})

		ctx.Redirect(http.StatusFound, "/logout")
	}
}

func (f *Frontend) UserGetExDelete(ctx *gin.Context) {

	accessToken := ctx.DefaultQuery("accessToken", "")

	if accessToken == "" {
		fmt.Println("ไม่ได้ส่ง accessToken มา")
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	user := aokmodel.Userlogin{}
	user = user.FindUserByAccesstoken(accessToken)
	fmt.Println("user: ", user)

	if user.Username == "" {
		fmt.Println("คนหาไม่เจอไอดีที่ต้องการลบ")
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	//ลบ Authorization ของเก่าก่อน
	tokenString1, err := ctx.Cookie("Authorization")
	if err == nil {
		ctx.SetCookie("Authorization", tokenString1, -1, "", "", false, true)
	}

	//
	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("MY_SECRET_KEY")))
	if err != nil {
		fmt.Println("Sign Token")
		ctx.HTML(http.StatusOK, "frontend/login.html", gin.H{
			"title":   "Age Of Khagan Thailand | Login",
			"user":    user,
			"bg":      "/public/data/img/LOGIN-BG.png",
			"message": "Token โปรดลองใหม่",
		})
		return
	}

	// SetCookie
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	// Redirect
	fmt.Println("บันทึก Token สำเร็จ")
	ctx.Redirect(http.StatusFound, "/profile/deletion")
}

func (f *Frontend) UserEmailVerifySend(user, Id, email string) {

	from := mail.NewEmail("AOK-TH", "info@ageofkhaganth.com")
	subject := "AOK-TH Verifying your email address."
	to := mail.NewEmail("AOK-TH", email)

	plainTextContent := `
	hello. %s 
	Please verify email
	You’re almost there! We sent an email to Click here to verify your email address. https://%s/email/verify/?idCode=%s
	Just click on the link in that email to complete your singup. If you don’t see it, you may need to check your spam folder.
	`

	htmlContent := `
	<html>
		<head></head>
		<body>
			<p>hello. %s </p>
			<p>Please verify email</p>
			<p>You’re almost there! We sent an email to <a href="https://%s/email/verify/?idCode=%s"><u>Click here to verify your email address.</u></a></p>
			<p></p>
			<p>Just click on the link in that email to complete your singup. If you don’t see it, you may need to check your spam folder.</p>
		</body>
	</html>
	`

	plainTextContent = fmt.Sprintf(plainTextContent, user, "ageofkhaganth.com", Id)
	htmlContent = fmt.Sprintf(htmlContent, user, "ageofkhaganth.com", Id)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)

	if err != nil {
		fmt.Println("ไม่สำเร็จ")
		log.Println(err)
	} else {
		fmt.Println("UserEmailVerifySend สำเร็จ")
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

func (f *Frontend) UserTokenCodeSend(user, email, tokencode string) {

	from := mail.NewEmail("AOK-TH", "info@ageofkhaganth.com")
	subject := "AOK-TH Item Code."
	to := mail.NewEmail("AOK-TH", email)

	plainTextContent := `
	Hello : %s
	<br>
	AOK-TH Item Code : %s
	`

	htmlContent := `
	Hello : %s
	<br>
	AOK-TH Item Code : %s
	`

	plainTextContent = fmt.Sprintf(plainTextContent, user, tokencode)
	htmlContent = fmt.Sprintf(htmlContent, user, tokencode)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("UserTokenCodeSend สำเร็จ")
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

func (f *Frontend) SmtpUserTokenCodeSend(user, email, tokencode string) {
	// Define the SMTP server configuration
	smtpHost := os.Getenv("smtpHost")
	smtpPort := 587
	smtpUsername := os.Getenv("smtpUsername")
	smtpPassword := os.Getenv("smtpPassword")

	// Define the email content
	from := os.Getenv("smtpUsername")
	to := email

	subject := "AOK-TH Item Code."
	body := fmt.Sprintf("Hello : %s\r\nAOK-TH Item Code : %s", user, tokencode)

	// Construct the email message
	message := []byte("To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	// Establish a connection with the SMTP server
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, from, []string{to}, message)
	if err != nil {
		fmt.Println("sendemail err", err)
	} else {
		fmt.Println("Email sent successfully!")

	}
}

func (f *Frontend) UserEmailVerify(ctx *gin.Context) {
	icode := ctx.DefaultQuery("idCode", "")

	if icode == "" {
		ctx.Redirect(http.StatusFound, "/")
	}

	log := aokmodel.Userlogin{}
	log.FindUserById(icode)

	if log.Username == "" {
		ctx.String(http.StatusOK, "Error1.")

	} else {

		if log.Isemailverified == 0 {
			log.Isemailverified = 1

			db.AOK_DB.Updates(&log)

			ctx.HTML(http.StatusOK, "frontend/profileverify.html", gin.H{
				"title": "Age Of Khagan Thailand | Login",
				"user":  log,
				"bg":    "/public/data/img/LOGIN-BG.png",
				"done":  "on",
			})

			return

		} else {
			ctx.String(http.StatusOK, "Error2.")

		}
	}
}

func (f *Frontend) UserGetMonster(ctx *gin.Context) {

	id := ctx.Param("id")
	fmt.Println("idmap", id)

	visit := model.LogWeb{
		DataType:  "Maps",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	//

	mob1 := []string{}
	mob2 := []string{}
	linl1 := ""
	linl2 := ""
	titlename := ""
	titlediscr := ""

	if id == "1" {
		mob1 = []string{
			"/public/data/img/map/1/mob1/Wroc 01 Vulture.png",
			"/public/data/img/map/1/mob1/Wroc 02 Bandit Worrior.png",
			"/public/data/img/map/1/mob1/Wroc 03 Titan.png",
			"/public/data/img/map/1/mob1/Wroc 04 Cave Man.png",
			"/public/data/img/map/1/mob1/Wroc 05 Ettin.png",
			"/public/data/img/map/1/mob1/Wroc 06 Valkyrie.png",
			"/public/data/img/map/1/mob1/Wroc 07 Ghast.png",
			"/public/data/img/map/1/mob1/Wroc 08 Skeleton Soldier.png",
			"/public/data/img/map/1/mob1/Wroc 09 Skeleton Archer.png",
			"/public/data/img/map/1/mob1/Wroc 10 Charon.png",
			"/public/data/img/map/1/mob1/Wroc 11 Broo.png",
			"/public/data/img/map/1/mob1/Wroc 12 Bug Bear.png",
			"/public/data/img/map/1/mob1/Wroc 13 Cursed Ettin.png",
		}
		mob2 = []string{
			"/public/data/img/map/1/mob2/Wroc 01 Vulture.png",
			"/public/data/img/map/1/mob2/Wroc 02 Bandit Worrior.png",
			"/public/data/img/map/1/mob2/Wroc 03 Titan.png",
			"/public/data/img/map/1/mob2/Wroc 04 Cave Man.png",
			"/public/data/img/map/1/mob2/Wroc 05 Ettin.png",
			"/public/data/img/map/1/mob2/Wroc 06 Valkyrie.png",
			"/public/data/img/map/1/mob2/Wroc 07 Ghast.png",
			"/public/data/img/map/1/mob2/Wroc 08 Skeleton Soldier.png",
			"/public/data/img/map/1/mob2/Wroc 09 Skeleton Archer.png",
			"/public/data/img/map/1/mob2/Wroc 10 Charon.png",
			"/public/data/img/map/1/mob2/Wroc 11 Broo.png",
			"/public/data/img/map/1/mob2/Wroc 12 Bug Bear.png",
			"/public/data/img/map/1/mob2/Wroc 13 Cursed Ettin.png",
		}
		linl1 = "/maps/map/1"
		linl2 = "/maps/mob/1"
		titlename = "WROCLAW FORTRESS"
		titlediscr = "จุดยุทธศาสตร์แห่งการรวมตัวของชนเผ่า Durlukin เพื่อเตรียมความพร้อมในการโจมตีกับกองทัพ Nurin"
	} else if id == "2" {
		mob1 = []string{
			"/public/data/img/map/2/mob1/Kara 01 White Wolf.png",
			"/public/data/img/map/2/mob1/Kara 02 Red Wolf.png",
			"/public/data/img/map/2/mob1/Kara 03 Zombie.png",
			"/public/data/img/map/2/mob1/Kara 04 Desetion Swordsman.png",
			"/public/data/img/map/2/mob1/Kara 05 Desertion Archer.png",
			"/public/data/img/map/2/mob1/Kara 06 Desertion Spear.png",
			"/public/data/img/map/2/mob1/Kara 07 Warrior.png",
			"/public/data/img/map/2/mob1/Kara 08 Sorcerer.png",
			"/public/data/img/map/2/mob1/Kara 09 Halbue Seniors.png",
			"/public/data/img/map/2/mob1/Kara 10 Elder.png",
			"/public/data/img/map/2/mob1/Kara 12 Godochoong (เขียว).png",
			"/public/data/img/map/2/mob1/Kara 12 Godochoong (ส้ม).png",
			"/public/data/img/map/2/mob1/Kara 13 Godochoong (แดง).png",
			"/public/data/img/map/2/mob1/Kara 14 Gyochoogsin(S).png",
			"/public/data/img/map/2/mob1/Kara 15 Gyochoogsin(M).png",
			"/public/data/img/map/2/mob1/Kara 16 Sanso.png",
			"/public/data/img/map/2/mob1/Kara 17 Bisasa Guisabso.png",
		}
		mob2 = []string{
			"/public/data/img/map/2/mob2/Kara 01 White Wolf.png",
			"/public/data/img/map/2/mob2/Kara 02 Red Wolf.png",
			"/public/data/img/map/2/mob2/Kara 03 Zombie.png",
			"/public/data/img/map/2/mob2/Kara 04 Desetion Swordsman.png",
			"/public/data/img/map/2/mob2/Kara 05 Desertion Archer.png",
			"/public/data/img/map/2/mob2/Kara 06 Desertion Spear.png",
			"/public/data/img/map/2/mob2/Kara 07 Warrior.png",
			"/public/data/img/map/2/mob2/Kara 08 Sorcerer.png",
			"/public/data/img/map/2/mob2/Kara 09 Halbue Seniors.png",
			"/public/data/img/map/2/mob2/Kara 10 Elder.png",
			"/public/data/img/map/2/mob2/Kara 12 Godochoong (เขียว).png",
			"/public/data/img/map/2/mob2/Kara 12 Godochoong (ส้ม).png",
			"/public/data/img/map/2/mob2/Kara 13 Godochoong (แดง).png",
			"/public/data/img/map/2/mob2/Kara 14 Gyochoogsin(S).png",
			"/public/data/img/map/2/mob2/Kara 15 Gyochoogsin(M).png",
			"/public/data/img/map/2/mob2/Kara 16 Sanso.png",
			"/public/data/img/map/2/mob2/Kara 17 Bisasa Guisabso.png",
		}
		linl1 = "/maps/map/2"
		linl2 = "/maps/mob/2"
		titlename = "KHARAKORUM"
		titlediscr = "" //"จุดยุทธศาสตร์แห่งการรวมตัวของชนเผ่า Nurin เพื่อเตรียมความพร้อมในการโจมตีกับกองทัพ Durlukin"

	} else if id == "3" {
		mob1 = []string{
			"/public/data/img/map/3/mob1/Lub 01 Seicken.png",
			"/public/data/img/map/3/mob1/Lub 02 Seicken Archer.png",
			"/public/data/img/map/3/mob1/Lub 03 Darer Bear.png",
			"/public/data/img/map/3/mob1/Lub 04 Griffon.png",
			"/public/data/img/map/3/mob1/Lub 05 Frost Salamanda.png",
			"/public/data/img/map/3/mob1/Lub 06 Death Bringer.png",
			"/public/data/img/map/3/mob1/Lub 07 Iwarse.png",
			"/public/data/img/map/3/mob1/Lub 08 Minotaur.png",
			"/public/data/img/map/3/mob1/Lub 09 Frost Worm.png",
			"/public/data/img/map/3/mob1/Lub 10 Ice Golem.png",
			"/public/data/img/map/3/mob1/Lub 11 Bone Iwarse.png",
			"/public/data/img/map/3/mob1/Lub 12 White Dragon.png",
			"/public/data/img/map/3/mob1/Lub 13 Mountain Ice Golem.png",
		}
		mob2 = []string{
			"/public/data/img/map/3/mob2/Lub 01 Seicken.png",
			"/public/data/img/map/3/mob2/Lub 02 Seicken Archer.png",
			"/public/data/img/map/3/mob2/Lub 03 Darer Bear.png",
			"/public/data/img/map/3/mob2/Lub 04 Griffon.png",
			"/public/data/img/map/3/mob2/Lub 05 Frost Salamanda.png",
			"/public/data/img/map/3/mob2/Lub 06 Death Bringer.png",
			"/public/data/img/map/3/mob2/Lub 07 Iwarse.png",
			"/public/data/img/map/3/mob2/Lub 08 Minotaur.png",
			"/public/data/img/map/3/mob2/Lub 09 Frost Worm.png",
			"/public/data/img/map/3/mob2/Lub 10 Ice Golem.png",
			"/public/data/img/map/3/mob2/Lub 11 Bone Iwarse.png",
			"/public/data/img/map/3/mob2/Lub 12 White Dragon.png",
			"/public/data/img/map/3/mob2/Lub 13 Mountain Ice Golem.png",
		}
		linl1 = "/maps/map/3"
		linl2 = "/maps/mob/3"
		titlename = "LUBLIN MONGOL FORTRESS"
		titlediscr = "" //"เมืองแห่งหิมะพื้นที่สำหรับนักรบในการต่อต้านเหล่ามอนสเตอร์ที่แข็งแกร่งและ ชั่วร้าย"

	} else if id == "4" {
		mob1 = []string{
			"/public/data/img/map/4/mob1/Iron 01 Black Beetle.png",
			"/public/data/img/map/4/mob1/Iron 02 Wraith.png",
			"/public/data/img/map/4/mob1/Iron 03 Giant Bat.png",
			"/public/data/img/map/4/mob1/Iron 04 Gargoyle.png",
			"/public/data/img/map/4/mob1/Iron 06 Black Phantom.png",
			"/public/data/img/map/4/mob1/Iron 06 Gas Lion.png",
			"/public/data/img/map/4/mob1/Iron 07 Trol.png",
			"/public/data/img/map/4/mob1/Iron 08 Greed Dragon.png",
			"/public/data/img/map/4/mob1/Iron 09 Phantom of the Phantom.png",
		}
		mob2 = []string{
			"/public/data/img/map/4/mob2/Iron 01 Black Beetle.png",
			"/public/data/img/map/4/mob2/Iron 02 Wraith.png",
			"/public/data/img/map/4/mob2/Iron 03 Giant Bat.png",
			"/public/data/img/map/4/mob2/Iron 04 Gargoyle.png",
			"/public/data/img/map/4/mob2/Iron 06 Black Phantom.png",
			"/public/data/img/map/4/mob2/Iron 06 Gas Lion.png",
			"/public/data/img/map/4/mob2/Iron 07 Trol.png",
			"/public/data/img/map/4/mob2/Iron 08 Greed Dragon.png",
			"/public/data/img/map/4/mob2/Iron 09 Phantom of the Phantom.png",
		}
		linl1 = "/maps/map/4"
		linl2 = "/maps/mob/4"
		titlename = "IRON DUNGEON"
		titlediscr = "" //"เหมืองแร่ใต้หุบเขา Karpatian เหมืองแร่โบราณแห่งความท้าทาย กับสภาพของผู้คนที่เปลี่ยนไป ด้วยความโลภและเวทย์มนต์ดำ"

	} else if id == "5" {
		mob1 = []string{
			"/public/data/img/map/5/mob1/Lava 01 Skeleton Soldier.png",
			"/public/data/img/map/5/mob1/Lava 02 Peryton.png",
			"/public/data/img/map/5/mob1/Lava 03 Giant Scorpion.png",
			"/public/data/img/map/5/mob1/Lava 04 Myconid.png",
			"/public/data/img/map/5/mob1/Lava 05 Berserk Zapher.png",
			"/public/data/img/map/5/mob1/Lava 06 Poison Myconid.png",
		}
		mob2 = []string{
			"/public/data/img/map/5/mob2/Lava 01 Skeleton Soldier.png",
			"/public/data/img/map/5/mob2/Lava 02 Peryton.png",
			"/public/data/img/map/5/mob2/Lava 03 Giant Scorpion.png",
			"/public/data/img/map/5/mob2/Lava 04 Myconid.png",
			"/public/data/img/map/5/mob2/Lava 05 Berserk Zapher.png",
			"/public/data/img/map/5/mob2/Lava 06 Poison Myconid.png",
		}
		linl1 = "/maps/map/5"
		linl2 = "/maps/mob/5"
		titlename = "LAVA CANYON"
		titlediscr = "" //"สถานที่น่าค้นหาและมีเสน่ห์ รายล้อมไปด้วยมอนสเตอร์ผู้ปกป้องทรัพย์สมบัติล้ำค่า"

	} else {
		return
	}

	ctx.HTML(http.StatusOK, "frontend/mob.html", gin.H{
		"title":      "Age Of Khagan Thailand | Maps",
		"user":       user,
		"bg":         "/public/data/img/LOGIN-BG.png",
		"titlename":  titlename,
		"mob1":       mob1,
		"mob2":       mob2,
		"type":       "mob",
		"linl1":      linl1,
		"linl2":      linl2,
		"titlediscr": titlediscr,
	})
}

func (f *Frontend) UserGetMap(ctx *gin.Context) {

	id := ctx.Param("id")
	fmt.Println("idmap", id)

	visit := model.LogWeb{
		DataType:  "Maps",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	//

	map1 := []string{}
	map2 := []string{}
	linl1 := ""
	linl2 := ""
	titlename := ""
	titlediscr := ""

	if id == "1" {
		map1 = []string{
			"/public/data/img/map/1/map/1/map1(1).png",
			"/public/data/img/map/1/map/1/map1(2).png",
			"/public/data/img/map/1/map/1/map1(3).png",
			"/public/data/img/map/1/map/1/map1(4).png",
			"/public/data/img/map/1/map/1/map1(5).png",
			"/public/data/img/map/1/map/1/map1(6).png",
			"/public/data/img/map/1/map/1/map1(7).png",
		}

		map2 = []string{
			"/public/data/img/map/1/map/2/map1(1).png",
			"/public/data/img/map/1/map/2/map1(2).png",
			"/public/data/img/map/1/map/2/map1(3).png",
			"/public/data/img/map/1/map/2/map1(4).png",
			"/public/data/img/map/1/map/2/map1(5).png",
			"/public/data/img/map/1/map/2/map1(6).png",
			"/public/data/img/map/1/map/2/map1(7).png",
		}
		linl1 = "/maps/map/1"
		linl2 = "/maps/mob/1"
		titlename = "WROCLAW FORTRESS"
		titlediscr = "จุดยุทธศาสตร์แห่งการรวมตัวของชนเผ่า Durlukin เพื่อเตรียมความพร้อมในการโจมตีกับกองทัพ Nurin"
	} else if id == "2" {
		map1 = []string{
			"/public/data/img/map/2/map1/map2(1).png",
			"/public/data/img/map/2/map1/map2(2).png",
			"/public/data/img/map/2/map1/map2(3).png",
			"/public/data/img/map/2/map1/map2(4).png",
			"/public/data/img/map/2/map1/map2(5).png",
			"/public/data/img/map/2/map1/map2(6).png",
			"/public/data/img/map/2/map1/map2(7).png",
		}

		map2 = []string{
			"/public/data/img/map/2/map2/map2(1).png",
			"/public/data/img/map/2/map2/map2(2).png",
			"/public/data/img/map/2/map2/map2(3).png",
			"/public/data/img/map/2/map2/map2(4).png",
			"/public/data/img/map/2/map2/map2(5).png",
			"/public/data/img/map/2/map2/map2(6).png",
			"/public/data/img/map/2/map2/map2(7).png",
		}
		linl1 = "/maps/map/2"
		linl2 = "/maps/mob/2"
		titlename = "KHARAKORUM"
		titlediscr = "จุดยุทธศาสตร์แห่งการรวมตัวของชนเผ่า Nurin เพื่อเตรียมความพร้อมในการโจมตีกับกองทัพ Durlukin"

	} else if id == "3" {
		map1 = []string{
			"/public/data/img/map/3/map1/map2(1).png",
			"/public/data/img/map/3/map1/map2(2).png",
			"/public/data/img/map/3/map1/map2(3).png",
			"/public/data/img/map/3/map1/map2(4).png",
			"/public/data/img/map/3/map1/map2(5).png",
			"/public/data/img/map/3/map1/map2(6).png",
			"/public/data/img/map/3/map1/map2(7).png",
		}
		map2 = []string{
			"/public/data/img/map/3/map2/map2(1).png",
			"/public/data/img/map/3/map2/map2(2).png",
			"/public/data/img/map/3/map2/map2(3).png",
			"/public/data/img/map/3/map2/map2(4).png",
			"/public/data/img/map/3/map2/map2(5).png",
			"/public/data/img/map/3/map2/map2(6).png",
			"/public/data/img/map/3/map2/map2(7).png",
		}
		linl1 = "/maps/map/3"
		linl2 = "/maps/mob/3"
		titlename = "LUBLIN MONGOL FORTRESS"
		titlediscr = "เมืองแห่งหิมะพื้นที่สำหรับนักรบในการต่อต้านเหล่ามอนสเตอร์ที่แข็งแกร่งและชั่วร้าย"

	} else if id == "4" {
		map1 = []string{
			"/public/data/img/map/4/map1/map3(1).png",
			"/public/data/img/map/4/map1/map3(2).png",
			"/public/data/img/map/4/map1/map3(3).png",
			"/public/data/img/map/4/map1/map3(4).png",
			"/public/data/img/map/4/map1/map3(5).png",
		}
		map2 = []string{
			"/public/data/img/map/4/map2/map3(1).png",
			"/public/data/img/map/4/map2/map3(2).png",
			"/public/data/img/map/4/map2/map3(3).png",
			"/public/data/img/map/4/map2/map3(4).png",
			"/public/data/img/map/4/map2/map3(5).png",
		}
		linl1 = "/maps/map/4"
		linl2 = "/maps/mob/4"
		titlename = "IRON DUNGEON"
		titlediscr = "เหมืองแร่ใต้หุบเขา Karpatian เหมืองแร่โบราณแห่งความท้าทาย กับสภาพของผู้คนที่เปลี่ยนไป ด้วยความโลภและเวทย์มนต์ดำ"

	} else if id == "5" {
		map1 = []string{
			"/public/data/img/map/5/map1/map5(1).png",
			"/public/data/img/map/5/map1/map5(2).png",
			"/public/data/img/map/5/map1/map5(3).png",
			"/public/data/img/map/5/map1/map5(4).png",
			"/public/data/img/map/5/map1/map5(5).png",
		}
		map2 = []string{
			"/public/data/img/map/5/map2/map5(1).png",
			"/public/data/img/map/5/map2/map5(2).png",
			"/public/data/img/map/5/map2/map5(3).png",
			"/public/data/img/map/5/map2/map5(4).png",
			"/public/data/img/map/5/map2/map5(5).png",
		}
		linl1 = "/maps/map/5"
		linl2 = "/maps/mob/5"
		titlename = "LAVA CANYON"
		titlediscr = "สถานที่น่าค้นหาและมีเสน่ห์ รายล้อมไปด้วยมอนสเตอร์ผู้ปกป้องทรัพย์สมบัติล้ำค่า"

	} else {
		return
	}

	ctx.HTML(http.StatusOK, "frontend/map.html", gin.H{
		"title":      "Age Of Khagan Thailand | Maps",
		"user":       user,
		"bg":         "/public/data/img/LOGIN-BG.png",
		"titlename":  titlename,
		"map1":       map1,
		"map2":       map2,
		"type":       "map",
		"linl1":      linl1,
		"linl2":      linl2,
		"titlediscr": titlediscr,
	})
}

func (f *Frontend) UserNewPage(ctx *gin.Context) {
	idNew := ctx.Param("id")
	visit := model.LogWeb{
		DataType:  "NewPages",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	if idNew == "" {
		fmt.Println("reload")
		ctx.Redirect(http.StatusFound, "/")
	}
	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	newOpen := model.LogNews{}
	db.Conn.First(&newOpen, idNew)

	ctx.HTML(http.StatusOK, "frontend/newpage.html", gin.H{
		"title":   "Age Of Khagan Thailand | NewPages",
		"user":    user,
		"bg":      "/public/data/img/NewPage-BG.png",
		"newOpen": newOpen,
	})
}

func (f *Frontend) UserNewAll(ctx *gin.Context) {

	visit := model.LogWeb{
		DataType:  "NewAll",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	newall := []model.LogNews{}
	db.Conn.Order("created_at DESC").Find(&newall)

	ctx.HTML(http.StatusOK, "frontend/newall.html", gin.H{
		"title":  "Age Of Khagan Thailand | NewPages",
		"user":   user,
		"bg":     "/public/data/img/NewPage-BG.png",
		"newall": newall,
	})
}

// privacypolicy
func (f *Frontend) UserGetPrivacypolicy(ctx *gin.Context) {
	visit := model.LogWeb{
		DataType:  "visit",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	ctx.HTML(http.StatusOK, "frontend/privacypolicy.html", gin.H{
		"title": "Age Of Khagan Thailand | PrivacyPolicy",
		"user":  user,
		"bg":    "/public/data/img/CLASS_BG.png",
	})
}

// Service
func (f *Frontend) UserGetService(ctx *gin.Context) {
	visit := model.LogWeb{
		DataType:  "visit",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	ctx.HTML(http.StatusOK, "frontend/service.html", gin.H{
		"title": "Age Of Khagan Thailand | Service",
		"user":  user,
		"bg":    "/public/data/img/CLASS_BG.png",
	})
}

// UserGetLogout logs the user out
func (f *Frontend) UserGetLogout(ctx *gin.Context) {
	tokenString, err := ctx.Cookie("Authorization")
	if err == nil {
		ctx.SetCookie("Authorization", tokenString, -1, "", "", false, true)
	}
	ctx.Redirect(http.StatusFound, "/home")
}

func (f *Frontend) UserGetDownload(ctx *gin.Context) {

	visit := model.LogWeb{
		DataType:  "DownloadPC",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// linkpc := model.Setup{}
	// db.Conn.Where("Type=?", "pc").First(&linkpc)

	// //"Age Of Khagan 0.0.0.1 Setup (x64).exe"
	// link := fmt.Sprintf("public/download/%s", linkpc.Value)

	// if linkpc.Value == "#" {
	// 	ctx.Redirect(http.StatusTemporaryRedirect, "/")
	// } else {
	// 	ctx.Redirect(http.StatusTemporaryRedirect, link)
	// }

	ctx.Redirect(http.StatusTemporaryRedirect, "https://files.ageofkhaganth.com/Age_Of_Khagan_0.0_0.1%20Setup_(x64).exe")
	//https://file.ageofkhaganth.com/Age Of Khagan 0.0.0.1 Setup (x64).exe

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

	username1 := strings.Split(email, "@")[0]

	//ตรวจสอบไอดีในระบบ
	if err := db.AOK_DB.First(&aokmodel.Userlogin{}, "username = ?", username1).Error; err == nil {
		fmt.Println("ไอดีซ้ำ")
		ctx.HTML(http.StatusOK, "frontend/auth.html", gin.H{
			"title":  "Age Of Khagan Thailand",
			"tirle1": "ไอดีนี้ มีอยู่ในระบบ ไม่สารถใช้ไอดีนี้ได้",
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
	db.Conn.Model(&model.LogRegister{}).Where("sub = ?", idcode).Updates(model.LogRegister{Status: "Google", Username: username1, Password: passSig})

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

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	userID := ctx.DefaultQuery("username", "-")
	email := ctx.DefaultQuery("email", "-")
	pass := ctx.DefaultQuery("password", "-")
	repass := ctx.DefaultQuery("repassword", "-")

	data := cRegister{
		Username:   userID,
		Email:      email,
		Password:   pass,
		Repassword: repass,
	}

	if userID == "-" || email == "-" || pass == "-" || repass == "-" {
		ctx.HTML(http.StatusOK, "frontend/register.html", gin.H{
			"title":   "Age Of Khagan | Custom Registration",
			"name":    "กรอกข้อมูลให้ครบ",
			"data":    data,
			"bg":      "/public/data/img/REGISTER-BG.png",
			"user":    user,
			"message": "",
		})
		return
	}

	//ตรวจสอบพาสตรงกัน
	if pass != repass {
		ctx.HTML(http.StatusOK, "frontend/register.html", gin.H{
			"title":   "Age Of Khagan | Custom Registration",
			"name":    "Password ไม่ตรงกัน",
			"data":    data,
			"bg":      "/public/data/img/REGISTER-BG.png",
			"user":    user,
			"message": "",
		})
		return
	}

	//ตรวจสอบไอดีในระบบ
	if err := db.AOK_DB.First(&aokmodel.Userlogin{}, "username = ?", userID).Error; err == nil {
		ctx.HTML(http.StatusOK, "frontend/register.html", gin.H{
			"title":   "Age Of Khagan | Custom Registration",
			"name":    "Username มีอยู่ในระบบแล้ว โปรดลองใหม่",
			"data":    data,
			"bg":      "/public/data/img/REGISTER-BG.png",
			"user":    user,
			"message": "",
		})
		return
	}

	errs := unit.Validate(map[string]interface{}{
		"username":       userID,
		"password":       pass,
		"password_check": repass,
		"email":          email,
	}, map[string]string{

		"username":       "required|min:6|max:15|alphanum",
		"password":       "required|min:6|max:15|alphanum",
		"password_check": "required|same:password",
		"email":          "required|email",
	})

	if errs != nil {
		ctx.HTML(http.StatusOK, "frontend/register.html", gin.H{
			"title":   "Age Of Khagan | Custom Registration",
			"name":    "ตรวจสอบข้อมูลผิดพลาดโปรดลองใหม่",
			"data":    data,
			"bg":      "/public/data/img/REGISTER-BG.png",
			"user":    user,
			"message": "",
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

	// ตรวจสอบ การรับของรางวัล ลงทะเบรียนล่วงหน้าโดน เช็ค อีเมล์ ไม่ซ้ำ
	fmt.Println("Email", email)
	if err := db.AOK_DB.First(&aokmodel.Userlogin{}, "email = ?", email).Error; err != nil {
		TokenCode := model.LogTokenregister{
			Username: userID,
			Email:    email,
			Tokenid:  unit.GenerateSecureToken(32),
			Status:   0,
			IP:       ctx.ClientIP(),
		}

		if err := db.Conn.Save(&TokenCode).Error; err != nil {

			fmt.Println("ไอดีหรืออีเมล์เคยรับไปแล้ว ไม่สามารถบันทึกข้อมูลได้")
		} else {
			f.UserTokenCodeSend(TokenCode.Username, TokenCode.Email, TokenCode.Tokenid)
			//f.SmtpUserTokenCodeSend(TokenCode.Username, TokenCode.Email, TokenCode.Tokenid)
			fmt.Println("ลงทะเบียนล่วงหน้าสำเร็จ")
		}
	} else {
		fmt.Println("เช็คลงทะเบียนล่วงหน้าเมล์ซ้ำ ไม่ได้รับสิทธิ์")
	}

	//บันทึกลงฐานข้อมูล
	logid := aokmodel.Userlogin{
		Id:       idcode,
		Username: userID,
		Password: passSig,
		Email:    email,
	}
	fmt.Println("logid", logid)
	if err := db.AOK_DB.Save(&logid).Error; err != nil {
		ctx.HTML(http.StatusOK, "frontend/register.html", gin.H{
			"title":   "Age Of Khagan | Custom Registration",
			"name":    "บันทึกลงฐานข้อมูลไม่สำเร็จ Error",
			"data":    data,
			"bg":      "/public/data/img/REGISTER-BG.png",
			"user":    user,
			"message": "",
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

	// UserEmail ** **** *VerifySend
	//f.UserEmailVerifySend(logid.Username, logid.Id, logid.Email)

	//ctx.Redirect(http.StatusFound, "/")
	ctx.HTML(http.StatusOK, "frontend/register.html", gin.H{
		"title":   "Age Of Khagan | Custom Registration",
		"name":    "",
		"data":    data,
		"bg":      "/public/data/img/REGISTER-BG.png",
		"user":    user,
		"message": "open",
	})
}

func (f *Frontend) Auth_facebook_login(ctx *gin.Context) {
	OauthConf.ClientID = os.Getenv("facebookclientID")
	OauthConf.ClientSecret = os.Getenv("facebookclentSecret")

	URL, err := url.Parse(OauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	// fmt.Println("URL1", URL)
	parameters := url.Values{}
	parameters.Add("client_id", OauthConf.ClientID)
	parameters.Add("scope", strings.Join(OauthConf.Scopes, ","))
	parameters.Add("redirect_uri", OauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", OauthStateString)
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	// fmt.Println("URL2", url)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (f *Frontend) Auth_facebook_call(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")

	if state != OauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", OauthStateString, state)
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	if code == "" {
		fmt.Println("Code not found..")
		return
	} else {
		token, err := OauthConf.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
			return
		}

		resp, err := http.Get("https://graph.facebook.com/me?access_token=" + url.QueryEscape(token.AccessToken))
		if err != nil {
			fmt.Printf("Get: %s\n", err)
			return
		}
		defer resp.Body.Close()

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("ReadAll: %s\n", err)
			return
		}

		var m Message
		if err := json.Unmarshal([]byte(response), &m); err != nil {
			fmt.Println("err:", err.Error())
		}
		fmt.Println("Message", m)
		ctx.HTML(http.StatusOK, "frontend/authfacebook.html", gin.H{
			"title":  "Age Of Khagan Thailand | Account",
			"email":  "",
			"name":   m.Name,
			"imgsrc": "",
			"status": "",
			"sub":    m.Id,
		})
	}
}

func (f *Frontend) Auth_facebook_regis(ctx *gin.Context) {
	fullname := ctx.DefaultQuery("fullname", "-")
	idcode := ctx.DefaultQuery("idcode", "-")
	username1 := ctx.DefaultQuery("username", "-")
	pass := ctx.DefaultQuery("password", "-")
	repass := ctx.DefaultQuery("repassword", "-")

	fmt.Println("idcode: ", idcode)
	fmt.Println("username1: ", username1)
	fmt.Println("pass: ", pass)
	fmt.Println("repass: ", repass)

	//ตรวจสอบไอดีในระบบ
	if err := db.AOK_DB.First(&aokmodel.Userlogin{}, "username = ?", username1).Error; err == nil {
		ctx.HTML(http.StatusOK, "frontend/authfacebook.html", gin.H{
			"title":  "Age Of Khagan Thailand | Facebook",
			"tirle1": "ไอดีนี้ มีอยู่ในระบบ ไม่สารถใช้ไอดีนี้ได้",
			"status": "true",
		})
		return
	}

	//ตรวจสอบพาสตรงกัน
	if pass != repass {
		ctx.HTML(http.StatusOK, "frontend/authfacebook.html", gin.H{
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
		Email:    "",
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

	//บันทึก Log  LogRegister Wait
	db.Conn.Save(&model.LogRegister{
		Sub:      idcode,
		Email:    "",
		Name:     fullname,
		Img:      "",
		Username: username1,
		Password: passSig,
		Status:   "FaceBook",
	})

	ctx.HTML(http.StatusOK, "frontend/auth.html", gin.H{
		"title":  "Age Of Khagan Thailand | Sign up successfully",
		"tirle1": "Sign up Successfully",
		"email":  username1,
		"pass":   pass,
		"name":   fullname,
		"imgsrc": "/public/data/รวมไฟล์ 2D by มีน/Standy Rol/cleric.png",
		"status": "true",
	})
}
