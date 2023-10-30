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
	"time"
)

/* mogoDB schema init rule */

var default_admin_schema = []any{
	RuleInfo{
		Name:      "admin",
		Authority: "/department/querys",
		// Type:      "管理",
		// Operation: "Query Scan",
	},
	RuleInfo{
		Name:      "admin",
		Authority: "/log/remove",
		// Type:      "管理",
		// Operation: "Create Delete Scan Query Update",
	},
	RuleInfo{
		Name:      "admin",
		Authority: "/log/query",
		// Type:      "系统",
		// Operation: "Query",
	},
	RuleInfo{
		Name:      "admin",
		Authority: "/user/filter",
		// Type:      "管理",
		// Operation: "Query Delete Scan Update",
	}}

//TODO write the new python resource to the program
// Rule{
// 	Name:   "admin",
// 	Authority: "",
// }

var default_member_schma = []any{
	// RuleInfo{
	// 	Name:      "member",
	// 	Authority: "",
	// 	// Type:      "功能",
	// 	// Operation: "Query Scan",
	// },
	// RuleInfo{
	// 	Name:      "member",
	// 	Authority: "层级图",
	// 	// Type:      "功能",
	// 	// Operation: "Query Scan",
	// },
	// RuleInfo{
	// 	Name:      "member",
	// 	Authority: "文件切分",
	// 	// Type:      "功能",
	// 	// Operation: "Query Scan",
	// },
}

var default_departmentManage_schema = []any{
	RuleInfo{
		Name:      "departmentManage",
		Authority: "/user/filter",
	},
	RuleInfo{
		Name:      "departmentManage",
		Authority: "/department/querys",
	},
}
var default_audit_schema = []any{
	RuleInfo{
		Name:      "audit",
		Authority: "/log/query",
		// Type:      "系统",
		// Operation: "Query Scan",
	},
	RuleInfo{
		Name:      "audit",
		Authority: "/log/remove",
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

var defualt_department_schema = []any{DepartmentInfo{
	Name:        "案件服务部",
	Description: "案件服务部",
	CreateTime:  time.Now().Unix(),
	Leaders:     []string{"xiyang"},
}, DepartmentInfo{
	Name:        "财务部",
	Description: "财务部",
	CreateTime:  time.Now().Unix(),
	Leaders:     []string{"xiaoli"},
}, DepartmentInfo{
	Name:        "售前部门",
	Description: "售前部门",
	CreateTime:  time.Now().Unix(),
	Leaders:     []string{"xiaozhang"},
}}

func initDepartmentSchema() {
	getDepartmentCollection().Drop(context.Background())
	_, err := getDepartmentCollection().InsertMany(context.TODO(), defualt_department_schema)
	if err != nil {
		panic(fmt.Sprintf("初始化department schema出错:%s", err.Error()))
	}
	log.Println("初始化department Collections成功")

}

// *********************************************** init the database to use **************************************************************//
func InitMogoSchema() {
	// initDB()
	initRuleSchema()
	// 为了防止在分布式条件下，更新后的其他服务器配置文件未能及时更新从而导致的数据错误,只在user collections 中追加user数据
	initUserSchema()
	// initLogSchema()
	initFuncMapSchema()
	initDepartmentSchema()
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
	// ruleCollections = append(ruleCollections, default_member_schma...)
	ruleCollections = append(ruleCollections, default_audit_schema...)
	ruleCollections = append(ruleCollections, default_departmentManage_schema...)
	//in order to create the database
	// _, err := getRuleCollection()(context.Background(), map[string]struct{}{})
	// if err != nil {
	// 	panic(err.Error())
	// }
	err := getRuleCollection().Drop(context.Background())
	if err != nil {
		panic(fmt.Sprintf("drop the rule schema collection error: %v", err))
	}
	_, err = getRuleCollection().InsertMany(context.Background(), ruleCollections)
	if err != nil {
		panic(fmt.Sprintf("insert the rule schema collection error: %v", err))
	}
	log.Printf("成功初始化Rule collections:%v", ruleCollections)
}

var default_users_schema = []any{
	UserInfo{
		Username:    "admin",
		Authority:   "admin",
		RealName:    "管理员",
		Password:    utils.Encrypt("admin"),
		PhoneNumber: "18080705675",
		Department:  "无",
		CreateTime:  time.Now().Unix(),
	},
	UserInfo{
		Username:    "audit",
		Authority:   "audit",
		Password:    utils.Encrypt("audit"),
		PhoneNumber: "18080705675",
		RealName:    "审计员",
		Department:  "无",
		CreateTime:  time.Now().Unix(),
	},
	UserInfo{
		Username:    "member",
		Authority:   "member",
		Password:    utils.Encrypt("member"),
		PhoneNumber: "18080705675",
		RealName:    "会员",
		Department:  "无",
		CreateTime:  time.Now().Unix(),
	},
	UserInfo{
		Username:    "案件服务",
		Authority:   "member",
		Password:    utils.Encrypt("member"),
		PhoneNumber: "18080705675",
		RealName:    "会员",
		Department:  "案件服务部",
		CreateTime:  time.Now().Unix(),
	},
	UserInfo{
		Username:    "财务部",
		Authority:   "member",
		Password:    utils.Encrypt("member"),
		PhoneNumber: "18080705675",
		RealName:    "会员",
		Department:  "财务部",
		CreateTime:  time.Now().Unix(),
	},
	UserInfo{
		Username:    "售前部门",
		Authority:   "member",
		Password:    utils.Encrypt("member"),
		PhoneNumber: "18080705675",
		RealName:    "会员",
		Department:  "售前部门",
		CreateTime:  time.Now().Unix(),
	},
}

var default_funcmaps_schema = []any{
	FuncMap{
		Function:   "null",
		Command:    "null",
		OSType:     Undefiend,
		Type:       None,
		Additional: "null",
	},
}

func initUserSchema() {
	// var users = []any{}
	getUserCollection().Drop(context.Background())
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
func initFuncMapSchema() {
	getFuncMapCollection().Drop(context.Background())
	_, err := getFuncMapCollection().InsertMany(context.Background(), default_funcmaps_schema)
	if err != nil {
		panic(fmt.Sprintf("初始化funcmap document失败:%s", err.Error()))
	}
	log.Printf("成功初始化funcmap collections:%v", default_funcmaps_schema...)
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
