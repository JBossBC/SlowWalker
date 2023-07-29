package infrastructure

type Operate interface {
	// Valid() bool
}

type BaseOperate struct {
	Operator string         `json:"operator"`
	Params   map[string]any `json:"params"`
	Function string         `json:"function"`
	// websocket connection callback func, make sure the message be consumed
	CallBack func(any) any `json:"-"`
}

// func (base *BaseOperate) Valid() bool {
// 	return true
// }

type OperateOption func(*BaseOperate)

func WithParams(params map[string]any) OperateOption {
	return func(op *BaseOperate) {
		op.Params = params
	}
}

func WithCallBack(callback func(any) any) OperateOption {
	return func(op *BaseOperate) {
		op.CallBack = callback
	}
}

func NewOperate(operator string, function string, options ...OperateOption) BaseOperate {
	base := new(BaseOperate)
	for i := 0; i < len(options); i++ {
		options[i](base)
	}
	base.Operator = operator
	base.Function = function
	return *base
}

// type LinuxOperate struct {
// 	base *BaseOperate
// }

// func (linux *LinuxOperate) Valid() bool {

// }

// func NewLinuxOperate(operator string, function string, options ...OperateOption) LinuxOperate {
// 	op := new(LinuxOperate)
// 	base := newBaseOperate(operator, function)
// 	op.base = &base
// 	for i := 0; i < len(options); i++ {
// 		option := options[i]
// 		option(op.base)
// 	}
// 	return *op
// }

// type WindowsOperator struct {
// }
