package dao

type OperateType int

const (
	// represent the task should be running when it is created
	Sustainment OperateType = iota
	//  represent the task will be close when it runs at a moment
	ShortTerm
)

type Operate interface {
	// Valid() bool
	GetFunction() string
	GetOperator() string
	GetParams() []string
	// GetCallBack() func(any) any
	GetOperateType() OperateType
	GetMedium() []string
}

type BaseOperate struct {
	Operate
	// Command  string   `json:"command"`
	Operator string `json:"operator"`
	// the complete params including the params template and actual params for the inputing of user
	Params []string `json:"params"`
	//operate function according to the rule collection
	Function string `json:"function"`
	// // websocket connection callback func, make sure the message be consumed
	// CallBack func(any) `json:"-"`
	// the operateType
	OperateType OperateType `json:"operatetype"`
	// the share resource for need to handler
	Medium []string `json:"medium"`
}

func (base BaseOperate) GetOperateType() OperateType {
	return base.OperateType
}
func (base BaseOperate) GetMedium() []string {
	return base.Medium
}

func (base BaseOperate) GetFunction() string {
	return base.Function
}
func (base BaseOperate) GetOperator() string {
	return base.Operator
}
func (base BaseOperate) GetParams() []string {
	return base.Params
}

// func (base *BaseOperate) GetCallBack() func(any) any {
// 	return base.CallBack
// }

// func (base *BaseOperate) Valid() bool {
// 	return true
// }

type OperateOption func(*BaseOperate)

func WithParams(params []string) OperateOption {
	return func(op *BaseOperate) {
		op.Params = params
	}
}

func WithMedium(medium []string) OperateOption {
	return func(op *BaseOperate) {
		op.Medium = medium
	}
}

// ignore the callback
// func WithCallBack(callback func(any) any) OperateOption {
// 	return func(op *BaseOperate) {
// 		op.CallBack = callback
// 	}
// }

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
