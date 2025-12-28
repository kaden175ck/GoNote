package io_test

import (
	"dqq/go/basic/io"
	"fmt"
	"testing"
	"time"
)

func TestWriteFile(t *testing.T) {
	io.WriteFile()
}

func TestReadFile(t *testing.T) {
	io.ReadFile()
}

func TestWriteFileWithBuffer(t *testing.T) {
	io.WriteFileWithBuffer()
}

func TestReadFileWithBuffer(t *testing.T) {
	io.ReadFileWithBuffer()
}

func TestSplitFile(t *testing.T) {
	imgFile := "../img/logo.png"
	io.SplitFile(imgFile, "../img/图像分割", 4)
}

func TestMergeFile(t *testing.T) {
	io.MergeFile("../img/图像分割", "../img/图像合并.png")
}

func TestCreateFile(t *testing.T) {
	io.CreateFile("../data/poem.txt")
}

func TestWalkDir(t *testing.T) {
	io.WalkDir("../data")
}

func TestCompress(t *testing.T) {
	io.Compress("../img/logo.png", "../img/logo.png.gzip", io.GZIP)
	io.Decompress("../img/logo.png.gzip", "../data/logo.png", io.GZIP)
}

func TestLog(t *testing.T) {
	logger := io.NewLogger("../data/biz.log")
	io.Log(logger)
}

func TestSLog(t *testing.T) {
	logger := io.NewSLogger("../data/sbiz.log")
	io.SLog(logger)
}

func TestSysCall(t *testing.T) {
	io.SysCall()
}

func TestJson(t *testing.T) {
	io.JsonSerialize()
}

func TestRegex(t *testing.T) {
	io.UseRegex()
}

func TestLimitReader(t *testing.T) {
	io.LimitReader()
}

func TestMultiReader(t *testing.T) {
	io.MultiReader()
}

func TestMultiWriter(t *testing.T) {
	io.MultiWriter()
}

func TestTeeReader(t *testing.T) {
	io.TeeReader()
}

func TestPipe(t *testing.T) {
	io.PipeIO()
}

func TestBufferedFileWriter(t *testing.T) {
	t1 := time.Now()
	io.WriteDirect("../data/no_buffer.txt")
	t2 := time.Now()
	io.WriteWithBuffer("../data/with_buffer.txt")
	t3 := time.Now()
	fmt.Printf("不用缓冲耗时%dms，用缓冲耗时%dms\n", t2.Sub(t1).Milliseconds(), t3.Sub(t2).Milliseconds())
}

// go test -v ./io -run=^TestWriteFile$ -count=1
// go test -v ./io -run=^TestReadFile$ -count=1
// go test -v ./io -run=^TestWriteFileWithBuffer$ -count=1
// go test -v ./io -run=^TestReadFileWithBuffer$ -count=1
// go test -v ./io -run=^TestCreateFile$ -count=1
// go test -v ./io -run=^TestWalkDir$ -count=1
// go test -v ./io -run=^TestSplitFile$ -count=1
// go test -v ./io -run=^TestMergeFile$ -count=1
// go test -v ./io -run=^TestJson$ -count=1
// go test -v ./io -run=^TestCompress$ -count=1
// go test -v ./io -run=^TestLog$ -count=1
// go test -v ./io -run=^TestSLog$ -count=1
// go test -v ./io -run=^TestSysCall$ -count=1
// go test -v ./io -run=^TestRegex$ -count=1
// go test -v ./io -run=^TestLimitReader$ -count=1
// go test -v ./io -run=^TestMultiReader$ -count=1
// go test -v ./io -run=^TestMultiWriter$ -count=1
// go test -v ./io -run=^TestTeeReader$ -count=1
// go test -v ./io -run=^TestPipe$ -count=1
// go test -v ./io -run=^TestBufferedFileWriter$ -count=1
