package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	// "bufio"
	_ "github.com/lib/pq"
	// "io/ioutil"
	"strconv"
	// "os"
	// "encoding/base64"
	"../util"
	"path/filepath"
	// "github.com/grafana/grafana/pkg/setting"
)

var postgresqlDB *sql.DB
var schemaName = "composer"
var ownerName = "g_composer"
var username = ""

type postgresql struct{}

func (m postgresql) init(host string, port string, user string, password string, dbname string) {
	fmt.Println("postgresql init")
	fmt.Println(host)
	fmt.Println(port)
	fmt.Println(user)
	fmt.Println(password)
	fmt.Println(dbname)
	fmt.Println(postgresqlDB)
	username = user
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	postgresqlDB = db
	fmt.Println(reflect.TypeOf(db))
	fmt.Println(postgresqlDB)
}

func (m postgresql) initByConnectStr(connectStr string) {
	fmt.Println("postgresql init connectStr")
	fmt.Println(connectStr)
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		panic(err)
	}
	postgresqlDB = db
	fmt.Println(reflect.TypeOf(db))
	fmt.Println(postgresqlDB)
}

func (m postgresql) createTable() {
	// fmt.Println("postgresql init connectStr")
	// fmt.Println(connectStr)
	// create table if not exists
	if isWisePaas {
		createSchema()
	}
	fmt.Println("createTable")
	sql_table := "CREATE TABLE IF NOT EXISTS " + orgTable + "(id serial NOT NULL,version TEXT NOT NULL,name TEXT NOT NULL ,PRIMARY KEY (id));"
	_, err := postgresqlDB.Exec(sql_table)

	if err != nil {
		panic(err)
	}

	sql_table = "CREATE TABLE IF NOT EXISTS " + orgPreferenceTable + "(id serial NOT NULL,version TEXT NOT NULL,org_id INT NOT NULL,data TEXT NOT NULL,updated_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() ,PRIMARY KEY (id));"
	_, err = postgresqlDB.Exec(sql_table)
	if err != nil {
		panic(err)
	}

	sql_table = "CREATE TABLE IF NOT EXISTS " + userTable + "(id serial NOT NULL,version TEXT NOT NULL,email TEXT NOT NULL,name TEXT NOT NULL,password TEXT NOT NULL,org_id_default INT NOT NULL,is_admin BOOL NOT NULL,created_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),updated_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() ,PRIMARY KEY (id));"
	_, err = postgresqlDB.Exec(sql_table)
	if err != nil {

		fmt.Println(err)

		panic(err)
	}

	sql_table = "CREATE TABLE IF NOT EXISTS " + orgUserTable + "(id serial NOT NULL,org_id INT NOT NULL,ht_user_id INT NOT NULL,role TEXT NOT NULL,created_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),updated_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() ,PRIMARY KEY (id));"
	_, err = postgresqlDB.Exec(sql_table)
	if err != nil {

		fmt.Println(err)

		panic(err)
	}

	sql_table = "CREATE TABLE IF NOT EXISTS " + sketchingElementTable + "(id serial NOT NULL,version TEXT NOT NULL,title TEXT NOT NULL,data TEXT NOT NULL,data_type TEXT NOT NULL,org_id INT NOT NULL,folder_id INT NOT NULL,created_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),updated_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),updated_by TEXT NOT NULL,created_by TEXT NOT NULL ,PRIMARY KEY (title,folder_id,org_id));"
	_, err = postgresqlDB.Exec(sql_table)
	if err != nil {

		fmt.Println(err)

		panic(err)
	}

	sql_table = "CREATE TABLE IF NOT EXISTS " + sketchingBoardFolderTable + "(id serial NOT NULL,version TEXT NOT NULL,name TEXT NOT NULL,data_type TEXT NOT NULL,org_id INT NOT NULL,created_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() ,PRIMARY KEY (name,org_id));"
	_, err = postgresqlDB.Exec(sql_table)
	if err != nil {

		fmt.Println(err)

		panic(err)
	}

	sql_table = "CREATE TABLE IF NOT EXISTS " + sketchingBoardImageTable + "(id serial NOT NULL,version TEXT NOT NULL,name TEXT NOT NULL,data BYTEA NOT NULL,data_type TEXT NOT NULL,org_id INT NOT NULL,folder_id INT NOT NULL,updated_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() ,PRIMARY KEY (name,folder_id,org_id));"
	_, err = postgresqlDB.Exec(sql_table)
	if err != nil {

		fmt.Println(err)

		panic(err)
	}
	fmt.Println("createTable ok")
}

