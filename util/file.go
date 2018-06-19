package util

import (
	"path/filepath"
	"strings"
)

func IsFile(name string) bool {
	ext := strings.TrimLeft(filepath.Ext(name), ".")
	if ext != "" {
		return true
	}
	return false
}

func IsImg(name string) bool {
	ext := strings.TrimLeft(filepath.Ext(name), ".")
	if ext == "png" || ext == "gif" || ext == "jpg" || ext == "jpeg" || ext == "svg" {
		return true
	}
	return false
}

func GetDirAndFilename(str string) string {
	dropExt := filepath.Ext(str)
	result := string(str[0 : len(str)-len(dropExt)])
	return result
}

func GetImgNameAndImgDataByImageDataB64(fileName string, imageDataB64 string) (string, string) {
	const imgPrefix = "data:image/"
	const imgSurffix = ";"
	const imgSeparator = ","
	var imgName string
	var imgData string
	start := strings.Index(imageDataB64, imgPrefix)
	end := strings.Index(imageDataB64, imgSurffix)
	separator := strings.Index(imageDataB64, imgSeparator)
	if start == -1 || end == -1 || separator == -1 {
		return "", ""
	} else {
		imgName = imageDataB64[start+len(imgPrefix) : end]
		imgData = imageDataB64[separator+len(imgSeparator):]

		imgName = fileName + "." + imgName
	}
	return imgName, imgData
}
