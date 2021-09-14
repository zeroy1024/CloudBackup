package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func RefreshAccessToken(refreshToken string) UserToken {
	url := "https://api.aliyundrive.com/token/refresh"
	data := "{\"refresh_token\": \"" + refreshToken + "\"}"
	client := http.Client{}
	request, _ := http.NewRequest("POST", url, strings.NewReader(data))
	request.Header.Add("User-Agent", "None")
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Printf("RefreshToken 请求异常")
		panic("RefreshToken 请求异常")
	}

	body, _ := ioutil.ReadAll(response.Body)

	var token UserToken
	_ = json.Unmarshal(body, &token)

	return token
}

func GetUserToken() UserToken {
	config := LoadConfig()
	userToken := config.UserToken

	parse, _ := time.Parse("2006-01-02T15:04:05Z", userToken.ExpiresTime)
	expireTimestamp := parse.Unix()
	nowTimestamp := time.Now().Unix()

	if userToken.AccessToken != "" && (expireTimestamp-nowTimestamp) <= int64(30*60) {
		userToken = RefreshAccessToken(userToken.RefreshToken)
		UpdateUserToken(userToken)
		fmt.Println("刷新token")
	}

	return userToken
}

func CreateFolder(accessToken string, fileId string, driveId string, folderName string) File {
	url := "https://api.aliyundrive.com/adrive/v2/file/createWithFolders"
	data := map[string]interface{}{
		"drive_id":        driveId,
		"parent_file_id":  fileId,
		"name":            folderName,
		"type":            "folder",
		"check_name_mode": "refuse",
	}
	dataBytes, _ := json.Marshal(data)

	client := http.Client{}
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(dataBytes)))
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Content-Type", "application/json")

	response, _ := client.Do(request)
	body, _ := ioutil.ReadAll(response.Body)

	var folderInfo File
	_ = json.Unmarshal(body, &folderInfo)

	return folderInfo
}

func GetFileList(accessToken string, fileId string, driveId string) FileList {
	url := "https://api.aliyundrive.com/adrive/v3/file/list"

	data := map[string]interface{}{
		"drive_id":       driveId,
		"parent_file_id": fileId,
		"limit":          200,
		"fields":         "*",
		"url_expire_sec": 1600,
	}
	dataBytes, _ := json.Marshal(data)

	client := http.Client{}
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(dataBytes)))
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	request.Header.Add("Authorization", accessToken)
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Println("文件列表请求异常: " + err.Error())
	}

	body, _ := ioutil.ReadAll(response.Body)
	var fileList FileList
	_ = json.Unmarshal(body, &fileList)

	return fileList
}

func GetFile(accessToken string, fileId string, driveId string) File {
	url := "https://api.aliyundrive.com/v2/file/get"
	data := map[string]interface{}{
		"drive_id": driveId,
		"file_id":  fileId,
	}
	dataBytes, _ := json.Marshal(data)

	client := http.Client{}
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(dataBytes)))
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Println("文件列表请求异常: " + err.Error())
	}

	body, _ := ioutil.ReadAll(response.Body)

	var file File
	_ = json.Unmarshal(body, &file)

	return file
}

func DownloadFile(fileName string, url string) {
	client := http.Client{}

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Referer", "https://www.aliyundrive.com/")
	response, _ := client.Do(request)

	body, _ := ioutil.ReadAll(response.Body)

	_ = ioutil.WriteFile("./Downloads/"+fileName, body, 0755)
}

func UploadFile(accessToken string, fileId string, driveId string, filePath string) File {
	var filePathSplit []string
	if strings.Index(filePath, "/") != -1 {
		filePathSplit = strings.Split(filePath, "/")
	} else {
		filePathSplit = strings.Split(filePath, "\\")
	}

	fileName := filePathSplit[len(filePathSplit)-1]
	file, _ := ioutil.ReadFile(filePath)
	fileSize := len(file)

	url := "https://api.aliyundrive.com/adrive/v2/file/createWithFolders"
	data := map[string]interface{}{
		"drive_id":        driveId,
		"parent_file_id":  fileId,
		"name":            fileName,
		"type":            "file",
		"check_name_mode": "auto_rename",
		"part_info_list": []map[string]int{
			{
				"part_number": 1,
			},
		},
		"size": fileSize,
	}
	dataBytes, _ := json.Marshal(data)

	client := http.Client{}
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(dataBytes)))
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Content-Type", "application/json")

	response, _ := client.Do(request)
	body, _ := ioutil.ReadAll(response.Body)

	var parentDetail ParentDetail
	_ = json.Unmarshal(body, &parentDetail)

	f, _ := os.Open(filePath)
	Upload(f, parentDetail.PartInfoList[0].UploadUrl)
	fileInfo := Complete(accessToken, parentDetail.DriveId, parentDetail.FileId, parentDetail.UploadId)

	return fileInfo
}

func Upload(file *os.File, url string) {

	client := http.Client{}
	request, _ := http.NewRequest("PUT", url, file)
	request.Header.Add("Referer", "https://www.aliyundrive.com/")

	_, _ = client.Do(request)
}

func Complete(accessToken string, driveId string, fileId string, uploadId string) File {
	url := "https://api.aliyundrive.com/v2/file/complete"
	data := map[string]interface{}{
		"drive_id":  driveId,
		"file_id":   fileId,
		"upload_id": uploadId,
	}
	dataBytes, _ := json.Marshal(data)

	client := http.Client{}
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(dataBytes)))
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Content-Type", "application/json")
	response, _ := client.Do(request)
	body, _ := ioutil.ReadAll(response.Body)

	var file File
	_ = json.Unmarshal(body, &file)

	return file
}
