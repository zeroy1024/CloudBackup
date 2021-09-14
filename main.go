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

	//rootId := "610cc439004f2feb6e154465b6b4feec399f36d3"
	//sha1Sum := sha1sum("./Downloads/rufus-3.14.exe")

	/*userToken := plugins.GetUserToken()

	cloudBackupFolder := GetCloudBackupFolder(userToken)

	fmt.Println(cloudBackupFolder.FileId)*/

	/*fileList := plugins.GetFileList(userToken.AccessToken, rootId, userToken.DriveId)

	for i := range fileList.Items {
		if fileList.Items[i].Name == "rufus-3.14.exe" {
			fmt.Println(strings.ToLower(fileList.Items[i].ContentHash) == strings.ToLower(sha1Sum))
			fmt.Println(fileList.Items[i].ContentHashName)
			fmt.Println(fileList.Items[i].Crc64Hash)
		}
	}*/

	/*pwd, _ := os.Getwd()
	fileList, _ := ioutil.ReadDir(pwd)

	for i := range fileList {
		fmt.Println(pwd + fileList[i].Name())
		fmt.Println(fileList[i].IsDir())
	}*/

	/*watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./plugins/AliDrive.go")
	if err != nil {
		log.Fatal(err)
	}
	<-done*/
	/*rootId := "613a050373fd0da4e101492eae9ae0f6d3c33d3c"
	userToken := plugins.GetUserToken()
	fileInfo := plugins.UploadFile(userToken.AccessToken, rootId, userToken.DriveId, "./Downloads/rufus-3.14.exe")
	fmt.Println(fileInfo)*/
}
