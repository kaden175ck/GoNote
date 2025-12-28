package main

import (
	"fmt"
	"strings"
)

func DifferentNil() {
	var slc1 []int = nil
	var slc2 []bool = nil
	var mp map[int]bool = nil
	var ch chan int = nil
	var fk func(int) bool = nil
	var p *int = nil
	var err error = nil
	// 以上这些nil，全都互不相等

	// fmt.Println(slc1 == slc2) // 编译错误，虽然都是slice，但类型还是不同，不能进行==比较
	// fmt.Println(slc1 == mp)   // 编译错误，类型不同，不能进行==比较

	var ifc1 any
	var ifc2 any
	ifc1 = slc1
	ifc2 = slc2
	fmt.Println(ifc1 == ifc2)
	ifc1 = slc2
	ifc2 = mp
	fmt.Println(ifc1 == ifc2)
	ifc1 = mp
	ifc2 = ch
	fmt.Println(ifc1 == ifc2)
	ifc1 = ch
	ifc2 = fk
	fmt.Println(ifc1 == ifc2)
	ifc1 = fk
	ifc2 = p
	fmt.Println(ifc1 == ifc2)
	ifc1 = p
	ifc2 = err
	fmt.Println(ifc1 == ifc2)
}

type TP int

// TP没有实现error接口，但是*TP实现了error接口
func (tp *TP) Error() string {
	return "TP"
}

func NewTP() *TP {
	return nil
}

func main() {
	DifferentNil()
	fmt.Println(strings.Repeat("*", 20))

	var tp *TP = nil
	fmt.Println(tp.Error()) // 可以调用nil pointer的方法
	fmt.Println(tp == nil)  // true
	var err error
	//fmt.Println(err.Error()) // err是nil interface，调用其Error()方法会导致运行时错误

	err = tp                 // err绑定了具体的类型，err不再是nil interface
	fmt.Println(err.Error()) // 因为err不是nil interface，所以调用其Error()方法是有输出的
	fmt.Println(err == nil)  // false

	err = NewTP()            // err绑定了具体的类型，err不再是nil interface
	fmt.Println(err.Error()) // 因为err不是nil interface，所以调用其Error()方法是有输出的
	fmt.Println(err == nil)  // false

	var p *int = nil
	if p == nil {
		fmt.Println("p是nil pointer")
	}
	var a interface{} = p // a绑定了具体的类型，a不再是nil interface
	if a != nil { // True
		fmt.Println("a不是nil interface")
	}
}
