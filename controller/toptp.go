package controller

import (
	"ccrd/db"
	"ccrd/dto"
	"ccrd/model"
	"ccrd/model/aokmodel"
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

// หน้าเช็ค UserCheck ID ที่กรอกเข้ามา
func UserCheck(ctx *gin.Context) {

	user := ctx.Param("user")
	userId := ctx.DefaultQuery("username", "nil")

	// fmt.Println("user: ", user)
	// fmt.Println("userId: ", userId)

	// เช็คไอดี

	topup := aokmodel.Userlogin{}
	if err := db.AOK_DB.Where("username =?", userId).First(&topup).Error; err != nil {
		ctx.HTML(http.StatusOK, "frontend/topup.html", gin.H{
			"title":      "Age Of Khagan | Topup เติมเงิน",
			"status":     "false",
			"userId":     "",
			"bg_success": "",
			"message":    "ไอดีไม่มีอยู่ในระบบ โปรดลองใหม่อีกครั้ง",
		})
	}

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

// ลูก ค้ากด ออเดอร์ ออกไป ให้ Razer
func Payment(ctx *gin.Context) {

	usernameId := ctx.Query("usernameId")
	channel := ctx.Query("channel")
	price := ctx.Query("price")

	if usernameId == "" || channel == "" || price == "" {
		ctx.HTML(http.StatusOK, "frontend/topup.html", gin.H{
			"title":      "Age Of Khagan | Topup เติมเงิน",
			"status":     "false",
			"userId":     "",
			"bg_success": "",
			"message":    "ไอดีไม่มีอยู่ในระบบ โปรดลองใหม่อีกครั้ง",
		})
		return
	}

	topup := aokmodel.Userlogin{}
	if err := db.AOK_DB.Where("username =?", usernameId).First(&topup).Error; err != nil {
		ctx.HTML(http.StatusOK, "frontend/topup.html", gin.H{
			"title":      "Age Of Khagan | Topup เติมเงิน",
			"status":     "false",
			"userId":     "",
			"bg_success": "",
			"message":    "ไอดีไม่มีอยู่ในระบบ โปรดลองใหม่อีกครั้ง",
		})
	}

	h := md5.New()
	io.WriteString(h, strconv.Itoa(rand.Int()))
	orderid := hex.EncodeToString(h.Sum(nil))

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

	RequestTopup := model.LogTopup{
		DataType:  "RequestTopup",
		UserId:    topup.Username,
		Txid:      "",
		Orderid:   orderid,
		Status:    "Order",
		Detail:    "",
		Channel:   channel,
		Price:     price,
		Sig:       sumSig,
		IPAddress: ctx.ClientIP(),
	}
	if err := db.Conn.Save(&RequestTopup).Error; err != nil {
		fmt.Println("RequestTopup Error", err.Error())
		return
	}

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
		//

		//ดึง เลข ออเดอร์ จากตารางมาเทียบ
		data := model.LogTopup{}
		if err := db.Conn.Where("orderid = ?", request.Orderid).Where("data_type = ?", "RequestTopup").First(&data).Error; err != nil {
			fmt.Println("ค้นหาเลข Orderid ไม่เจอ")
			ctx.JSON(http.StatusOK, dto.TopupResponse{
				Txid:   request.Txid,
				Status: "609",
			})
			return
		}

		//รอดำเนินการ บันทึกเพิ่มอีก log ในส่วนของ NotificationTopup Status:    "Wait"
		db.Conn.Save(&model.LogTopup{
			DataType:  "NotificationTopup",
			UserId:    data.UserId,
			Txid:      request.Txid,
			Orderid:   request.Orderid,
			Status:    "Wait",
			Detail:    "",
			Channel:   request.Channel,
			Price:     request.Amount + request.Currency,
			Sig:       request.Sig,
			IPAddress: ctx.ClientIP(),
		})

		// เงินที่จะเติม
		caseint, err := strconv.Atoi(request.Amount)
		if err != nil {
			fmt.Println("str to int ไม่ได้ ", data.Price)
			ctx.JSON(http.StatusOK, dto.TopupResponse{
				Txid:   request.Txid,
				Status: "609",
			})
			return
		}

		//ดึงเงินที่อยู่ใน id นั้น
		idcash := aokmodel.Userlogin{}
		db.AOK_DB.First(&idcash, "username = ?", data.UserId)

		idcash.Cash += caseint

		if err := db.AOK_DB.Model(&aokmodel.Userlogin{}).Where("username = ?", idcash.Username).Update("cash", idcash.Cash).Error; err != nil {
			fmt.Println("บันทึกแคชไม่สำเร็จ", err.Error())
			ctx.JSON(http.StatusOK, dto.TopupResponse{
				Txid:   request.Txid,
				Status: "609",
			})
		}

		//รอดำเนินการ บันทึกเพิ่มอีก log ในส่วนของ NotificationTopup Status:"Success"
		db.Conn.Model(&model.LogTopup{}).Where("orderid = ?", request.Orderid).Where("data_type = ?", "NotificationTopup").Where("status = ?", "Wait").Update("status", "Success")

		// Send  200  Ok Success
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
		return
	}
}

// Redirect PayProcess
func (t *Topup) PayProcess(ctx *gin.Context) {

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
		fmt.Println("PayProcess: ", "Success")
		ctx.HTML(http.StatusOK, "frontend/topupdon.html", gin.H{
			"title":    "Age Of Khagan | Success.",
			"sum":      "Success",
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
