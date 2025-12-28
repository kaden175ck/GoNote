package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

type Comparable interface {
	~int32 | int64
}

type Apple[T cmp.Ordered] struct{}

func (Apple[T]) getBigger(a, b T) T {
	if a > b {
		return a
	} else {
		return b
	}
}

type GetUserRequest struct{}
type GetBookRequest struct{}

func httpRPC[T GetUserRequest | GetBookRequest](request T) {
	url := "http://127.0.0.1/"
	tp := reflect.TypeOf(request)
	switch tp.Name() {
	case "GetUserRequest":
		url += "user"
	case "GetBookRequest":
		url += "book"
		// default:
		// 	panic("不支持的request类型")
	}
	fmt.Println(url)
	bs, _ := json.Marshal(request)
	http.Post(url, "application/json", bytes.NewReader(bs))
}

func main33() {
	httpRPC(GetUserRequest{})

	a := Apple[int32]{}
	a.getBigger(3, 6)
}
