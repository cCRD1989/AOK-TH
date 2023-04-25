package controller

import (
	"ccrd/model/aokmodel"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Guides struct{}

func (g *Guides) GetMap(ctx *gin.Context) {

	var usr, _ = ctx.Get("user")
	user, _ := usr.(aokmodel.Userlogin)

	maps := ctx.Param("maps")

	mapName := ""
	mapMonter := []string{}
	mapid := ""
	var mapDes []string

	if maps == "wroclawfortress" {
		mapName = "wroclaw fortress"
		mapMonter = mapPart("1")
		mapid = "/public/data/Maps/1.png"

		mapDes = []string{"ชนเผ่า Durlukin", "จุดยุทธศาสตร์แห่งการรวมตัวของชนเผ่า Durlukin", "เพื่อเตรียมความพร้อมในการโจมตีกับ", "กองทัพ Nurin"}
	} else if maps == "kharakorum" {
		mapName = "kharakorum"
		mapMonter = mapPart("10")
		mapid = "/public/data/Maps/10.png"
		mapDes = []string{"ป้อมปราการเอก", "พื้นที่สำหรับนักรบในการต่อต้านเหล่ามอนสเตอร์", "ที่แข็งแกร่งและชั่วร้าย"}
	} else if maps == "lublinmongolfortress" {
		mapName = "lublin mongol fortress"
		mapMonter = mapPart("2")
		mapid = "/public/data/Maps/2.png"
		mapDes = []string{"เหมืองแร่ใต้หุบเขา Karpatian", "เหมืองแร่โบราณแห่งความท้าทาย", "กับสภาพของผู้คนที่เปลี่ยนไป", "ด้วยความโลภและเวทย์มนต์ดำ"}
	} else if maps == "irondungeon" {
		mapName = "iron dungeon"
		mapMonter = mapPart("3")
		mapid = "/public/data/Maps/3.png"
		mapDes = []string{"ลาวาใต้พิภพ", "สถานที่น่าค้นหาและมีเสน่ห์", "รายล้อมไปด้วยมอนสเตอร์", "ผู้ปกป้องทรัพย์สมบัติล้ำค่า"}
	} else if maps == "lavacanyon" {
		mapName = "lava canyon"
		mapMonter = mapPart("4")
		mapid = "/public/data/Maps/4.png"
		mapDes = []string{"ชนเผ่า Nurin", "จุดยุทธศาสตร์แห่งการรวมตัวของชนเผ่า Nurin", "เพื่อเตรียมความพร้อมในการโจมตีกับ", "กองทัพ Durlukin"}
	} else {
		ctx.Redirect(http.StatusOK, "/")
		return
	}

	ctx.HTML(http.StatusOK, "frontend/maps.html", gin.H{
		"title":     "Age Of Khagan | " + mapName,
		"user":      user,
		"mapName":   mapName,
		"mapMonter": mapMonter,
		"mapid":     mapid,
		"mapDes":    mapDes,
	})

}

func mapPart(p string) []string {

	var pathall = []string{}
	if err := filepath.Walk(fmt.Sprintf("public/data/Maps/%s/", p), func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ".png" {
			path = "/" + strings.Replace(path, "\\", "/", -1)
			pathall = append(pathall, path)
		}
		return nil
	}); err != nil {
		log.Fatal("filepath.Walk", err.Error())
	}

	return pathall
}
