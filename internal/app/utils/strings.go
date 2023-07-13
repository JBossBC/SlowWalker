package utils

import "strings"

// MergeStr: if the value[i] is "" , the value will be replace for "*"
func MergeStr(value ...string) string {
	sb := strings.Builder{}
	for i := 0; i < len(value); i++ {
		str := value[i]
		if str == "" {
			str = "*"
		}
		sb.WriteString(value[i])
	}
	return sb.String()
}

func IgnoreQuotationMarks(str string) string {
	return strings.ReplaceAll(str, "\"", "")
}
