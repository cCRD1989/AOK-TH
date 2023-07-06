package controller

import (
	"ccrd/db"
	"ccrd/dto"
	"ccrd/model"
	"ccrd/model/aokmodel"
	"ccrd/unit"
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

// เปิดหน้าเติมเงิน
func Paytopups(ctx *gin.Context) {

	visit := model.LogWeb{
		DataType:  "visit",
		IPAddress: ctx.ClientIP(),
	}
	db.Conn.Save(&visit)

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	if user.Username != "" {
		ctx.HTML(http.StatusOK, "frontend/topup.html", gin.H{
			"title": "Age Of Khagan Thailand | เติมเงิน",
			"user":  user,
			"ff":    "",
			"bg":    "/public/data/img/TOPUP_BG.png",
		})
	} else {
		ctx.HTML(http.StatusOK, "frontend/login.html", gin.H{
			"title": "Age Of Khagan Thailand | Login",
			"user":  user,
			"bg":    "/public/data/img/TOPUP_BG.png",
		})
	}
}

func GetBonusBanking(ctx *gin.Context) {

	channel := ctx.DefaultPostForm("channel", "")

	errs := unit.Validate(map[string]interface{}{
		"channel": channel,
	}, map[string]string{

		"channel": "required|alphanum",
	})
	if errs != nil {
		return
	}

	Bonus := model.Bankingbonus{}
	if err := db.Conn.Where("channel =?", channel).First(&Bonus).Error; err != nil {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"Bonus":  strconv.Itoa(Bonus.Bonus),
	})
}

// ทำรายการ
func PaytopupsAddPoint(ctx *gin.Context) {

	usernameId := ctx.PostForm("username")
	price := ctx.PostForm("price") + "THB"
	channel := ctx.PostForm("channel")

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	errs := unit.Validate(map[string]interface{}{
		"username": usernameId,
		"price":    price,
		"channel":  channel,
	}, map[string]string{

		"username": "required|alphanum",
		"price":    "required|alphanum",
		"channel":  "required|alphanum",
	})
	if errs != nil {
		ctx.HTML(http.StatusOK, "frontend/topup.html", gin.H{
			"title":   "Age Of Khagan | เติมเงิน",
			"message": "ข้อมูลไม่ถูกต้อง",
			"ff":      "",
			"user":    user,
			"bg":      "/public/data/img/TOPUP_BG.png",
		})
		return
	}

	if usernameId == user.Username {

		if usernameId == "" || channel == "" || price == "" {
			ctx.HTML(http.StatusOK, "frontend/topup.html", gin.H{
				"title":   "Age Of Khagan | เติมเงิน",
				"message": "ข้อมูลไม่ถูกต้อง",
				"ff":      "",
				"user":    user,
				"bg":      "/public/data/img/TOPUP_BG.png",
			})
			return
		}

		topup := aokmodel.Userlogin{}
		if err := db.AOK_DB.Where("username =?", usernameId).First(&topup).Error; err != nil {
			ctx.HTML(http.StatusOK, "frontend/topup.html", gin.H{
				"title":   "Age Of Khagan | Topup เติมเงิน",
				"message": "ไอดีไม่มีอยู่ในระบบ โปรดลองใหม่อีกครั้ง",
				"ff":      "",
				"user":    user,
				"bg":      "/public/data/img/TOPUP_BG.png",
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
			ctx.HTML(http.StatusOK, "frontend/topup.html", gin.H{
				"title":   "Age Of Khagan | Topup เติมเงิน",
				"message": "RequestTopup Error",
				"ff":      "",
				"user":    user,
				"bg":      "/public/data/img/TOPUP_BG.png",
			})
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

		//ctx.Redirect(http.StatusTemporaryRedirect, url)
		ctx.Redirect(http.StatusFound, urlA.String())

	} else {
		ctx.Redirect(http.StatusFound, "/")
	}
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

// ลูกค้ากด ออเดอร์ ออกไป ให้ Razer
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
	fmt.Println("notification data all :", request, ">>>>>>>>>>>>", request.Detail)

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
			ctx.Status(http.StatusBadRequest)
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
			Bonus:     0,
			Sig:       request.Sig,
			IPAddress: ctx.ClientIP(),
		})

		if request.Status == "200" {

			// รายการเงินที่จะเติม ตามที่ลูกค้ากดเข้ามา
			caseint, err := strconv.Atoi(request.Amount)
			if err != nil {
				fmt.Println("str to int ไม่ได้ ", data.Price)
				ctx.Status(http.StatusBadRequest)
				return
			}

			//Bonus โปรแกรมเติมเงิน ที่มีโปบนัสเพิ่ม %

			BonusTopup := model.Bankingbonus{}
			if err := db.Conn.Where("Channel = ?", request.Channel).First(&BonusTopup).Error; err != nil {
				fmt.Println("ค้นหาโบนัสไม่เจอ", err.Error())
				ctx.Status(http.StatusBadRequest)
				return
			}

			//ดึงเงินที่อยู่ใน id นั้น
			idcash := aokmodel.Userlogin{}
			db.AOK_DB.First(&idcash, "username = ?", data.UserId)

			CASH := caseint + int(float64(caseint)*(float64(BonusTopup.Bonus)/float64(100)))

			fmt.Println(">>>>>>>>>>>>>>>>>>", BonusTopup)
			fmt.Println(">>>>>>>>>>>>>>>>>>", request.Channel)
			fmt.Println(">>>>>>>>>>>>>>>>>>", BonusTopup.Bonus)
			fmt.Println(">>>>>>>>>>>>>>>>>>", CASH)

			log_cash := model.LogMailTopup{
				Eventid:    "9",
				Senderid:   idcash.Id,
				Sendername: "SYSTEM",
				Receiverid: idcash.Id,
				Title:      "Bonus Pre Topup",
				Content:    "คุณได้รับ (Bonus Pre Topup) จำนวน " + strconv.Itoa(CASH) + " CASH",
				Gold:       0,
				Cash:       CASH,
				Currencies: " ",
				Items:      " ",
			}

			if err := db.AOK_DB.Save(&log_cash).Error; err != nil {
				fmt.Println("AOKบันทึกแคชไม่สำเร็จ", err.Error())
				ctx.Status(http.StatusBadRequest)
				return
			}

			if err := db.Conn.Save(&log_cash).Error; err != nil {
				fmt.Println("LOGบันทึกแคชไม่สำเร็จ", err.Error())
				ctx.Status(http.StatusBadRequest)
				return
			}

			db.Conn.Save(&model.Topuprecheck{
				UserId:    data.UserId,
				Txid:      request.Txid,
				Orderid:   request.Orderid,
				Status:    "Done",
				Detail:    request.Detail,
				Channel:   request.Channel,
				Price:     request.Amount,
				Bonus:     strconv.Itoa(BonusTopup.Bonus),
				Sig:       request.Sig,
				IPAddress: ctx.ClientIP(),
			})

			//ของเก่า
			// if err := db.AOK_DB.Model(&aokmodel.Userlogin{}).Where("username = ?", idcash.Username).Update("cash", idcash.Cash).Error; err != nil {
			// 	fmt.Println("บันทึกแคชไม่สำเร็จ", err.Error())
			// 	ctx.Status(http.StatusBadRequest)
			// 	return
			// }

			//รอดำเนินการ บันทึกเพิ่มอีก log ในส่วนของ NotificationTopup Status:"Success"
			db.Conn.Model(&model.LogTopup{}).Where("orderid = ?", request.Orderid).Where("data_type = ?", "NotificationTopup").Where("status = ?", "Wait").Update("status", "Success")
			db.Conn.Model(&model.LogTopup{}).Where("orderid = ?", request.Orderid).Where("data_type = ?", "NotificationTopup").Where("status = ?", "Success").Update("bonus", strconv.Itoa(BonusTopup.Bonus))

			//
			// Send  200  Ok Success
			ctx.JSON(http.StatusOK, dto.TopupResponse{
				Txid:   request.Txid,
				Status: "200",
			})
		} else {

			//รอดำเนินการ บันทึกเพิ่มอีก log ในส่วนของ NotificationTopup Status:"Success"
			db.Conn.Model(&model.LogTopup{}).Where("orderid = ?", request.Orderid).Where("data_type = ?", "NotificationTopup").Where("status = ?", "Wait").Update("status", "Failed")

			ctx.JSON(http.StatusOK, dto.TopupResponse{
				Txid:   request.Txid,
				Status: "400",
			})

		}

	} else {
		fmt.Println("Sig ไม่ตรง")
		fmt.Println("data", data)
		fmt.Println("old", request.Sig)
		fmt.Println("new", sumSig)
		ctx.Status(http.StatusBadRequest)
		return
	}
}

