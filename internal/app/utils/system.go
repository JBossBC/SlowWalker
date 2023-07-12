package utils

import "unsafe"

func GetOperationBit() int {
	bitSize := unsafe.Sizeof(0) * 8
	return int(bitSize)
}
