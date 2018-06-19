package database

import (
	"fmt"
)

type sqlite struct{}

func (m sqlite) init(host string, port string, user string, password string, dbname string) {
	fmt.Println("sqlite init")
}

func (m sqlite) initByConnectStr(connectStr string) {
	fmt.Println("sqlite initByConnectStr-----")
	fmt.Println(connectStr)
	fmt.Println("----------------------------")
}

func (m sqlite) createHTFolder(name string, orgId int64) bool {
	return true
}

func (m sqlite) saveSketchingElement(title string,data string,dataType string, orgId int64,folderName string,user string) bool {
	return true
}

func (m sqlite) saveSketchingBoardImage(name string,data string,dataType string, orgId int64,folderName string) bool {
	return true
}

func (m sqlite) createTable() {

}

func (m sqlite) getExploreMap(dataType string,orgId int64) map[string]interface{} {
	a := make(map[string]interface{})
	return a
}

func (m sqlite) getSketchingBoardImage(name string,dataType string, orgId int64)(bool, int, string, string, int, string) {
	return false,0,"","",0,""
}

func (m sqlite) getSketchingElementTable(name string,dataType string, orgId int64)(returnErr bool, id int, version string,data string,folderId int,createdTime string,updatedTime string,updatedBy string,createdBy string) {
	return 
}

func (m sqlite) updateSketchingBoardFolderName(oldName string,newName string,dataType string, orgId int64) bool {
	return true
}

func (m sqlite) updateFileName(oldName string,newName string,dataType string,folderName string, orgId int64) bool {
	return true
}

func (m sqlite) deleteSketchingBoardFolder(name string,dataType string, orgId int64)(bool){
	return true
}

func (m sqlite) deleteFileName(name string,dataType string,folderName string, orgId int64) bool {
	return true
}

// func (m sqlite) saveWidgetbuilderWidget(name string, imageDataB64 string, contentJson string, orgId int64) bool {
// 	return true
// }

// func (m sqlite) saveComponent(fileName string, imageDataB64 string, combJsonString string, orgID int64) (ok bool) {
// 	return true
// }

// func (m sqlite) removeComponent(fileName string, orgID int64) (ok bool) {
// 	return true
// }

// func (m sqlite) getWidgetbuilderTree(orgId int64) []map[string]interface{} {
// 	a := make([]map[string]interface{}, 5)
// 	return a
// }

// func (m sqlite) getComponentTree(orgId int64) []map[string]interface{} {
// 	a := make([]map[string]interface{}, 5)
// 	return a
// }

// func (m sqlite) removeWidgetbuilderWidget(name string, orgId int64) bool {
// 	return true
// }

// func (m sqlite) removeWidgetbuilderFolder(name string, orgId int64) bool {
// 	return true
// }

// func (m sqlite) isExistWidgetbuilderWidget(name string, orgId int64) bool {
// 	return true
// }

// func (m sqlite) renameWidgetbuilder(oldFilePath string, newName string, orgId int64) bool {
// 	return true
// }

// func (m sqlite) upload(name string, value string, folderPath string, orgID int64) (ok bool) {
// 	return true
// }

// func (m sqlite) getImgByNamePathOrgid(name string, folderPath string, orgID int64) (img string, ok bool) {
// 	return "", false
// }

// func (m sqlite) showImagesPath(path string, orgID int64, recursive bool, constraints string) (lists []string, err error) {
// 	return nil, nil
// }

// func (m sqlite) getGifImagedetail(fileName string, orgId int64) (bool, map[string]interface{}) {
// 	return true, make(map[string]interface{})
// }

// func (m sqlite) getImagesInfoByFolderPath(folderPath string, orgId int64) (bool, []map[string]interface{}) {
// 	return true, make([]map[string]interface{}, 0)
// }

// func (m sqlite) insertDefaultImagesByOrgID(orgID int64) (ok bool) {
// 	return true
// }

// func (m sqlite) deleteDefaultImagesByOrgID(orgID int64) (ok bool) {
// 	return true
// }

// func (m sqlite) getWidgetbuilderContentJson(name string, orgId int64) (bool, string) {
// 	return false, ""
// }

// func (m sqlite) isExistImage(name string, orgId int64, folderPath string) bool {
// 	return false
// }

// func (m sqlite) getComponentContentJson(name string, orgId int64) (bool, string) {
// 	return false, ""
// }

// func (m sqlite) getOrgNameByOrgId(orgId int64) string {
// 	return ""
// }

// func (m sqlite) getWidgetbuilderList(orgId int64) (bool, []string) {
// 	a := make([]string, 0)
// 	return false, a
// }