// Redirect PayProcess จาก Razer รายงานผลการเติมเงิน
func (t *Topup) PayProcess(ctx *gin.Context) {

	// ตรวจสอบ User Cookie
	usr, _ := ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

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
	//a10a5bc6f1b8e430fd5eb2fde1773d4f a729ff0bbb3e14e9c5f97f92a91db667 200 Success bbl 100 THB 9f042e5b2be0147160d5533fe32ee38a
	fmt.Println("Redirect PayProcess data all: ", request)

	if request.Status == "200" {
		fmt.Println("PayProcess: ", "Success")

		ctx.HTML(http.StatusOK, "frontend/topup.html", gin.H{
			"title":    "Age Of Khagan | Success.",
			"user":     user,
			"txid":     request.Txid,
			"orderid":  request.Orderid,
			"status":   request.Status,
			"detail":   request.Detail,
			"channel":  request.Channel,
			"amount":   request.Amount,
			"currency": request.Currency,
			"sig":      request.Sig,
			"ff":       "ok",
			"bg":       "/public/data/img/TOPUP_BG.png",
		})
	} else {
		fmt.Println("PayProcess: ", "Failed")
		ctx.HTML(http.StatusOK, "frontend/topup.html", gin.H{
			"title":    "Age Of Khagan | Failed.",
			"user":     user,
			"txid":     "Failed",
			"orderid":  "Failed",
			"status":   "Failed",
			"detail":   "Failed",
			"channel":  "Failed",
			"amount":   "Failed",
			"currency": "Failed",
			"sig":      "Failed",
			"ff":       "nook",
			"bg":       "/public/data/img/TOPUP_BG.png",
		})
	}
}
