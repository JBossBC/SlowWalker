package utils

import (
	"fmt"
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
		fmt.Printf("输入的label为%v\n", label)
		label = strings.Replace(label, " ", "|", -1)
		label = strings.Replace(label, ",", "|", -1)
		label = strings.Replace(label, "，", "|", -1)
		label = strings.Replace(label, "。", "|", -1)
		label = strings.Replace(label, ".", "|", -1)
		labelSs = append(labelSs, label)
	}
	fmt.Printf("输出的labels为%v\n", labelSs)
	return labelSs
}