func (m postgresql) createHTFolder(name string, orgId int64) bool {
	// name = "/" + name
	// name = filepath.FromSlash(name) // 平台处理
	err, isExist, _ := isExistHTFolder(name, orgId)
	if err {
		return false
	} else {
		if isExist {
			return false
		} else {
			// return insertWidgetbuilderFolderPostgresql(name, orgId)
			_, insertFolderResult := insertHTFolderPostgresql(name, orgId)
			return insertFolderResult
		}
	}
}

func (m postgresql) saveSketchingElement(title string, data string, dataType string, orgId int64, folderName string, user string) bool {
	err, folderId, _, _ := queryFolderByNameAndOrgId(folderName, orgId)
	if err == true {
		return false
	}
	_, result := saveSketchingElementExe(title, data, dataType, orgId, folderId, user)
	return result
}

func (m postgresql) saveSketchingBoardImage(name string, data string, dataType string, orgId int64, folderName string) bool {
	err, folderId, _, _ := queryFolderByNameAndOrgId(folderName, orgId)
	fmt.Println(folderId)
	if err == true {
		return false
	}
	result := saveSketchingBoardImageExe(name, data, dataType, orgId, folderId)
	return result
}

func (m postgresql) getExploreMap(dataType string, orgId int64) map[string]interface{} {
	folderRes, folders := querySketchingBoardFolderByDataTypeAndOrgId(dataType, orgId)
	if folderRes {
		elementRes, elements := querySketchingElementByDataTypeAndOrgId(dataType, orgId)
		if elementRes {
			imgRes, imgs := querySketchingBoardImageByDataTypeAndOrgId(dataType, orgId)
			if imgRes {
				return mergeExplore(folders, elements, imgs)
			}
		}
	}
	return make(map[string]interface{})
}

func (m postgresql) getSketchingElementTable(name string, dataType string, orgId int64) (returnErr bool, id int, version string, data string, folderId int, createdTime string, updatedTime string, updatedBy string, createdBy string) {
	folderName := dataType
	if filepath.Dir(name) != "." {
		folderName = filepath.Join(folderName,filepath.Dir(name))
	}
	queryName := filepath.Base(name)
	err, isExist, existId := isExistHTFolder(folderName, orgId)
	if err {
		return true, 0, "", "", 0, "", "", "", ""
	} else {
		folderId := fmt.Sprintf("%v", existId)
		if isExist {
			return querySketchingElement(queryName, dataType, orgId, folderId)
		}
	}
	return true, 0, "", "", 0, "", "", "", ""
}

