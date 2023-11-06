package test

import (
	"replite_web/internal/app/infrastructure"
	"testing"
)

func TestMeiliSearch(t *testing.T) {
	meiliProvider := infrastructure.MeiliSearchProvider{}
	err := meiliProvider.AddDocuments("searchTest.json", "WebFFF") //test add
	if err != nil {
		panic(err)
	}
	err = meiliProvider.UpdateSettingFilters("WebFFF") //test setting
	if err != nil {
		panic(err)
	}
	//rsp := meiliProvider.SearchFunctions([]string{"go"}, "input", "WebFFF") //test search
	//fmt.Println(rsp)
	//err = meiliProvider.DeleteAllDocuments("test") //test delete
	//if err != nil {
	//	panic(err)
	//}

}
