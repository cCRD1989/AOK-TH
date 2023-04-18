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
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zalando/gin-oauth2/google"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
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

var (
	// AuthURL:  "https://www.facebook.com/v16.0/dialog/oauth",
	// TokenURL: "https://graph.facebook.com/v16.0/oauth/access_token",
	OauthConf = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "https://ageofkhagan/auth/facebookCall",
		//RedirectURL: "https://localhost/",
		Scopes:   []string{"public_profile", "email"},
		Endpoint: facebook.Endpoint,
	}
	OauthStateString = "thisshouldberandom"
)

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
	ctx.HTML(http.StatusOK, "frontend/index.html", gin.H{
		"title": "Age Of Khagan Thailand",
		"user":  user,
	})
}

func (f *Frontend) UserGetLogin(ctx *gin.Context) {

	Form := aokmodel.Userlogin{Username: ctx.PostForm("username"), Password: ctx.PostForm("password")}
	fmt.Println("Form: ", Form)

	user := aokmodel.Userlogin{}
	user = user.FindUserByName(Form.Username)
	fmt.Println("user: ", user)
	// Check ID
	if user.Id == "" {
		ctx.Redirect(http.StatusFound, "/")
		fmt.Println("Check ID")
		return
	}

	// CompareHashAndPassword MD5
	if unit.HashMD5(Form.Password) != user.Password {
		ctx.Redirect(http.StatusFound, "/")
		fmt.Println("Check Pass")
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
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	// SetCookie
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	// Redirect
	fmt.Println("บันทึก Token สำเร็จ")
	ctx.Redirect(http.StatusFound, "/")
}

// UserGetLogout logs the user out
func (f *Frontend) UserGetLogout(ctx *gin.Context) {
	tokenString, err := ctx.Cookie("Authorization")
	if err == nil {
		ctx.SetCookie("Authorization", tokenString, -1, "", "", false, true)
	}
	ctx.Redirect(http.StatusFound, "/")
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
