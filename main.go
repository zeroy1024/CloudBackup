package main

import (
	"CloudBackup/plugins"
	"fmt"
	"io/ioutil"
)

func main() {
	config := plugins.LoadConfig()
	userToken := plugins.GetUserToken()
	cloudBackupFolder := plugins.CreateFolder(userToken.AccessToken, "root", userToken.DriveId, config.CloudBackupFolderName)
	for i := range config.BackupPath {
		fileList, _ := ioutil.ReadDir(config.BackupPath[i])
		splitPath := plugins.SplitPath(config.BackupPath[i])
		folderName := splitPath[len(splitPath)-1]
		subFolderInfo := plugins.CreateFolder(userToken.AccessToken, cloudBackupFolder.FileId, cloudBackupFolder.DriveId, folderName)

		fmt.Println(subFolderInfo)
		for j := range fileList {
			filePath := fmt.Sprintf("%s\\%s", config.BackupPath[i], fileList[j].Name())
			fileInfo := plugins.UploadFile(userToken.AccessToken, subFolderInfo.FileId, subFolderInfo.DriveId, filePath)
			fmt.Printf("文件 %s 上传完毕, 大小 %d.\n", fileInfo.Name, fileInfo.Size)
		}
	}
}
