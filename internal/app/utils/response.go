package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var HELPLESS_ERROR_RESPONSE []byte

func init() {
	immutableError, err := json.Marshal(NewFailedResponse("系统错误"))
	if err != nil {
		panic(fmt.Sprintf("尝试序列化response失败:%s", err.Error()))
	}
	HELPLESS_ERROR_RESPONSE = immutableError

}

type Response interface {
	SerializeJSON() []byte
}

type BaseResponse struct {
	State   bool   `json:"state"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func (response BaseResponse) SerializeJSON() []byte {
	result, err := json.Marshal(response)
	if err != nil {
		return HELPLESS_ERROR_RESPONSE
	}
	return result
}
func NewSuccessResponse(data any) Response {
	return BaseResponse{
		State:   true,
		Data:    data,
		Message: "成功",
	}
}

func NewFailedResponse(message string) Response {
	return BaseResponse{
		State:   false,
		Data:    nil,
		Message: message,
	}
}

func Not_Granted_Error(request *http.Request) []byte {
	log.Printf("禁止访问 %v", request)
	response := NewFailedResponse("你没有权限访问")
	data := response.SerializeJSON()
	return data
}
