package uc

import "math/rand"

const codeStr = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"

const defaultCodeLength = 6

const codeCacheTTL = 300

// generateSmsCode 生成验证码
func generateSmsCode(length ...int) string {
	codeLength := defaultCodeLength
	if len(length) > 0 {
		codeLength = length[0]
	}
	str := ""
	for i := 0; i < codeLength; i++ {
		randIndex := rand.Intn(36)
		str += string(codeStr[randIndex])
	}
	return str
}
