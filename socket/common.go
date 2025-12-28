package transport

import (
	"log"
	"os"
)

func CheckError(err error) {
	if err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}
}

var MAGIC = []byte{13, 7, 29, 83, 113} //一般取素数

type AddRequest struct {
	RequestId int
	A         int
	B         int
}

type AddResponse struct {
	RequestId int
	Sum       int
}
