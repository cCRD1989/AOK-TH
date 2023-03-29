package controller

import (
	"ccrd/dto"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Topup struct{}

func Paytopups(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "frontend/topup.html", gin.H{
		"title": "Age Of Khagan Thailand | Topup เติมเงิน",
		"css":   "topup.css",
	})
}

func UserCheck(ctx *gin.Context) {

	user := ctx.Param("user")
	userId := ctx.DefaultQuery("username", "nil")

	fmt.Println("user: ", user)
	fmt.Println("userId: ", userId)

	if user == "user" {

		ctx.HTML(http.StatusOK, "frontend/topup.html", gin.H{
			"title":      "Age Of Khagan | Topup เติมเงิน",
			"status":     "true",
			"userId":     userId,
			"bg_success": "bg-success",
		})
	} else {
		fmt.Println("ไม่เจอข้อมูลไดๆ")
	}

	// ctx.HTML(http.StatusOK, "frontend/topup.html", gin.H{
	// 	"title": "Age Of Khagan | Topup เติมเงิน",
	// })
}

func Payment(ctx *gin.Context) {

	usernameId := ctx.Query("usernameId")
	channel := ctx.Query("channel")
	price := ctx.Query("price")

	h := md5.New()
	io.WriteString(h, strconv.Itoa(rand.Int()))
	orderid := hex.EncodeToString(h.Sum(nil))

	//fmt.Println("orderid:", orderid)

	if usernameId == "" || channel == "" || price == "" {
		// ctx.JSON(http.StatusBadRequest, gin.H{
		// 	"error": "input error.",
		// })
		
		return
	}

	var forr = os.Getenv("FOR") + "-" + orderid
	var operator = ""
	var sid = os.Getenv("SID")
	var uid = usernameId
	var SECRET_KEY = os.Getenv("SECRET_KEY")

	var urlA *url.URL
	var err error
	if channel == "truewallet" {
		urlA, err = url.Parse("https://sea-api.gold-sandbox.razer.com/ewallet/pay?channel=&for=&orderid=&sid=&uid=&price=&sig=")
		if err != nil {
			log.Fatal("RUL Payment error :", err)
		}
	} else {
		urlA, err = url.Parse("https://sea-api.gold-sandbox.razer.com/ibanking/pay?channel=&for=&orderid=&sid=&uid=&price=&sig=")
		if err != nil {
			log.Fatal("RUL Payment error :", err)
		}
	}

	data := channel + forr + operator + orderid + price + sid + uid + SECRET_KEY

	h = md5.New()
	io.WriteString(h, data)
	sumSig := hex.EncodeToString(h.Sum(nil))

	urladdpara := urlA.Query()
	urladdpara.Set("channel", channel)
	urladdpara.Set("for", forr)
	urladdpara.Set("orderid", orderid)
	urladdpara.Set("sid", sid)
	urladdpara.Set("uid", uid)
	urladdpara.Set("price", price)
	urladdpara.Set("sig", sumSig)

	urlA.RawQuery = urladdpara.Encode()

	fmt.Println("sendURL: ", urlA.String())
	ctx.Redirect(http.StatusTemporaryRedirect, urlA.String())
}

// notification Paytopup
func (t *Topup) Paytopup(ctx *gin.Context) {

	fmt.Println("notification In")
	// var request dto.TopupRequest
	// if err := ctx.ShouldBindJSON(&request); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	// 	return
	// }

	request := dto.TopupRequest{
		Txid:     ctx.Query("txid"),
		Orderid:  ctx.Query("orderid"),
		Status:   ctx.Query("status"),
		Detail:   ctx.Query("detail"),
		Channel:  ctx.Query("channel"),
		Amount:   ctx.Query("amount"),
		Currency: ctx.Query("currency"),
		Sig:      ctx.Query("sig"),
	}
	fmt.Println("notification data all :", request)

	data := request.Amount + request.Channel + request.Currency + request.Detail + request.Orderid + request.Status + request.Txid + os.Getenv("SECRET_KEY")

	h := md5.New()
	io.WriteString(h, data)
	sumSig := hex.EncodeToString(h.Sum(nil))

	if request.Sig == sumSig {
		//
		//Code..
		fmt.Println("Sig ตรง ผ่าน")
		//
		ctx.JSON(http.StatusOK, dto.TopupResponse{
			Txid:   request.Txid,
			Status: "200",
		})
	} else {
		fmt.Println("Sig ไม่ตรง")
		fmt.Println("data", data)
		fmt.Println("old", request.Sig)
		fmt.Println("new", sumSig)
		ctx.JSON(http.StatusOK, dto.TopupResponse{
			Txid:   request.Txid,
			Status: "609",
		})
	}
}

// Redirect PayProcess
func (t *Topup) PayProcess(ctx *gin.Context) {

	fmt.Println("Redirect PayProcess")

	request := dto.TopupRequest{
		Txid:     ctx.Query("txid"),
		Orderid:  ctx.Query("orderid"),
		Status:   ctx.Query("status"),
		Detail:   ctx.Query("detail"),
		Channel:  ctx.Query("channel"),
		Amount:   ctx.Query("amount"),
		Currency: ctx.Query("currency"),
		Sig:      ctx.Query("sig"),
	}
	fmt.Println("Redirect PayProcess data all: ", request)

	if request.Status == "200" {
		fmt.Println("PayProcess: ", "Succeeding")
		ctx.HTML(http.StatusOK, "frontend/topupdon.html", gin.H{
			"title":    "Age Of Khagan | Succeeding.",
			"sum":      "Succeeding",
			"txid":     request.Txid,
			"orderid":  request.Orderid,
			"status":   request.Status,
			"detail":   request.Detail,
			"channel":  request.Channel,
			"amount":   request.Amount,
			"currency": request.Currency,
			"sig":      request.Sig,
		})
	} else {
		fmt.Println("PayProcess: ", "Failed")
		ctx.HTML(http.StatusOK, "frontend/topupdon.html", gin.H{
			"title":    "Age Of Khagan | Failed.",
			"sum":      "Failed",
			"txid":     "Failed",
			"orderid":  "Failed",
			"status":   "Failed",
			"detail":   "Failed",
			"channel":  "Failed",
			"amount":   "Failed",
			"currency": "Failed",
			"sig":      "Failed",
		})
	}

}
