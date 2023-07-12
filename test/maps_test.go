package test

import (
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
	"testing"
)

func TestMapsFormat(t *testing.T) {
	result := utils.Format(dao.Rule{Name: "xiyang", Authority: "xiyang"})
	if result == nil {
		t.Fatal("format model error ")
	}
	for k, v := range result {
		println(k, v)
	}
}
