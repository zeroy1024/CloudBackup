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

	UpdateUserToken(token)
	fmt.Println("UserToken 已更新!")

	return token
}

func GetUserToken() UserToken {
	config := LoadConfig()
	userToken := config.UserToken
	if userToken.RefreshToken == "" {
		var refreshToken string
		fmt.Print("请输入RefreshToken: ")
		_, _ = fmt.Scanln(&refreshToken)
		userToken = RefreshAccessToken(refreshToken)
	} else {
		parse, _ := time.Parse("2006-01-02T15:04:05Z", userToken.ExpiresTime)
		expireTimestamp := parse.Unix()
		nowTimestamp := time.Now().Unix()

		if (expireTimestamp - nowTimestamp) <= int64(30*60) {
			userToken = RefreshAccessToken(userToken.RefreshToken)
		}
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
	filePathSplit := SplitPath(filePath)

	fileName := filePathSplit[len(filePathSplit)-1]
	file, _ := ioutil.ReadFile(filePath)
	fileSize := len(file)

	url := "https://api.aliyundrive.com/adrive/v2/file/createWithFolders"
	data := map[string]interface{}{
		"drive_id":        driveId,
		"parent_file_id":  fileId,
		"name":            fileName,
		"type":            "file",
		"check_name_mode": "refuse",
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

	var fileInfo File
	var uploadResult File
	_ = json.Unmarshal(body, &uploadResult)
	if uploadResult.Exist == true {
		fileInfo = GetFile(accessToken, uploadResult.FileId, uploadResult.DriveId)
		if strings.ToUpper(fileInfo.ContentHash) != strings.ToUpper(Sha1Sum(filePath)) {
			DeleteFile(accessToken, fileInfo.FileId, fileInfo.DriveId)
			fileInfo = UploadFile(accessToken, fileId, driveId, filePath)
		}
		return fileInfo
	}

	f, _ := os.Open(filePath)
	Upload(f, uploadResult.PartInfoList[0].UploadUrl)
	fileInfo = Complete(accessToken, uploadResult.DriveId, uploadResult.FileId, uploadResult.UploadId)

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

func DeleteFile(accessToken string, fileId string, driveId string) {
	url := "https://api.aliyundrive.com/v2/recyclebin/trash"
	data := map[string]interface{}{
		"drive_id": driveId,
		"file_id":  fileId,
	}
	dataBytes, _ := json.Marshal(data)

	client := http.Client{}
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(dataBytes)))
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Content-Type", "application/json")
	_, _ = client.Do(request)
}

func Batch(accessToken string, batchRequestList []BatchRequest) {
	url := "https://api.aliyundrive.com/v2/batch"

	marshal, _ := json.Marshal(batchRequestList)

	data := map[string]interface{}{
		"requests": marshal,
		"resource": "file",
	}
	dataBytes, _ := json.Marshal(data)

	client := http.Client{}
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(dataBytes)))
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Content-Type", "application/json")

	response, _ := client.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(body)
}
