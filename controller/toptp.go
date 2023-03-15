package controller

import (
	"ccrd/dto"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Topup struct{}

func Paytopups(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "frontend/topup.html", gin.H{
		"title": "Age Of Khagan Thailand | Topup เติมเงิน",
		"css":   "topup.css",
	})
}
func (t *Topup) Paytopup(ctx *gin.Context) {

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

	//
	// Code...
	//
	//------------------------------------------------------------------------------
	// -- แบบที่ 1
	// var channel = request.Channel
	// var forr = ""
	// var operator = ""
	// var orderid = request.Orderid
	// var price = request.Amount
	// var sid = os.Getenv("SID")
	// var uid = ""
	// var SECRET_KEY = os.Getenv("SECRET_KEY")
	// data := channel + forr + operator + orderid + price + sid + uid + SECRET_KEY

	//------------------------------------------------------------------------------
	// -- แบบที่ 2
	var channel = request.Channel
	var forr = os.Getenv("FOR") + "-" + request.Orderid
	var operator = ""
	var orderid = request.Orderid
	var price = request.Amount + request.Currency
	var sid = os.Getenv("SID")
	var uid = ""
	var SECRET_KEY = os.Getenv("SECRET_KEY")
	data := channel + forr + operator + orderid + price + sid + uid + SECRET_KEY
	//------------------------------------------------------------------------------

	h := md5.New()
	io.WriteString(h, data)
	sumSig := hex.EncodeToString(h.Sum(nil))

	if request.Sig == sumSig {
		//
		//Code..
		//
		ctx.JSON(http.StatusOK, dto.TopupResponse{
			Txid:   request.Txid,
			Status: "200",
		})
	} else {
		fmt.Println("data", data)
		fmt.Println("old", request.Sig)
		fmt.Println("new", sumSig)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "issue Sig",
		})
	}

}
