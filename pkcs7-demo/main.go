package main

import (
	"fmt"
)

// PKCS7Padding 对数据进行 PKCS7 填充
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := 0; i < padding; i++ {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

// PKCS7UnPadding 去除 PKCS7 填充
func PKCS7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("数据为空")
	}
	padding := int(data[length-1])
	if padding > length {
		return nil, fmt.Errorf("无效的填充")
	}
	return data[:length-padding], nil
}

func main() {
	fmt.Println("=== PKCS7 填充测试 ===")
	blockSize := 16

	// 测试数据切片
	testCases := []string{
		"123456789",
		"src",
		"Hello",
		"1234567890123456",
	}

	for i, testData := range testCases {
		fmt.Printf("\n测试 %d: %s\n", i+1, testData)
		src := []byte(testData)

		// 填充
		padded := PKCS7Padding(src, blockSize)
		fmt.Printf("原始: %s (长度: %d)\n", src, len(src))
		fmt.Printf("填充: %v (长度: %d)\n", padded, len(padded))

		// 去填充
		unpadded, _ := PKCS7UnPadding(padded)
		fmt.Printf("还原: %s\n", unpadded)
	}
}
