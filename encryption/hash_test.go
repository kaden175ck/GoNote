package encryption_test

import (
	"dqq/go/basic/encryption"
	"fmt"
	"testing"
	"time"
)

var (
	BigFile = "D:\\download\\go1.25.1.windows-amd64.msi"
)

func TestHash(t *testing.T) {
	data := "123456"
	hs := encryption.Sha1(data)
	fmt.Println("SHA-1", hs, len(hs))
	hm := encryption.Md5(data)
	fmt.Println("MD5", hm, len(hm))
}

func TestCreateSha256OfSmallFile(t *testing.T) {
	hash, err := encryption.CreateSha256OfSmallFile(BigFile)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("CreateSha256OfSmallFile", hash)
}

func TestCreateSha256OfBigFile(t *testing.T) {
	begin := time.Now()
	hash, err := encryption.CreateSha256OfBigFile(BigFile, 10<<20)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("CreateSha256OfBigFile", hash, "use time", time.Since(begin).Milliseconds())
}

// go test -v ./encryption -run=^TestHash$ -count=1

// go test -v ./encryption -run=^TestCreateSha256OfSmallFile$ -count=1
// go test -v ./encryption -run=^TestCreateSha256OfBigFile$ -count=1
// certutil -hashfile "D:\\download\\go1.25.1.windows-amd64.msi" SHA256
