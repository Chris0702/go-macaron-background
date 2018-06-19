package database

import (
	// "encoding/base64"
	"fmt"
	// "image/gif"
	// "io/ioutil"
	// "os"
	"../util"
	"path"
	"path/filepath"
	"strings"
)

type Database interface {
	init(host string, port string, user string, password string, dbname string)
	initByConnectStr(connectStr string)
	createHTFolder(name string, orgId int64) bool
	createTable()
	saveSketchingElement(title string, data string, dataType string, orgId int64, folderName string, user string) bool
	saveSketchingBoardImage(name string, data string, dataType string, orgId int64, folderName string) bool
	getExploreMap(dataType string, orgId int64) map[string]interface{}
	getSketchingBoardImage(name string, dataType string, orgId int64) (bool, int, string, string, int, string)
	getSketchingElementTable(name string, dataType string, orgId int64) (bool, int, string, string, int, string, string, string, string)
	updateSketchingBoardFolderName(oldName string, newName string, dataType string, orgId int64) bool
	updateFileName(oldName string, newName string, dataType string, folderName string, orgId int64) bool
	deleteSketchingBoardFolder(name string, dataType string, orgId int64) bool
	deleteFileName(name string, dataType string, folderName string, orgId int64) bool
}

// id serial NOT NULL,
// 	version TEXT NOT NULL,
// 	title TEXT NOT NULL,
// 	data TEXT NOT NULL,
// 	data_type TEXT NOT NULL,
// 	org_id INT NOT NULL,
// 	folder_id INT NOT NULL,
// 	created_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
// 	,updated_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
// 	updated_by TEXT NOT NULL,created_by TEXT NOT NULL

var db Database
var orgTable = "org"

// var orgTable = "user"
var orgPreferenceTable = "org_preference"

// var userTable = "user1"
var userTable = "ht_user"
var orgUserTable = "org_user"
var sketchingElementTable = "sketching_element"
var sketchingBoardFolderTable = "sketching_board_folder"
var sketchingBoardImageTable = "sketching_board_image"
var version = "1.0.0"
var isWisePaas = false

func InitDB(dbType string, host string, port string, user string, password string, dbname string, isWisePaasIng bool) {
	isWisePaas = isWisePaasIng
	switch dbType {
	case "postgres":
		fmt.Println("postgres")
		db = postgresql{}
		db.init(host, port, user, password, dbname)
	case "sqlite":
		fmt.Println("sqlite")
		db = sqlite{}
		db.init(host, port, user, password, dbname)
	default:
		fmt.Println("default")
	}
}

func InitByConnectStr(dbType string, connectStr string) {
	switch dbType {
	case "postgres":
		fmt.Println("postgres")
		db = postgresql{}
		db.initByConnectStr(connectStr)
	case "sqlite":
		fmt.Println("sqlite")
		db = sqlite{}
		db.initByConnectStr(connectStr)
	default:
		fmt.Println("default")
	}
}

func CreateTable() {
	db.createTable()
	createDefaultFolder(1)
}

func CreateHTFolder(name string, orgId int64) bool {
	return db.createHTFolder(name, orgId)
}

func GetExploreMap(dataType string, orgId int64) map[string]interface{} {
	return db.getExploreMap(dataType, orgId)
}

func SaveSketchingBoardImage(title string, data string, orgId int64, folderName string) bool {
	dataType := getDataType(folderName)
	ext := strings.TrimLeft(filepath.Ext(title), ".")
	if ext == "svg"{
		imgName := 	filepath.Base(title)
		db.saveSketchingBoardImage(imgName, data, dataType, orgId, folderName)	

	}else{
		imgName, imgData := util.GetImgNameAndImgDataByImageDataB64(util.GetDirAndFilename(title), data)
		db.saveSketchingBoardImage(imgName, imgData, dataType, orgId, folderName)
	}
	return true
}

func SaveSketchingElement(title string, data string, orgId int64, folderName string, user string) bool {
	dataType := getDataType(folderName)
	return db.saveSketchingElement(title, data, dataType, orgId, folderName, user)
}

func GetSketchingElementTable(name string, dataType string, orgId int64) (bool, int, string, string, int, string, string, string, string) {
	return db.getSketchingElementTable(name, dataType, orgId)
}

func GetSketchingBoardImage(name string, dataType string, orgId int64) (bool, int, string, string, int, string) {
	return db.getSketchingBoardImage(name, dataType, orgId)
}

func Rename(oldName string, newName string, orgId int64) bool {
	if util.IsFile(oldName) {
		dataType := getDataType(oldName)
		folderName := path.Dir(oldName)
		return db.updateFileName(path.Base(oldName), path.Base(newName), dataType, folderName, orgId)
	} else {
		dataType := getDataType(oldName)
		return db.updateSketchingBoardFolderName(oldName, newName, dataType, orgId)
	}
}

func Remove(name string, orgId int64) bool {
	if util.IsFile(name) {
		dataType := getDataType(name)
		folderName := path.Dir(name)
		return db.deleteFileName(path.Base(name), dataType, folderName, orgId)
	} else {
		dataType := getDataType(name)
		return db.deleteSketchingBoardFolder(name, dataType, orgId)
	}
}

func createDefaultFolder(orgId int64) {
	db.createHTFolder("assets", orgId)
	db.createHTFolder("components", orgId)
	db.createHTFolder("displays", orgId)
	db.createHTFolder("symbols", orgId)
}

func getDataType(folderName string) string {
	var result = ""
	for i := 0; i < len(folderName); i++ {
		if string(folderName[i]) == "/" {
			i = len(folderName)
		} else {
			result = result + string(folderName[i])
		}
	}
	return result
}

func filterDataTypePath(dataPath string) string {
	var result = ""
	filter := true
	for i := 0; i < len(dataPath); i++ {
		if filter == false {
			result = result + string(dataPath[i])
		}
		if string(dataPath[i]) == "/" {
			filter = false
		}
	}
	return result
}

func checkErr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	} else {
		return false
	}
}
