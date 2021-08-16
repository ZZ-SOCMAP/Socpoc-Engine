package utils

import (
	"sort"
	"unsafe"
)

// SortMapKeys map keys to string array and sort
func SortMapKeys(m map[string]string) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// StrToBytes string to bytes
func StrToBytes(s string) []byte {
	t := (*[2]uintptr)(unsafe.Pointer(&s))
	b := [3]uintptr{t[0], t[1], t[1]}
	return *(*[]byte)(unsafe.Pointer(&b))
}

// BytesToStr bytes to string
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
