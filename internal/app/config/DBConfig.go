package config

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

type DataBase struct {
	XMLName     xml.Name    `xml:"dbConfig"`
	MongoConfig MongoConfig `xml:"mongoConfig"`
	RedisConfig RedisConfig `xml:"redisConfig"`
}
type MongoConfig struct {
	URL      string `xml:"url"`
	Database string `xml:"database"`
	Init     string `xml:"init"`
}

type RedisConfig struct {
	Address  string `xml:"address"`
	Username string `xml:"username"`
	Passowrd string `xml:"password"`
	Database string `xml:"database"`
	Init     string `xml:"init"`
}

var DEFAULT_DB_CONFIG = fmt.Sprint(".", string(filepath.Separator), "configs", string(filepath.Separator), "db.xml")

var DBConfig *DataBase

func init() {
	config, err := os.Open(DEFAULT_DB_CONFIG)
	if err != nil {
		panic(fmt.Sprintf("database config file %s is error: %s", DEFAULT_DB_CONFIG, err.Error()))
	}
	DBConfig = new(DataBase)
	err = xml.NewDecoder(bufio.NewReader(config)).Decode(DBConfig)
	if err != nil {
		panic(fmt.Sprintf("analysis the xml format error: %s", err))
	}
	//TODO if the init is true,config will renew the database ,and update the init value to keep the config file
}
