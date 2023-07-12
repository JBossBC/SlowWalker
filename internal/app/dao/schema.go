package dao

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"replite_web/internal/app/config"
)

/* mogoDB schema init rule */

var default_admin_schema = []any{
	Rule{
		Name:      "admin",
		Authority: "fileSystem",
	},
	Rule{
		Name:      "admin",
		Authority: "log",
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
	},
	Rule{
		Name:      "member",
		Authority: "levelGraph",
	},
	Rule{
		Name:      "member",
		Authority: "fileCut",
	},
}
var default_audit_schema = []any{
	Rule{
		Name:      "audit",
		Authority: "log",
	},
}

// *********************************************** init the database to use **************************************************************//
func InitMogoSchema() {
	// initDB()
	initRuleSchema()
	// initUserSchema()
	// initLogSchema()
	//renew the db.xml the init state
	go func() {
		config.DBConfig.MongoConfig.Init = "true"
		bytes, _ := xml.Marshal(config.DBConfig)
		file, _ := os.OpenFile(config.DEFAULT_DB_CONFIG, os.O_TRUNC|os.O_WRONLY, 0755)
		file.Write(bytes)
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
	if err != nil {
		panic(fmt.Sprintf("insert the rule schema collection error: %v", err))
	}
}

// func initUserSchema() {
// 	var users = []any{}
// 	_, err := getUserCollection().InsertOne(context.Background(), map[string]struct{}{})
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	err = getUserCollection().Drop(context.Background())
// 	if err != nil {
// 		panic(fmt.Sprintf("drop the user schema collection error: %v", err))
// 	}
// 	_, err = getUserCollection().InsertMany(context.Background(), users)
// 	if err != nil {
// 		panic(fmt.Sprintf("insert the user schema collection error: %v", err))
// 	}
// }

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
