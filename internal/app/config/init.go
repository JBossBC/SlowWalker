package config

// var DEFAULT_COLLECTIONS_CONFIG string
// var DEFUALT_SERVER_CONFIG_FILE string
// var CollectionConfig *staticMap

// var environmentMap = map[string][]string{
// 	"develop": []string{
// 		fmt.Sprint(".", string(filepath.Separator), "configs", string(filepath.Separator), "collections.xml"),
// 		fmt.Sprint(".", string(filepath.Separator), "configs", string(filepath.Separator), "db.xml"),
// 		fmt.Sprint(".", string(filepath.Separator), "configs", string(filepath.Separator), "server.xml"),
// 	},
// 	"test": []string{
// 		fmt.Sprint("..", string(filepath.Separator), "configs", string(filepath.Separator), "collections.xml"),
// 		fmt.Sprint("..", string(filepath.Separator), "configs", string(filepath.Separator), "db.xml"),
// 		fmt.Sprint("..", string(filepath.Separator), "configs", string(filepath.Separator), "server.xml"),
// 	},
// 	"online": []string{
// 		fmt.Sprint(".", string(filepath.Separator), "configs", string(filepath.Separator), "collections.xml"),
// 		fmt.Sprint(".", string(filepath.Separator), "configs", string(filepath.Separator), "db.xml"),
// 		fmt.Sprint(".", string(filepath.Separator), "configs", string(filepath.Separator), "server.xml"),
// 	},
// }

// func Init() {
// 	str := strings.ToLower(os.Getenv("RepliteWebEnvironment"))
// 	if str == "" {
// 		str = "develop"
// 	}
// 	log.Printf("当前开发环境为:%s", str)
// 	values := environmentMap[str]
// 	DEFAULT_COLLECTIONS_CONFIG = values[0]
// 	DEFAULT_DB_CONFIG = values[1]
// 	DEFUALT_SERVER_CONFIG_FILE = values[2]
// }
