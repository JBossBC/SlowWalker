package infrastructure

type Operate interface {
}

type baseOperate struct {
	Operator string         `json:"operator"`
	Params   map[string]any `json:"params"`
	Command  string         `json:"command"`
	Function string         `json:"function"`
}
