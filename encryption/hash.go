package encryption

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func Sha1(data string) []byte {
	sha1 := sha1.New()
	sha1.Write([]byte(data))
	return sha1.Sum(nil)
}

func Md5(data string) []byte {
	md5 := md5.New()
	md5.Write([]byte(data))
	return md5.Sum(nil)
}

func CreateSha256OfSmallFile(FileName string) (string, error) {
	file, err := os.Open(FileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file) // 文件较小，可一次性读取文件的全部内容
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write(data)
	digest := hasher.Sum(nil)

	return hex.EncodeToString(digest), nil
}

func CreateSha256OfBigFile(FileName string, BlockSize int) (string, error) {
	file, err := os.Open(FileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()

	data := make([]byte, BlockSize)
	for {
		readBytes, err := file.Read(data)
		if err != nil {
			if err == io.EOF {
				if readBytes > 0 {
					hasher.Write(data[:readBytes])
				}
				break
			} else {
				return "", err
			}
		}
		hasher.Write(data[:readBytes]) // 不断往hasher里写入新内容，但hasher的内存开销是固定的（结构体里都是定长的数组），所以整体的内容开销跟文件大小没关系
	}
	digest := hasher.Sum(nil)

	return hex.EncodeToString(digest), nil
}
