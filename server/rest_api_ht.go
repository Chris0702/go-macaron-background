package server

import "fmt"
import "net/http"
import "gopkg.in/macaron.v1"
import db "../database"
import "path"
import "../util"
import "path/filepath"

func staticApi(w http.ResponseWriter, dataType string, name string) {
	if util.IsImg(name) {
		err, _, _, value, _, _ := db.GetSketchingBoardImage(name, dataType, 1)
		if err {
			sendMessage(w, "err")
		} else {
			sendImg(w, filepath.Base(name), value)
		}
	} else {
		err, _, _, value, _, _, _, _, _ := db.GetSketchingElementTable(name, dataType, 1)
		if err {
			sendMessage(w, "err")
		} else {
			sendMessage(w, value)
		}
	}
}

func displays(ctx *macaron.Context, w http.ResponseWriter, r *http.Request) {
	name := ctx.Params("*")
	staticApi(w, "displays", name)
}

func symbols(ctx *macaron.Context, w http.ResponseWriter, r *http.Request) {
	name := ctx.Params("*")
	staticApi(w, "symbols", name)
}

func components(ctx *macaron.Context, w http.ResponseWriter, r *http.Request) {
	name := ctx.Params("*")
	staticApi(w, "components", name)
}

func assets(ctx *macaron.Context, w http.ResponseWriter, r *http.Request) {
	name := ctx.Params("*")
	staticApi(w, "assets", name)
}

func explore(ctx *macaron.Context, w http.ResponseWriter, r *http.Request) {
	exploreType := ctx.Params("*")
	fmt.Println("----------explore----exploreType--------")
	fmt.Println(exploreType)
	res := db.GetExploreMap(exploreType, 1)
	JsonRes(w, res)
}

func upload(ctx *macaron.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------upload------------")
	var getPath = r.PostFormValue("path")
	var content = r.PostFormValue("content")
	var getExtName = path.Ext(getPath)
	var saveName = path.Base(getPath)
	var folder = path.Dir(getPath)
	var result = true
	var user = "user"
	var orgId int64	
	orgId = 1
	if util.IsImg(getExtName){
		result = db.SaveSketchingBoardImage(saveName, content, orgId, folder)
	} else {
		result = db.SaveSketchingElement(saveName, content, orgId, folder, user)
	}
	if result {
		sendMessage(w, "true")
	} else {
		sendMessage(w, "false")
	}
}

func mkdir(ctx *macaron.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------mkdir------------")
	res := db.CreateHTFolder(r.PostFormValue("path"), 1)
	if res {
		sendMessage(w, "true")
	} else {
		sendMessage(w, "false")
	}
}

func rename(ctx *macaron.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------rename------------")
	var oldName = r.PostFormValue("oldPath")
	var newName = r.PostFormValue("newPath")
	result := db.Rename(oldName, newName, 1)
	res := fmt.Sprintf("%v", result)
	sendMessage(w, res)
}

func remove(ctx *macaron.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------remove------------")
	var name = r.PostFormValue("path")
	result := db.Remove(name, 1)
	res := fmt.Sprintf("%v", result)
	sendMessage(w, res)
}
