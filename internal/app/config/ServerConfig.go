package config

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

//this meaning the server config

type ServerConfig struct {
	XMLName     xml.Name   `xml:"server"`
	Secret      string     `xml:"secret"`
	Port        string     `xml:"port"`
	SMSConfig   *SMSConfig `xml:"sms"`
	Environment string     `xml:"environment"`
}

type SMSConfig struct {
	// XMLName xml.Name `xml:"sms"`
	Key          string `xml:"key"`
	Secret       string `xml:"secret"`
	TemplateCode string `xml:"templateCode"`
}

var ServerConf *ServerConfig

var DEFUALT_SERVER_CONFIG_FILE = fmt.Sprint(".", string(filepath.Separator), "configs", string(filepath.Separator), "server.xml")

//  init the envrionment to make configuration

func init() {
	file, err := os.Open(DEFUALT_SERVER_CONFIG_FILE)
	if err != nil {
		panic(fmt.Sprintf("初始化server配置文件出错:%s", err.Error()))
	}

	ServerConf = new(ServerConfig)
	err = xml.NewDecoder(bufio.NewReader(file)).Decode(ServerConf)
	if err != nil {
		panic(fmt.Sprintf("xml解析server配置文件出错:%s", err.Error()))
	}
}
