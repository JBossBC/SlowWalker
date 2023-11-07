package infrastructure

type AIModel interface {
	Guess(title string, prompt string) ([]*ComponentInfo, error)
}

func GetAIModel() AIModel {
	return getChatgptModel()
}

type ComponentInfo struct {
	Element string `json:"element"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Value   string `json:"value"`
}
