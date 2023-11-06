package utils

import (
	"strings"
)

// MergeStr: if the value[i] is "" , the value will be replace for "*"
func MergeStr(value ...string) string {
	sb := strings.Builder{}
	for i := 0; i < len(value); i++ {
		str := value[i]
		if str == "" {
			str = "*"
		}
		sb.WriteString(str)
	}
	return sb.String()
}

func IgnoreQuotationMarks(str string) string {
	return strings.ReplaceAll(str, "\"", "")
}

func ParseLabel(labels []string) []string {
	var labelSs []string
	for _, label := range labels {
		label = strings.Replace(label, " ", "|", -1)
		label = strings.Replace(label, ",", "|", -1)
		label = strings.Replace(label, "，", "|", -1)
		label = strings.Replace(label, "。", "|", -1)
		label = strings.Replace(label, ".", "|", -1)
		labelSs = append(labelSs, label)
	}
	return labelSs
}
