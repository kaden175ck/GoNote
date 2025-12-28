package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	fileDir = "http/live/server/"
)

func GetFileSize(w http.ResponseWriter, r *http.Request) {
	// 获取请求参数
	FileName := r.PathValue("file_name")
	// 检查参数的合法性(安全性)
	if strings.Contains(FileName, "/") {
		http.Error(w, "非法的文件名", http.StatusBadRequest)
		return
	}
	// 打开文件
	file, err := os.Open(fileDir + FileName)
	if err != nil {
		http.Error(w, "文件不存在", http.StatusNotFound)
		return
	}
	defer file.Close()

	//获取文件的基本信息
	stat, _ := file.Stat()
	FileSize := stat.Size() // 文件的大小(B)
	// 响应头
	//w.Header().Add("file-size", strconv.FormatInt(FileSize, 10))
	// 响应体
	w.Write([]byte(strconv.FormatInt(FileSize, 10)))
	//w.WriteHeader(http.StatusOK)
}

func Download(w http.ResponseWriter, r *http.Request) {
	// 获取请求参数
	FileName := r.PathValue("file_name")
	// 检查参数的合法性(安全性)
	if strings.Contains(FileName, "/") {
		http.Error(w, "非法的文件名", http.StatusBadRequest)
		return
	}

	hash := Hash4File(fileDir + FileName)
	w.Header().Add("hash", hash) // 把文件的哈希值放到响应头

	// 打开文件
	file, err := os.Open(fileDir + FileName)
	if err != nil {
		http.Error(w, "文件不存在", http.StatusNotFound)
		return
	}
	defer file.Close()

	buffer := make([]byte, 2048)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF { // End Of File
				w.Write(buffer[:n]) // 响应体
			} else {
				log.Printf("读文件异常: %s", err)
			}
			break
		}
		w.Write(buffer[:n]) // 响应体
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /file_size/{file_name}", GetFileSize)
	mux.HandleFunc("GET /download/{file_name}", Download)
	if err := http.ListenAndServe("127.0.0.1:5678", mux); err != nil {
		panic(err)
	}
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

// go run ./http/live/server
