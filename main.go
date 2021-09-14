package main

import (
	"CloudBackup/plugins"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
)

func sha1sum(filePath string) string {
	file, _ := ioutil.ReadFile(filePath)
	return fmt.Sprintf("%x", sha1.Sum(file))
}

func GetCloudBackupFolder(userToken plugins.UserToken) plugins.File {
	rootList := plugins.GetFileList(userToken.AccessToken, "root", userToken.DriveId)
	var cloudBackupFolder plugins.File
	for i := range rootList.Items {
		if rootList.Items[i].Name == "CloudBackup" {
			cloudBackupFolder = rootList.Items[i]
		}
	}

	if cloudBackupFolder.FileId == "" {
		cloudBackupFolder = plugins.CreateFolder(userToken.AccessToken, "root", userToken.DriveId, "CloudBackup")
	}

	return cloudBackupFolder
}

func GetBackupPathList() []string {
	config := plugins.LoadConfig()
	return config.BackupPath
}

func main() {
	userToken := plugins.GetUserToken()
	cloudBackupFolder := GetCloudBackupFolder(userToken)
	backupPathList := GetBackupPathList()

	for i := range backupPathList {
		fileList, _ := ioutil.ReadDir(backupPathList[i])
		for j := range fileList {
			filePath := fmt.Sprintf("%s\\%s", backupPathList[i], fileList[j].Name())
			fileInfo := plugins.UploadFile(userToken.AccessToken, cloudBackupFolder.FileId, cloudBackupFolder.DriveId, filePath)
			fmt.Printf("文件 %s 上传完毕, 大小 %d.\n", fileInfo.Name, fileInfo.Size)
		}
	}
}