func (m postgresql) updateSketchingBoardFolderName(oldName string, newName string, dataType string, orgId int64) bool {
	stmt := fmt.Sprintf("UPDATE " + sketchingBoardFolderTable + " SET name=$1" + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND name = '" + oldName + "'" + " AND data_type = '" + dataType + "'")
	sqlcmd, err := postgresqlDB.Prepare(stmt)
	if err != nil {
		return false
	}
	_, err = sqlcmd.Exec(newName)
	if err != nil {
		return false
	}
	return true
}

func (m postgresql) updateFileName(oldName string, newName string, dataType string, folderName string, orgId int64) bool {
	err, folderId, _, _ := querySketchingFolderByNameAndOrgIdAndDataType(folderName, dataType, orgId)
	if err || folderId <= 0 {
		return false
	}
	folderIdReq := fmt.Sprintf("%v", folderId)
	if util.IsImg(oldName) {
		return updateSketchingBoardImageName(oldName, newName, dataType, orgId, folderIdReq)
	}
	return updateSketchingElementName(oldName, newName, dataType, orgId, folderIdReq)
}

func (m postgresql) getSketchingBoardImage(name string, dataType string, orgId int64) (returnErr bool, id int, version string, data string, folderId int, updatedTime string) {
	folderName := dataType
	if filepath.Dir(name) != "." {
		folderName = filepath.Join(folderName,filepath.Dir(name))
	}
	queryName := filepath.Base(name)
	err, isExist, existId := isExistHTFolder(folderName, orgId)
	if err {
		return true, 0, "", "", 0, ""
	} else {
		folderId := fmt.Sprintf("%v", existId)
		if isExist {
			return querySketchingBoardImage(queryName, dataType, orgId, folderId)
		}
	}
	return true, 0, "", "", 0, ""
}

func (m postgresql) deleteSketchingBoardFolder(name string, dataType string, orgId int64) bool {
	err, isExist, folderId := isExistHTFolder(name, orgId)
	if err || !isExist {
		return false
	}
	deleteRes := deleteSketchingBoardFolderByIDPostgresql(folderId)
	if !deleteRes {
		return false
	}
	deleteRes = deleteSketchingElementByFolderIdPostgresql(folderId)
	if !deleteRes {
		return false
	}
	deleteRes = deleteSketchingBoardImageTableByFolderIdPostgresql(folderId)
	if !deleteRes {
		return false
	}
	return true
}

func (m postgresql) deleteFileName(name string, dataType string, folderName string, orgId int64) bool {
	err, folderId, _, _ := querySketchingFolderByNameAndOrgIdAndDataType(folderName, dataType, orgId)
	if err || folderId <= 0 {
		return false
	}
	if util.IsImg(name) {
		return deleteSketchingBoardImageTableByTitleFolderIdOrgIdPostgresql(name, folderId, orgId)
	}
	return deleteSketchingElementByTitleFolderIdOrgIdPostgresql(name, folderId, orgId)
}

func createSchema() {
	fmt.Println("createSchema")
	sql_table := "CREATE SCHEMA IF NOT EXISTS \"" + schemaName + "\" AUTHORIZATION " + ownerName + ";ALTER ROLE \"" + username + "\" SET search_path TO \"" + schemaName + "\";"
	_, err := postgresqlDB.Exec(sql_table)

	if err != nil {
		fmt.Println(sql_table)
		panic(err)
	}
}

func deleteSketchingElementByTitleFolderIdOrgIdPostgresql(name string, folderId int, orgId int64) (ok bool) {
	stmt := fmt.Sprintf("DELETE FROM %v WHERE folder_id=$1 AND title =$2 AND org_id =$3 ", sketchingElementTable)
	sqlcmd, err := postgresqlDB.Prepare(stmt)
	if err != nil {
		return false
	}
	_, err = sqlcmd.Exec(folderId, name, orgId)
	if err != nil {
		return false
	}
	return true
}

func deleteSketchingBoardImageTableByTitleFolderIdOrgIdPostgresql(name string, folderId int, orgId int64) (ok bool) {
	stmt := fmt.Sprintf("DELETE FROM %v WHERE folder_id=$1 AND name =$2 AND org_id =$3 ", sketchingBoardImageTable)
	sqlcmd, err := postgresqlDB.Prepare(stmt)
	if err != nil {
		return false
	}
	_, err = sqlcmd.Exec(folderId, name, orgId)
	if err != nil {
		return false
	}
	return true
}

func deleteSketchingBoardFolderByIDPostgresql(id int) (ok bool) {
	stmt := fmt.Sprintf("DELETE FROM %v WHERE id=$1", sketchingBoardFolderTable)
	sqlcmd, err := postgresqlDB.Prepare(stmt)
	if err != nil {
		return false
	}
	_, err = sqlcmd.Exec(id)
	if err != nil {
		return false
	}
	return true
}

func deleteSketchingElementByFolderIdPostgresql(folderId int) (ok bool) {
	stmt := fmt.Sprintf("DELETE FROM %v WHERE folder_id=$1", sketchingElementTable)
	sqlcmd, err := postgresqlDB.Prepare(stmt)
	if err != nil {
		return false
	}
	_, err = sqlcmd.Exec(folderId)
	if err != nil {
		return false
	}
	return true
}

func deleteSketchingBoardImageTableByFolderIdPostgresql(folderId int) (ok bool) {
	stmt := fmt.Sprintf("DELETE FROM %v WHERE folder_id=$1", sketchingBoardImageTable)
	sqlcmd, err := postgresqlDB.Prepare(stmt)
	if err != nil {
		return false
	}
	_, err = sqlcmd.Exec(folderId)
	if err != nil {
		return false
	}
	return true
}

func updateSketchingElementName(oldName string, newName string, dataType string, orgId int64, folderIdReq string) bool {
	stmt := fmt.Sprintf("UPDATE " + sketchingElementTable + " SET title=$1" + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND title = '" + oldName + "'" + " AND data_type = '" + dataType + "'" + " AND folder_id = '" + folderIdReq + "'")
	sqlcmd, err := postgresqlDB.Prepare(stmt)
	if err != nil {
		return false
	}
	_, err = sqlcmd.Exec(newName)
	if err != nil {
		return false
	}
	return true
}

func updateSketchingBoardImageName(oldName string, newName string, dataType string, orgId int64, folderIdReq string) bool {
	stmt := fmt.Sprintf("UPDATE " + sketchingBoardImageTable + " SET name=$1" + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND name = '" + oldName + "'" + " AND data_type = '" + dataType + "'" + " AND folder_id = '" + folderIdReq + "'")
	sqlcmd, err := postgresqlDB.Prepare(stmt)
	if err != nil {
		return false
	}
	_, err = sqlcmd.Exec(newName)
	if err != nil {
		return false
	}
	return true
}

func querySketchingFolderByNameAndOrgIdAndDataType(name string, dataType string, orgId int64) (returnErr bool, id int, version string, createdTime string) {
	stmt := "SELECT id, version,created_time FROM " + sketchingBoardFolderTable + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND name = '" + name + "'" + " AND data_type = '" + dataType + "'"
	rows, err := postgresqlDB.Query(stmt)
	defer rows.Close()
	if checkErr(err) {
		return true, 0, "", ""
	} else {
		for rows.Next() {
			err = rows.Scan(&id, &version, &createdTime)
			if err != nil {
				fmt.Println(err)
				return true, 0, "", ""
			}
		}
		return
	}
}

func querySketchingElement(name string, dataType string, orgId int64, folderIdReq string) (returnErr bool, id int, version string, data string, folderId int, createdTime string, updatedTime string, updatedBy string, createdBy string) {
	stmt := "SELECT id, version,data,folder_id,created_time,updated_time,updated_by,created_by FROM " + sketchingElementTable + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND title = '" + name + "'" + " AND data_type = '" + dataType + "'" + " AND folder_id = '" + folderIdReq + "'"
	rows, err := postgresqlDB.Query(stmt)
	defer rows.Close()
	if checkErr(err) {
		return true, 0, "", "", 0, "", "", "", ""
	} else {
		for rows.Next() {
			err = rows.Scan(&id, &version, &data, &folderId, &createdTime, &updatedTime, &updatedBy, &createdBy)
			if err != nil {
				fmt.Println(err)
				return true, 0, "", "", 0, "", "", "", ""
			}
		}
		return
	}
}

func querySketchingBoardImage(name string, dataType string, orgId int64, folderIdReq string) (returnErr bool, id int, version string, data string, folderId int, updatedTime string) {
	stmt := "SELECT id,version,data,folder_id,updated_time FROM " + sketchingBoardImageTable + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND name = '" + name + "'" + " AND data_type = '" + dataType + "'" + " AND folder_id = '" + folderIdReq + "'"
	rows, err := postgresqlDB.Query(stmt)
	defer rows.Close()
	if checkErr(err) {
		return true, 0, "", "", 0, ""
	} else {
		for rows.Next() {
			err = rows.Scan(&id, &version, &data, &folderId, &updatedTime)
			if err != nil {
				return true, 0, "", "", 0, ""
			}
		}
		return
	}
}

func mergeExplore(folderMap []map[string]interface{}, elementMap []map[string]interface{}, imgMap []map[string]interface{}) map[string]interface{} {
	mergeRes := make(map[string]interface{})
	for i := 0; i < len(folderMap); i++ {
		folderName := fmt.Sprintf("%v", folderMap[i]["name"])
		mergeExploreExe(mergeRes, filterDataTypePath(folderName), elementMap, imgMap, folderMap, folderName)
	}
	return mergeRes
}

func mergeExploreExe(merge map[string]interface{}, folderName string, elementMap []map[string]interface{}, imgMap []map[string]interface{}, folderMap []map[string]interface{}, orgFolderName string) {
	var folderDir = ""
	var isMerge = false
	var sw = false
	var localName = ""
	for i := 0; i < len(folderName); i++ {
		isMerge = true
		if string(folderName[i]) == "/" {
			sw = true
		} else {
			if sw {
				localName = localName + string(folderName[i])
			} else {
				folderDir = folderDir + string(folderName[i])
			}
		}
	}
	if isMerge {
		if merge[folderDir] ==nil{
			item := make(map[string]interface{})
			merge[folderDir] = item		
		} 
		if mergeItem, ok := merge[folderDir].(map[string]interface{}); ok {
			if localName != "" {
				mergeExploreExe(mergeItem, localName, elementMap, imgMap, folderMap, orgFolderName)
			} else {
				setExploreElementMap(mergeItem, orgFolderName, elementMap, folderMap)
				setExploreImgMap(mergeItem, orgFolderName, imgMap, folderMap)
			}
		}
	} else {
		setExploreElementMap(merge, orgFolderName, elementMap, folderMap)
		setExploreImgMap(merge, orgFolderName, imgMap, folderMap)
	}
}

func setExploreElementMap(merge map[string]interface{}, folderName string, elementMap []map[string]interface{}, folderMap []map[string]interface{}) {
	for i := 0; i < len(elementMap); i++ {
		elementFolderId := fmt.Sprintf("%v", elementMap[i]["folderId"])
		elementsFolderName := getFolderNameById(folderMap, elementFolderId)
		if elementsFolderName == folderName {
			mapProt := fmt.Sprintf("%v", elementMap[i]["title"])
			merge[mapProt] = true
		}
	}
}

func setExploreImgMap(merge map[string]interface{}, folderName string, imgMap []map[string]interface{}, folderMap []map[string]interface{}) {
	for i := 0; i < len(imgMap); i++ {
		imgFolderId := fmt.Sprintf("%v", imgMap[i]["folderId"])
		imgFolderName := getFolderNameById(folderMap, imgFolderId)
		if imgFolderName == folderName {
			mapProt := fmt.Sprintf("%v", imgMap[i]["name"])
			merge[mapProt] = true
		}
	}
}

func getFolderNameById(folderMap []map[string]interface{}, id string) string {
	folderName := ""
	for i := 0; i < len(folderMap); i++ {
		folderId := fmt.Sprintf("%v", folderMap[i]["id"])
		if folderId == id {
			folderName = fmt.Sprintf("%v", folderMap[i]["name"])
		}
	}
	return folderName
}

func querySketchingBoardImageByDataTypeAndOrgId(dataType string, orgId int64) (bool, []map[string]interface{}) {
	countRes, count := getSketchingBoardImageCountByDataTypeAndOrgId(dataType, orgId)
	if countRes == false {
		return false, make([]map[string]interface{}, 0)
	}
	result := make([]map[string]interface{}, count)
	countIndex := 0
	stmt := "SELECT id,version,name,data,data_type,folder_id,updated_time FROM " + sketchingBoardImageTable + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND data_type = '" + dataType + "'"
	rows, err := postgresqlDB.Query(stmt)
	defer rows.Close()
	if checkErr(err) {
		return false, make([]map[string]interface{}, 0)
	} else {
		for rows.Next() {
			item := make(map[string]interface{})
			var id, folderId int
			var version, name, data, dataType, updatedTime string
			err = rows.Scan(&id, &version, &name, &data, &dataType, &folderId, &updatedTime)
			if err != nil {
				return false, make([]map[string]interface{}, 0)
			}
			item["id"] = id
			item["version"] = version
			item["name"] = name
			item["data"] = data
			item["dataType"] = dataType
			item["folderId"] = folderId
			item["orgId"] = orgId
			item["updatedTime"] = updatedTime
			result[countIndex] = item
			countIndex++
		}
		return true, result
	}
}

func getSketchingBoardImageCountByDataTypeAndOrgId(dataType string, orgId int64) (bool, int) {
	stmt := "SELECT COUNT(*) as count FROM " + sketchingBoardImageTable + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND data_type = '" + dataType + "'"
	rows, err := postgresqlDB.Query(stmt)
	defer rows.Close()
	var count int
	if checkErr(err) {
		return false, 0
	} else {
		for rows.Next() {
			err := rows.Scan(&count)
			if checkErr(err) {
				return false, 0
			}
		}
	}
	return true, count
}

func querySketchingElementByDataTypeAndOrgId(dataType string, orgId int64) (bool, []map[string]interface{}) {
	countRes, count := getSketchingElementCountByDataTypeAndOrgId(dataType, orgId)
	if countRes == false {
		return false, make([]map[string]interface{}, 0)
	}
	result := make([]map[string]interface{}, count)
	countIndex := 0
	stmt := "SELECT id,version,title,data,data_type,folder_id,created_time,updated_time,updated_by,created_by FROM " + sketchingElementTable + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND data_type = '" + dataType + "'"
	rows, err := postgresqlDB.Query(stmt)
	defer rows.Close()
	if checkErr(err) {
		return false, make([]map[string]interface{}, 0)
	} else {
		for rows.Next() {
			item := make(map[string]interface{})
			var id, folderId int
			var version, title, data, dataType, createdTime, updatedTime, updatedBy, createdBy string
			err = rows.Scan(&id, &version, &title, &data, &dataType, &folderId, &createdTime, &updatedTime, &updatedBy, &createdBy)
			if err != nil {
				return false, make([]map[string]interface{}, 0)
			}
			item["id"] = id
			item["version"] = version
			item["title"] = title
			item["data"] = data
			item["dataType"] = dataType
			item["folderId"] = folderId
			item["orgId"] = orgId
			item["createdTime"] = createdTime
			item["updatedTime"] = updatedTime
			item["updatedBy"] = updatedBy
			item["createdBy"] = createdBy
			result[countIndex] = item
			countIndex++
		}
		return true, result
	}
}

func getSketchingElementCountByDataTypeAndOrgId(dataType string, orgId int64) (bool, int) {
	stmt := "SELECT COUNT(*) as count FROM " + sketchingElementTable + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND data_type = '" + dataType + "'"
	rows, err := postgresqlDB.Query(stmt)
	defer rows.Close()
	var count int
	if checkErr(err) {
		return false, 0
	} else {
		for rows.Next() {
			err := rows.Scan(&count)
			if checkErr(err) {
				return false, 0
			}
		}
	}
	return true, count
}

func querySketchingBoardFolderByDataTypeAndOrgId(dataType string, orgId int64) (bool, []map[string]interface{}) {
	countRes, count := getSketchingBoardFolderCountByDataTypeAndOrgId(dataType, orgId)
	if countRes == false {
		return false, make([]map[string]interface{}, 0)
	}
	result := make([]map[string]interface{}, count)
	countIndex := 0
	stmt := "SELECT id,version,name,created_time FROM " + sketchingBoardFolderTable + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND data_type = '" + dataType + "'"
	rows, err := postgresqlDB.Query(stmt)
	defer rows.Close()
	if checkErr(err) {
		return false, make([]map[string]interface{}, 0)
	} else {
		for rows.Next() {
			item := make(map[string]interface{})
			var id int
			var version, name, createdTime string
			err = rows.Scan(&id, &version, &name, &createdTime)
			if err != nil {
				return false, make([]map[string]interface{}, 0)
			}
			item["id"] = id
			item["version"] = version
			item["name"] = name
			item["orgId"] = orgId
			item["createdTime"] = createdTime
			result[countIndex] = item
			countIndex++
		}
		return true, result
	}
}

func getSketchingBoardFolderCountByDataTypeAndOrgId(dataType string, orgId int64) (bool, int) {
	stmt := "SELECT COUNT(*) as count FROM " + sketchingBoardFolderTable + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND data_type = '" + dataType + "'"
	rows, err := postgresqlDB.Query(stmt)
	defer rows.Close()
	var count int
	if checkErr(err) {
		return false, 0
	} else {
		for rows.Next() {
			err := rows.Scan(&count)
			if checkErr(err) {
				return false, 0
			}
		}
	}
	return true, count
}

func saveSketchingBoardImageExe(name string, data string, dataType string, orgId int64, folderId int) bool {
	_, isExist, isExistId := isExistSketchingBoardImage(name, orgId, folderId)
	if isExist {
		return updateSketchingBoardImage(isExistId, name, data, dataType, orgId, folderId)
	} else {
		return insertSketchingBoardImage(name, data, dataType, orgId, folderId)
	}
}

func updateSketchingBoardImage(id int, name string, data string, dataType string, orgId int64, folderId int) (ok bool) {
	stmt := fmt.Sprintf("UPDATE %v SET version=$1, name=$2, data=$3,data_type=$4,org_id=$5,folder_id=$6 WHERE id=%v", sketchingBoardImageTable, id)
	sqlcmd, err := postgresqlDB.Prepare(stmt)
	if err != nil {
		return false
	}
	_, err = sqlcmd.Exec(version, name, data, dataType, orgId, folderId)
	if err != nil {
		return false
	}
	return true
}

func insertSketchingBoardImage(name string, data string, dataType string, orgId int64, folderId int) (ok bool) {
	stmt := fmt.Sprintf("INSERT INTO %v (version,name, data, data_type, org_id,folder_id) VALUES ($1, $2, $3, $4, $5,$6);", sketchingBoardImageTable)
	sqlcmd, err := postgresqlDB.Prepare(stmt)
	if err != nil {
		return false
	}
	_, err = sqlcmd.Exec(version, name, data, dataType, orgId, folderId)
	if err != nil {
		return false
	}
	return true
}

func isExistSketchingBoardImage(name string, orgId int64, folderId int) (bool, bool, int) {
	rows, err := postgresqlDB.Query("SELECT id FROM " + sketchingBoardImageTable + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND name = '" + name + "'" + " AND folder_id = '" + strconv.Itoa(folderId) + "'")
	defer rows.Close()
	result := false
	var returnId = 0
	if checkErr(err) {
		return true, result, returnId
	} else {
		rowCount := 0
		for rows.Next() {
			var rowId int
			err = rows.Scan(&rowId)
			if err != nil {
			} else {
				returnId = rowId
			}
			rowCount++
		}
		if rowCount > 0 {
			result = true
		}
		return false, result, returnId
	}
}

func saveSketchingElementExe(title string, data string, dataType string, orgId int64, folderId int, user string) (int, bool) {
	_, isExist, isExistId := isExistSketchingElementPostgresql(title, orgId, folderId)
	if isExist {
		return updateSketchingElementPostgresql(isExistId, title, data, dataType, orgId, folderId, user)
	} else {
		return insertSketchingElementPostgresql(title, data, dataType, orgId, folderId, user)
	}
}

func updateSketchingElementPostgresql(id int, title string, data string, dataType string, orgId int64, folderId int, user string) (int, bool) {
	stmt := fmt.Sprintf("UPDATE %v SET title=$1, data=$2, data_type=$3, org_id=$4, folder_id=$5,updated_time=NOW(),updated_by=$6 WHERE id=%v", sketchingElementTable, id)
	sqlcmd, err := postgresqlDB.Prepare(stmt)
	if err != nil {
		return -1, false
	}
	_, err = sqlcmd.Exec(title, data, dataType, orgId, folderId, user)
	if err != nil {
		return -1, false
	}
	return id, true
}

func insertSketchingElementPostgresql(title string, data string, dataType string, orgId int64, folderId int, user string) (int, bool) {
	var id = -1
	stmt := fmt.Sprintf("INSERT INTO %v (version,title,data,data_type,org_id,folder_id,updated_by,created_by) VALUES ($1, $2,$3,$4,$5,$6,$7,$8) RETURNING id;", sketchingElementTable)
	err := postgresqlDB.QueryRow(stmt, version, title, data, dataType, orgId, folderId, user, user).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return id, false
	}
	return id, true
}

func isExistSketchingElementPostgresql(title string, orgId int64, folderId int) (bool, bool, int) {
	rows, err := postgresqlDB.Query("SELECT id FROM " + sketchingElementTable + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND title = '" + title + "'" + " AND folder_id = '" + strconv.Itoa(folderId) + "'")
	defer rows.Close()
	result := false
	var returnId = 0
	if checkErr(err) {
		return true, result, returnId
	} else {
		rowCount := 0
		for rows.Next() {
			var rowId int
			err = rows.Scan(&rowId)
			if err != nil {
			} else {
				returnId = rowId
			}
			rowCount++
		}
		if rowCount > 0 {
			result = true
		}
		return false, result, returnId
	}
}

func queryFolderByNameAndOrgId(name string, orgId int64) (returnErr bool, id int, version string, createdTime string) {
	stmt := "SELECT id,version,created_time FROM " + sketchingBoardFolderTable + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND name = '" + name + "'"
	rows, err := postgresqlDB.Query(stmt)
	defer rows.Close()
	if checkErr(err) {
		return true, 0, "", ""
	} else {
		for rows.Next() {
			err = rows.Scan(&id, &version, &createdTime)
			if err != nil {
				return true, 0, "", ""
			}
		}
		return
	}
}

func isExistHTFolder(name string, orgId int64) (bool, bool, int) {
	// rows, err := postgresqlDB.Query("SELECT id FROM " + sketchingBoardFolderTable + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND name = '" + name + "'")
	rows, err := postgresqlDB.Query("SELECT id FROM " + sketchingBoardFolderTable + " WHERE org_id = " + strconv.FormatInt(orgId, 10) + " AND (name = '" + strings.Replace(name, "\\", "/", -1) + "' OR name = '" + strings.Replace(name, "/", "\\", -1) + "')")
	defer rows.Close()
	result := false
	var returnId = 0
	if checkErr(err) {
		return true, result, returnId
	} else {
		rowCount := 0
		for rows.Next() {
			var rowId int
			err = rows.Scan(&rowId)
			if err != nil {
			} else {
				returnId = rowId
			}
			rowCount++
		}
		if rowCount > 0 {
			result = true
		}
		return false, result, returnId
	}
}

func insertHTFolderPostgresql(name string, orgId int64) (int, bool) {
	var id = -1
	if name == "" || name == "." {
		return 0, true
	}
	dataType := getDataType(name)
	stmt := fmt.Sprintf("INSERT INTO %v (name,version,org_id,data_type) VALUES ($1, $2,$3,$4) RETURNING id;", sketchingBoardFolderTable)
	err := postgresqlDB.QueryRow(stmt, name, version, orgId, dataType).Scan(&id)
	if err != nil {
		return id, false
	}
	return id, true
}
