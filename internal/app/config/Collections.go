package config

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"unsafe"
)

/*

  isolate  mongo database and server struct for mongo collection

  the configs/collections represent the relation for server data struct and mongo collection struct

*/

var collectionConfig *staticMap

func GetCollectionConfig() *staticMap {
	return collectionConfig
}

type Collections struct {
	XMLName xml.Name     `xml:"mongoTables"`
	Table   []Collection `xml:"table"`
}

type Collection struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

var DEFAULT_COLLECTIONS_CONFIG = fmt.Sprint(".", string(filepath.Separator), "configs", string(filepath.Separator), "collections.xml")

// var DEFAULT_COLLECTIONS_CONFIG = "configs/collections.xml"

func init() {
	file, err := os.Open(DEFAULT_COLLECTIONS_CONFIG)
	if err != nil {
		panic(fmt.Sprintf("cant find the collections config %s error: %s", DEFAULT_COLLECTIONS_CONFIG, err.Error()))
	}
	var cols = new(Collections)
	cols.Table = make([]Collection, 0, 5)
	err = xml.NewDecoder(bufio.NewReader(file)).Decode(&cols)
	// result, err := utils.XMLToMap(file)
	if err != nil {
		panic(fmt.Sprintf("解析collections配置文件(%s)出错:%s", DEFAULT_COLLECTIONS_CONFIG, err.Error()))
	}
	var result = make(map[string]string)
	for i := 0; i < len(cols.Table); i++ {
		col := cols.Table[i]
		result[col.Name] = col.Value
	}
	//TODO3 how to convert the map[string]string params to map[string]any
	collectionConfig = newStaticMap(result)
}

type staticMap struct {
	data map[string]string
	flag uintptr
}

func (staticMap *staticMap) Get(key string) (value any) {
	value = staticMap.data[key]
	return value
}

func (staticMap *staticMap) TryGet(key string) (value any, ok bool) {
	staticMap.nocopy()
	value, ok = staticMap.data[key]
	return value, ok
}

func (staticmap *staticMap) nocopy() {
	address := unsafe.Pointer(staticmap)
	if !atomic.CompareAndSwapUintptr(&staticmap.flag, uintptr(0), uintptr(address)) && uintptr(address) != staticmap.flag {
		panic("static map cant copy")
	}
}

func newStaticMap(mp map[string]string) *staticMap {
	// m := new(staticMap)
	// copyMap(&m.data, mp)
	m := new(staticMap)
	m.data = copyMap(mp)
	return m
}
func copyMap(src map[string]string) (dst map[string]string) {
	// if dst == nil {
	dst = make(map[string]string)
	// }
	for key, value := range src {
		dst[key] = value
	}
	return dst
}
