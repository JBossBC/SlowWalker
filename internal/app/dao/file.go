package dao

import (
	"time"
)

const default_file_redis_prefix = "files-"

const default_expire_times = 30 * time.Minute

// func AuthenticateFiles(username string, files []string) error {
// 	err := CreateSet(utils.MergeStr(default_file_redis_prefix, username), files, default_expire_times)
// 	if err != nil {
// 		log.Printf("添加files(username:%s,files:%s)到缓存中失败:%s", username, files, err.Error())
// 	}
// 	return err
// }
