package plugins

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func SaveConfig(config Config) {
	configByte, err := json.Marshal(config)
	if err != nil {
		log.Println("配置文件解析失败: " + err.Error())
		return
	}

	err = ioutil.WriteFile("./config.json", configByte, 0755)
	if err != nil {
		log.Println("写入配置文件出错: " + err.Error())
		return
	}
}

func LoadConfig() Config {
	configByte, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Println("读取配置文件出错")
	}

	var config Config
	_ = json.Unmarshal(configByte, &config)

	return config
}

func UpdateUserToken(userToken UserToken) {
	config := LoadConfig()
	config.UserToken = userToken
	SaveConfig(config)
}

func Sha1Sum(filePath string) string {
	file, _ := ioutil.ReadFile(filePath)
	return fmt.Sprintf("%x", sha1.Sum(file))
}

func SplitPath(path string) []string {
	var pathSplit []string
	if strings.Index(path, "/") != -1 {
		pathSplit = strings.Split(path, "/")
	} else {
		pathSplit = strings.Split(path, "\\")
	}

	return pathSplit
}
