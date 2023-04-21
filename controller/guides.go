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

	if maps == "wroclawfortress" {
		mapName = "wroclaw fortress"
		mapMonter = mapPart("1")
		mapid = "/public/data/Maps/1.png"
	} else if maps == "kharakorum" {
		mapName = "kharakorum"
		mapMonter = mapPart("10")
		mapid = "/public/data/Maps/10.png"
	} else if maps == "lublinmongolfortress" {
		mapName = "lublin mongol fortress"
		mapMonter = mapPart("2")
		mapid = "/public/data/Maps/2.png"
	} else if maps == "irondungeon" {
		mapName = "iron dungeon"
		mapMonter = mapPart("3")
		mapid = "/public/data/Maps/3.png"
	} else if maps == "lavacanyon" {
		mapName = "lava canyon"
		mapMonter = mapPart("4")
		mapid = "/public/data/Maps/4.png"
	} else {
		ctx.Redirect(http.StatusOK, "/")
		return
	}

	ctx.HTML(http.StatusOK, "frontend/guides.html", gin.H{
		"title":     "Age Of Khagan | " + mapName,
		"user":      user,
		"mapName":   mapName,
		"mapMonter": mapMonter,
		"mapid":     mapid,
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
