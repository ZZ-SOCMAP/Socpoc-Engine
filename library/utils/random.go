package utils

import (
	"math/rand"
	"time"
)

const (
	indexbits = 6
	indexmask = 1<<indexbits - 1
	indexmax  = 63 / indexbits
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	letternum = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// RandomStr 随机字符串
func RandomStr(n int, choices string) string {
	b := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i, cache, remain := n-1, r.Int63(), indexmax; i >= 0; {
		if remain == 0 {
			cache, remain = r.Int63(), indexmax
		}
		if idx := int(cache & indexmask); idx < len(choices) {
			b[i] = choices[idx]
			i--
		}
		cache >>= indexbits
		remain--
	}
	return string(b)
}

// RandLetters 随机小写字母
func RandLetters(n int) string {
	return RandomStr(n, lowercase)
}

// RandLetterNumbers 随机大小写字母和数字
func RandLetterNumbers(n int) string {
	return RandomStr(n, letternum)
}
