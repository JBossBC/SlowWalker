package test

import (
	"replite_web/internal/app/utils"
	"strings"
	"testing"
)

func TestIgnoreQM(t *testing.T) {
	result := utils.IgnoreQuotationMarks("\"xiyang\"")
	if strings.Compare(result, "xiyang") != 0 {
		t.Errorf("测试失败:%s", result)
	}
}
