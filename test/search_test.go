package test

import (
	"fmt"
	"replite_web/internal/app/infrastructure"
	"testing"
)

func TestMeiliSearch(t *testing.T) {
	meiliProvider := infrastructure.MeiliSearchProvider{}
	err := meiliProvider.AddDocuments("testMeiliSearch.json", "test") //test add
	if err != nil {
		panic(err)
	}
	err = meiliProvider.UpdateSettingFilters("test") //test setting
	if err != nil {
		panic(err)
	}
	rsp := meiliProvider.SearchFunctions([]string{"go"}, "input", "test") //test search
	fmt.Println(rsp)
	err = meiliProvider.DeleteAllDocuments("test") //test delete
	if err != nil {
		panic(err)
	}

}
