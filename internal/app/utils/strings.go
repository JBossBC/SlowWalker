package utils

import "strings"

func MergeStr(value ...string) string {
	sb := strings.Builder{}
	for i := 0; i < len(value); i++ {
		sb.WriteString(value[i])
	}
	return sb.String()
}

func IgnoreQuotationMarks(str string) string {
	return strings.ReplaceAll(str, "\"", "")
}
