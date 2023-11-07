package config

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

//this meaning the server config

type ServerConfig struct {
	XMLName     xml.Name     `xml:"server"`
	Secret      string       `xml:"secret"`
	Port        string       `xml:"port"`
	SMSConfig   *SMSConfig   `xml:"sms"`
	Kafka       *KafkaConfig `xml:"kafka"`
	MeiliSearch *MeiliSearch `xml:"meiliSearch"`
}

type KafkaConfig struct {
	Broker []string `xml:"broker"`
	Topic  []string `xml:"topic"`
}

type SMSConfig struct {
	// XMLName xml.Name `xml:"sms"`
	Key          string `xml:"key"`
	Secret       string `xml:"secret"`
	TemplateCode string `xml:"templateCode"`
}

type MeiliSearch struct {
	Key  string `xml:"key"`
	Host string `xml:"host"`
}

type PrometheusConfig struct {
	Server string `xml:"server"`
}

var serverConf *ServerConfig

var DEFUALT_SERVER_CONFIG_FILE = map[string]string{"develop": fmt.Sprint(".", string(filepath.Separator), "configs", string(filepath.Separator), "server.xml"), "test": fmt.Sprint("..", string(filepath.Separator), "configs", string(filepath.Separator), "server.xml")}

//  init the envrionment to make configuration

func init() {
	file, err := os.Open(DEFUALT_SERVER_CONFIG_FILE[string(CurEnviroment)])
	if err != nil {
		panic(fmt.Sprintf("初始化server配置文件出错:%s", err.Error()))
	}
	serverConf = new(ServerConfig)
	err = xml.NewDecoder(bufio.NewReader(file)).Decode(serverConf)
	if err != nil {
		panic(fmt.Sprintf("xml解析server配置文件出错:%s", err.Error()))
	}
	log.Printf("server的配置如下:%v", serverConf)
}

func GetServerConfig() *ServerConfig {
	return serverConf
}
