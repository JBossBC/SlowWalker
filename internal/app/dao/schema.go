package dao

import (
	"bufio"
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"replite_web/internal/app/config"
	"replite_web/internal/app/utils"
)

/* mogoDB schema init rule */

var default_admin_schema = []any{
	Rule{
		Name:      "admin",
		Authority: "fileSystem",
		Type:      "system",
	},
	Rule{
		Name:      "admin",
		Authority: "log",
		Type:      "system",
	},
	//TODO write the new python resource to the program
	// Rule{
	// 	Name:   "admin",
	// 	Authority: "",
	// }
}

var default_member_schma = []any{
	Rule{
		Name:      "member",
		Authority: "ipQuery",
		Type:      "function",
	},
	Rule{
		Name:      "member",
		Authority: "levelGraph",
		Type:      "function",
	},
	Rule{
		Name:      "member",
		Authority: "fileCut",
		Type:      "function",
	},
}
var default_audit_schema = []any{
	Rule{
		Name:      "audit",
		Authority: "log",
		Type:      "system",
	},
}

// TODO test to init
var default_platform_schema = []PlatForm{
	&RemotePlatForm{
		BasePlatForm{
			CoreType:    []Core{None},
			MechineType: Undefiend,
		},
	},
	&RemotePlatForm{
		BasePlatForm{
			CoreType:    []Core{CPU},
			MechineType: Linux,
		},
	},
	// &RemotePlatForm{
	// 	BasePlatForm{

	// 	}
	// }

	// &LocalPlatForm{
	// 	CoreType:    None,
	// 	MechineType: Undefiend,
	// 	Command:     "",
	// },
}

// *********************************************** init the database to use **************************************************************//
func InitMogoSchema() {
	// initDB()
	initRuleSchema()
	// 为了防止在分布式条件下，更新后的其他服务器配置文件未能及时更新从而导致的数据错误,只在user collections 中追加user数据
	initUserSchema()
	// initLogSchema()
	//renew the db.xml the init state
	go func() {
		config.DBConfig.MongoConfig.Init = "true"
		bytes, err := xml.MarshalIndent(config.DBConfig, "", "  ")
		if err != nil {
			log.Printf("序列化dbconfig出错%s", err.Error())
			return
		}
		// file.Write(bytes)
		// doc := etree.NewDocument()
		// if err := doc.ReadFromString(string(bytes)); err != nil {
		// 	log.Printf("格式化db配置文件时出错%s", err.Error())
		// 	return
		// }
		// prettyData, err := doc.WriteToString()
		// if err != nil {
		// 	log.Printf("格式化db配置文件时出错%s", err.Error())
		// 	return
		// }
		file, _ := os.OpenFile(config.DEFAULT_DB_CONFIG, os.O_TRUNC|os.O_WRONLY, 0755)
		writer := bufio.NewWriter(file)
		_, err = writer.Write(bytes)
		if err != nil {
			log.Printf("修改db配置文件出错%s", err.Error())
			return
		}
		err = writer.Flush()
		if err != nil {
			log.Printf("修改db配置文件出错%s", err.Error())
			return
		}
		log.Printf("修改db配置文件成功")
	}()
}

//	func initDB() {
//		getMongoClient().Database(dbConfig.MongoConfig.Database)
//	}
func initRuleSchema() {
	var ruleCollections = []any{}
	ruleCollections = append(ruleCollections, default_admin_schema...)
	ruleCollections = append(ruleCollections, default_member_schma...)
	ruleCollections = append(ruleCollections, default_audit_schema...)
	//in order to create the database
	_, err := getRuleCollection().InsertOne(context.Background(), map[string]struct{}{})
	if err != nil {
		panic(err.Error())
	}
	err = getRuleCollection().Drop(context.Background())
	if err != nil {
		panic(fmt.Sprintf("drop the rule schema collection error: %v", err))
	}
	_, err = getRuleCollection().InsertMany(context.Background(), ruleCollections)
	log.Printf("成功初始化Rule collections:%v", ruleCollections)
	if err != nil {
		panic(fmt.Sprintf("insert the rule schema collection error: %v", err))
	}
}

var default_users_schema = []any{
	User{
		Username:    "admin",
		Authority:   "admin",
		Password:    utils.Encrypt("admin"),
		PhoneNumber: "18080705675",
	},
	User{
		Username:    "audit",
		Authority:   "audit",
		Password:    utils.Encrypt("audit"),
		PhoneNumber: "18080705675",
	},
	User{
		Username:    "member",
		Authority:   "member",
		Password:    utils.Encrypt("member"),
		PhoneNumber: "18080705675",
	},
}

func initUserSchema() {
	// var users = []any{}
	_, err := getUserCollection().InsertMany(context.Background(), default_users_schema)
	if err != nil {
		panic(fmt.Sprintf("初始化user document失败:%s", err.Error()))
	}
	log.Printf("成功初始化user collections:%v", default_users_schema)
	// err = getUserCollection().Drop(context.Background())
	// if err != nil {
	// 	panic(fmt.Sprintf("drop the user schema collection error: %v", err))
	// }
	// _, err = getUserCollection().InsertMany(context.Background(), users)
	// if err != nil {
	// 	panic(fmt.Sprintf("insert the user schema collection error: %v", err))
	// }
}

// func initLogSchema() {
// 	var logs = []any{}
// 	_, err := getLogCollection().InsertOne(context.Background(), map[string]struct{}{})
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	err = getLogCollection().Drop(context.Background())
// 	if err != nil {
// 		panic(fmt.Sprintf("drop the log schema collection error: %v", err))
// 	}
// 	_, err = getLogCollection().InsertMany(context.Background(), logs)
// 	if err != nil {
// 		panic(fmt.Sprintf("insert the log schema collection error: %v", err))
// 	}
// }
