package test

import (
	"fmt"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
	"testing"
)

func TestIP(t *testing.T) {

	fmt.Println(dao.GetIPDao().QueryIP("127.0.0.1")) //输出：3，有几条记录就会输出几

	fmt.Println(utils.MergeStr("register-failed-", "127.0.0.1")) //输出：register-failed-127.0.0.1

	dao.GetIPDao().InsertIP("127.0.0.1")

}
