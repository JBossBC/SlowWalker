package utils

import "encoding/json"

var HELPLESS_ERROR_REQUEST = ""

type Request interface {
	MarshalJSON() string
}

type BaseRequest struct {
}

func (baseRequest *BaseRequest) MarshalJSON() string {
	result, err := json.Marshal(baseRequest)
	if err != nil {
		return HELPLESS_ERROR_REQUEST
	}
	return string(result)
}

type SMSRequest struct {
	BaseRequest
	Code string `json:"code"`
}

func (smsRequest *SMSRequest) MarshalJSON() string {
	result, err := json.Marshal(smsRequest)
	if err != nil {
		return HELPLESS_ERROR_REQUEST
	}
	return string(result)
}

func NewSMSRequest(code string) Request {
	return &SMSRequest{
		Code: code,
	}
}
