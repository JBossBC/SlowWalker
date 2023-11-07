package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
)

const defaultChatgptURL = "https://api.openai-proxy.com/v1/chat/completions"

type chatgpt struct {
	chatgptClient *http.Client
}

var (
	chatgptModel *chatgpt
	chatgptOnce  sync.Once
)

type Params struct {
	Key               string  `json:"-"`
	Model             string  `json:"model"`
	Max_tokens        int     `json:"max_tokens"`
	Temperature       float64 `json:"temperature"`
	Top_p             float64 `json:"top_p"`
	Presence_penalty  int     `json:"presence_penalty"`
	Frequency_penalty int     `json:"frequency_penalty"`
	N                 int     `json:"n"`
	Stream            bool    `json:"-"`
	Stop              string  `json:"-"`
	// Logit_bias        map[string]interface{} `gorm:"type:json" json:"logit_bias"`
	Logit_bias map[string]interface{} `json:"-"`
}

// init for database
// current can use default configure
var defaultParams Params = Params{
	Presence_penalty:  0,
	Frequency_penalty: 1,
	Temperature:       0.5,
	Top_p:             0.3,
	Model:             "gpt-3.5-turbo",
	Key:               "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",
	Max_tokens:        1024,
	N:                 1,
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	// Name    string `json:"name"`
}
type ChatResponse struct {
	Id      string        `json:"id"`
	Object  string        `json:"object"`
	Created time.Duration `json:"created"`
	Choices []Choice      `json:"choices"`
	Usage   Usage         `json:"usage"`
}

type Choice struct {
	Index         int64         `json:"index"`
	Message       OpenAIMessage `json:"message"`
	Finish_reason string        `json:"finish_reason"`
}

type Usage struct {
	Prompt_tokens     int64 `json:"prompt_tokens"`
	Completion_tokens int64 `json:"completion_tokens"`
	Total_tokens      int64 `json:"total_tokens"`
}

type Body struct {
	Params Params `json:"-"`
	// Params struct {
	// 	Key               string  `json:"-"`
	// 	Model             string  `json:"model"`
	// 	Max_tokens        int     `json:"max_tokens"`
	// 	Temperature       float64 `json:"temperature"`
	// 	Top_p             float64 `json:"top_p"`
	// 	Presence_penalty  int     `json:"presence_penalty"`
	// 	Frequency_penalty int     `json:"frequency_penalty"`
	// 	N                 int     `json:"n"`
	// 	Stream            bool    `json:"-"`
	// 	Stop              string  `json:"-"`
	// 	// Logit_bias        map[string]interface{} `gorm:"type:json" json:"logit_bias"`
	// 	Logit_bias map[string]interface{} `json:"-"`
	// }   `json:"params,inline"`
	//default keep the 10 length for per session
	Messages []OpenAIMessage `json:"messages"`
}

func getChatgptModel() *chatgpt {
	chatgptOnce.Do(func() {
		chatgptModel = new(chatgpt)
		chatgptModel.chatgptClient = new(http.Client)
	})
	return chatgptModel
}

var openAIMessageTemplate = []OpenAIMessage{
	OpenAIMessage{Role: "system", Content: "你现在是一名前端需求分析师"},
	OpenAIMessage{Role: "user", Content: "``` ```之中的内容代表客户提出的需求，请你将客户的需求用HTML表单内元素表达出来,因为客户可能表达多个HTML表单元素,所以最终的结果应该是JSON格式的数组,其中每个元素回答的模板应该是: element:xxx,type:xxx,title:xxx,value:xxx(多个值应该用\\符号分割),元素名只能有以下取值,()内代表一个取值,(input)(select)(textarea)(datalist)(fieldset),type只有当元素名为input的时候才会有值,分别是(file)(text)(del)(date)(radio)(email)(checkbox)(password)(number)(time)(week)(datetime)(search)(color) 比如说有例子如下客户:````我现在需要有一个可以传输文件的功能,名字是:IP文件,用户需要传入一个文件```,你的回答应该为: [{\"element\":\"input\",\"type\":\"file\",\"title\":\"IP文件\",\"value\":\"\"}]。客户:```用户可以选择是python还是Java或者是Golang,名字是功能类型，用户只能选择这三个中的其中一个```,你的回答应该是[{\"element\":\"input\",\"type\":\"radio\",\"title\":\"功能类型\",\"value\":\"python/Java/Golang\"}]。 你现在与客户正在交流%s的需求,客户:```%s```"},
}

func (chatgpt *chatgpt) Guess(title string, prompt string) (components []*ComponentInfo, err error) {
	defer func() {
		if panicError := recover(); panicError != nil {
			log.Printf("调用chatgpt服务出现系统级错误:%v\n", panicError)
			err = fmt.Errorf("系统出错")
		}
	}()
	content, err := chatgpt.sendPrompt(title, prompt)
	if err != nil {
		log.Println("调用chatgpt服务出错:", err)
		return nil, err
	}
	components = make([]*ComponentInfo, 0, 3)
	// err = json.Unmarshal([]byte(content), components)
	// if err != nil {
	// 	return nil, err
	// }
	err = parseJSON(content, components)
	if err != nil {
		log.Printf("解析chatgpt content出错:%s\n", content)
		return nil, err
	}
	return components, nil
}

func (chatgpt *chatgpt) sendPrompt(title string, prompt string) (content string, err error) {
	body := new(Body)
	body.Params = defaultParams
	body.Messages = make([]OpenAIMessage, len(openAIMessageTemplate))
	copy(body.Messages, openAIMessageTemplate)
	body.Messages[1].Content = fmt.Sprintf(body.Messages[1].Content, title, prompt)
	bodyByt, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest(http.MethodPost, defaultChatgptURL, bytes.NewBuffer(bodyByt))
	if err != nil {
		return "", err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", defaultParams.Key))
	response, err := chatgpt.chatgptClient.Do(request)
	if err != nil {
		return "", err
	}
	result, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("%v response read error:%s", result, err.Error())
	}
	var resultMap ChatResponse
	err = json.NewDecoder(bytes.NewReader([]byte(result))).Decode(&resultMap)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(resultMap.Choices[0].Message.Content), nil

}

func parseJSON(text string, components []*ComponentInfo) error {
	pattern := `\[.*?\]`
	re := regexp.MustCompile(pattern)
	match := re.FindString(text)
	if match != "" {
		err := json.Unmarshal([]byte(match), &components)
		if err != nil {
			return fmt.Errorf("无法解析JSON：%v", err)
		}
	}
	return nil
}
