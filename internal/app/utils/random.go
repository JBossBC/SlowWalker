package utils

import (
	"math/rand"
	"strconv"
	"time"
)

var DEFAULT_CODE_LENGTH = 1000000

// TODO 目前只能做到纳秒级随机
func NewRandomCode() string {
	// var randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
	// sb := strings.Builder{}
	// for i := 0; i < DEFAULT_CODE_LENGTH; i++ {
	// 	sb.WriteString(strconv.FormatInt(int64(randSource.Intn(10)), 10))
	// }
	rand.Seed(time.Now().UnixNano())
	return strconv.FormatInt(int64(rand.Intn(DEFAULT_CODE_LENGTH)), 10)
}
