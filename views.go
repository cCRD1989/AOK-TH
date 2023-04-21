package main

import (
	"ccrd/server/khanscr"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
)

func getPlayersOnlineCount() string {
	return "199 คน"
}
func getAllItemsBSV() map[int16]khanscr.BsvItem {
	return khanscr.GetAllItems()
}
func getItemNameById(id int16) string {
	return khanscr.GetItemName(id)
}

func getNameMonter(names string) string {
	data := strings.Split(names, "/")
	count := len(data)
	return strings.Split(data[count-1], ".")[0]
}

func createViews() multitemplate.Render {
	var fn = template.FuncMap{
		"getPlayersOnlineCount": getPlayersOnlineCount,
		"getAllItemsBSV":        getAllItemsBSV,
		"getItemNameById":       getItemNameById,
		"getNameMonter":         getNameMonter,
	}
	var r = multitemplate.New()
	var vtpath = filepath.Join("views", "templates")
	var dirs, err = ioutil.ReadDir("views/layouts/")
	checkAndPanic(err)
	for _, dir := range dirs {
		var dirName = dir.Name()
		layouts, err := filepath.Glob(fmt.Sprintf("views/layouts/%s/*.html", dirName))
		checkAndPanic(err)

		var templates = []string{}
		err = filepath.Walk(fmt.Sprintf("views/templates/%s/", dirName), func(path string, f os.FileInfo, err error) error {
			if filepath.Ext(path) == ".html" {
				templates = append(templates, path)
			}
			return nil
		})
		checkAndPanic(err)
		for _, tmpl := range templates {
			var tname = strings.Replace(tmpl, vtpath, "", 1)  // ลบพาทออก
			tname = strings.Replace(tname, "\\", "/", -1)[1:] //เปลี่ยนให้เป็นรูท
			log.Printf("[GIN-debug] %-6s %-25s --> %s\n", "VIEW", dirName, tname)
			r.AddFromFilesFuncs(tname, fn, append(layouts, tmpl)...)
			//r.AddFromFiles(tname, append(layouts, tmpl)...)
		}
	}
	return r
}

func checkAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}
