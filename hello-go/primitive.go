package main

import "fmt"

func main() {
	var MyName int
	fmt.Println(MyName)

	var a int = 10
	var b = a

	_ = b
	c := b
	a = c

}
