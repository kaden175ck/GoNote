package main

import (
	"crypto/sha256"
	"dqq/go/basic/http/live/util"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	fileDir = "http/live/client/"
)

func GetFileSize(FileName string) int {
	resp, err := http.Get("http://127.0.0.1:5678/file_size/" + FileName)
	if err != nil {
		log.Printf("GetFileSize failed: %s", err)
		return -1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, resp.Body)
		return -1
	}

	bs, _ := io.ReadAll(resp.Body)
	FileSize, err := strconv.Atoi(string(bs))
	if err != nil {
		log.Printf("返回的文件大小不是纯数字：%s", string(bs))
		return -1
	}
	return FileSize
}

func Download(FileName string) {
	FileSize := GetFileSize(FileName)
	if FileSize <= 0 {
		return
	}
	log.Printf("文件总大小是 %d B\n", FileSize)

	file, err := os.OpenFile(fileDir+FileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("open file failed: %s", err)
		return
	}
	defer file.Close()

	resp, err := http.Get("http://127.0.0.1:5678/download/" + FileName)
	if err != nil {
		log.Printf("download file failed: %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, resp.Body)
		return
	}

	// 从响应头里读出文件哈希值
	hash := resp.Header.Get("hash")

	bar := util.Bar{}
	bar.NewOptionWithGraph(0, int64(FileSize), "#") // 初始化进度条
	acc := 0                                        // 累积下载的字节数
	buffer := make([]byte, 2048)
	for {
		n, err := resp.Body.Read(buffer)
		if err != nil {
			if err == io.EOF {
				file.Write(buffer[:n])
				acc += n
				bar.Play(int64(acc)) // 刷新进度条
			} else {
				log.Printf("read response body failed: %s", err)
			}
			break
		}
		file.Write(buffer[:n])
		acc += n
		bar.Play(int64(acc)) // 刷新进度条
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println() // 换行

	myHash := Hash4File(fileDir + FileName)
	if myHash != hash {
		log.Printf("文件下载数据有丢失, %s != %s", myHash, hash)
	} else {
		log.Println("文件下载完全没问题", myHash)
	}
}

func main() {
	Download("通天.mp4")
}

func Hash4File(FileName string) string {
	// 打开文件
	file, err := os.Open(FileName)
	if err != nil {
		return ""
	}
	defer file.Close()

	hasher := sha256.New()
	buffer := make([]byte, 2048)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF { // End Of File
				hasher.Write(buffer[:n]) // 响应体
			} else {
				log.Printf("读文件异常: %s", err)
			}
			break
		}
		hasher.Write(buffer[:n]) // 响应体
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

// go run ./http/live/client
