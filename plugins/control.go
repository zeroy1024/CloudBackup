package plugins

import (
	"encoding/json"
	"io/ioutil"
	"log"
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


