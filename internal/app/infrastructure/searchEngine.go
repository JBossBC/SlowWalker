package infrastructure

import (
	"encoding/json"
	"replite_web/internal/app/config"
	"replite_web/internal/app/utils"

	"io"
	"os"

	"github.com/meilisearch/meilisearch-go"
)

type MeiliSearch struct {
}

func getMeiliSearchClient() *MeiliSearch {
	return new(MeiliSearch)
}

// to get meilisearchClient
func (m MeiliSearch) getMeiliSearchClient() *meilisearch.Client {
	return meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   string(config.GetServerConfig().MeiliSearch.Host), //server address
		APIKey: string(config.GetServerConfig().MeiliSearch.Key),  //API key
	})
}

// add Documents to sure index
func (m MeiliSearch) AddDocuments(documentsFile string, indexName string) error {
	jsonFile, _ := os.Open(documentsFile)
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var documents []map[string]interface{}
	err := json.Unmarshal(byteValue, &documents)
	if err != nil {
		return err
	}
	client := m.getMeiliSearchClient()
	_, err = client.Index(indexName).AddDocuments(documents)
	if err != nil {
		panic(err)
	}
	return err
}

// TODO should set the two field to help meilisearch building the index which is helpful to searching('label' and 'description')
//func (m MeiliSearch) SearchDocuments2(label []string, description string) (response []interface{}, err error) {
//	client := m.getMeiliSearchClient()
//	labels := queryStrings(label)
//	searchReq := &meilisearch.SearchRequest{
//		AttributesToSearchOn: []string{"label", "description"}, //just search in the two fields
//		Filter: [][]string{
//			labels,
//		},
//	}
//	resp, err := client.Index("WebFunction").Search(description, searchReq)
//	fmt.Println(resp.Hits)
//	return resp.Hits, err
//}

// TODO should set the two field to help meilisearch building the index which is helpful to searching('label' and 'description')
func (m MeiliSearch) SearchDocuments(label []string, description string) (response utils.Response) {
	client := m.getMeiliSearchClient()
	labels := queryStrings(label)
	searchReq := &meilisearch.SearchRequest{
		AttributesToSearchOn: []string{"label", "description"}, //just search in the two fields
		Filter: [][]string{
			labels,
		},
	}
	resp, err := client.Index("WebFunction").Search(description, searchReq)
	if err != nil {
		return utils.NewFailedResponse("搜索失败！")
	}
	return utils.NewSuccessResponse(resp)
}

// format the labels
func queryStrings(labels []string) []string {
	var queryStrings []string
	for _, label := range labels {
		queryString := "label = \"" + label + "\""
		queryStrings = append(queryStrings, queryString)
	}
	return queryStrings //return example :[]string{"label = \"golang_scripts\"", "label = \"python_scripts\""}
}

// delete index involve all documents belong the index
func (m MeiliSearch) DeleteAllDocuments(indexName string) error {
	client := m.getMeiliSearchClient()
	_, err := client.DeleteIndex(indexName)
	return err
}

// setting fliters
func (m MeiliSearch) UpdateSettingFilters() error {
	settings := meilisearch.Settings{
		FilterableAttributes: []string{"label"},
	}
	client := m.getMeiliSearchClient()
	_, err := client.Index("WebFunction").UpdateSettings(&settings)
	if err != nil {
		return err
	}
	return nil
}
