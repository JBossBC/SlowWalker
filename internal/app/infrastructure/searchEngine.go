package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"replite_web/internal/app/config"
	"replite_web/internal/app/utils"
	"sync"

	"io"
	"os"

	"github.com/meilisearch/meilisearch-go"
)

type MeiliSearchProvider struct {
	meiliSearchClient *meilisearch.Client
}

var (
	meiliSearchProvider *MeiliSearchProvider
	meiliSearchOnce     sync.Once
)

func getMeiliSearchProvider() *MeiliSearchProvider {
	meiliSearchOnce.Do(func() {
		meiliSearchProvider = new(MeiliSearchProvider)
		meiliSearchProvider.meiliSearchClient = getMeiliSearchClient()
	})
	return meiliSearchProvider
}

func getMeiliSearchClient() *meilisearch.Client {
	return meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   config.GetServerConfig().MeiliSearch.Host, //server address
		APIKey: config.GetServerConfig().MeiliSearch.Key,  //API key
		//Host:   "http://localhost:7700",
		//APIKey: "123456789",
	})
}

// add Documents to sure index
func (m MeiliSearchProvider) AddDocuments(documentsFile string, indexName string) error {
	jsonFile, _ := os.Open(documentsFile)
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var documents []map[string]interface{}
	err := json.Unmarshal(byteValue, &documents)
	if err != nil {
		return err
	}
	_, err = m.meiliSearchClient.Index(indexName).AddDocuments(documents)
	if err != nil {
		log.Println(err)
	}
	return err
}

// TODO should set the two field to help meilisearch building the index which is helpful to searching('label' and 'description')
func (m MeiliSearchProvider) SearchFunctions(label []string, description string, index string) (response utils.Response) {
	searchReq := &meilisearch.SearchRequest{
		AttributesToSearchOn: []string{"label", "description", "title"}, //just search in the two fields
	}
	resp, err := m.meiliSearchClient.Index(index).Search(description, searchReq)
	regStr := string(".*" + queryStrings(label) + ".*")
	pattern := regexp.MustCompile(regStr) // 定义正则匹配模式
	filteredHits := make([]interface{}, 0)
	for _, hit := range resp.Hits {
		if label, ok := hit.(map[string]interface{})["label"].(string); ok {
			if pattern.MatchString(label) {
				filteredHits = append(filteredHits, hit)
			}
		}
	}
	if err != nil {
		log.Println(err)
		return utils.NewFailedResponse("搜索失败！")
	}
	return utils.NewSuccessResponse(filteredHits)
}

// format the labels
func queryStrings(labels []string) string {
	var strS string
	for _, label := range labels {
		if len(labels) == 1 {
			strS = label // 在每个标签之后添加空格
			return strS
		}
		strS += label + "|"
	}
	return strS
}

// delete index involve all documents belong the index
func (m MeiliSearchProvider) DeleteAllDocuments(indexName string) error {
	_, err := m.meiliSearchClient.DeleteIndex(indexName)
	return err
}

// setting fliters
func (m MeiliSearchProvider) UpdateSettingFilters(index string) error {
	settings := meilisearch.Settings{
		FilterableAttributes: []string{"label"},
	}
	_, err := m.meiliSearchClient.Index(index).UpdateSettings(&settings)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	meiliConfig := meilisearch.ClientConfig{
		Host:   config.GetServerConfig().MeiliSearch.Host, //server address
		APIKey: config.GetServerConfig().MeiliSearch.Key,  //API key
		//Host:   "http://localhost:7700",
		//APIKey: "123456789",
	}
	client := meilisearch.NewClient(meiliConfig)
	_, err := client.Health()
	if err != nil {
		panic(fmt.Sprintf("connecting the MeiliSearch within %s seconds is error: %v", meiliConfig.Timeout, err))
	}
	log.Printf("Connected to MeiliSearch")
}
