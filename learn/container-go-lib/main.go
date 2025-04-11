package main

import (
	"container/list"
	"fmt"
)

func main() {
	// 使用list.New()直接初始化
	l1 := list.New()
	l1.PushFront(1)
	fmt.Println(l1.Front().Value) // 1

	l1.InsertBefore(2, l1.Front())

	// 使用list.List{}延迟初始化
	l2 := list.List{}
	l2.PushFront(2)
	fmt.Println(l2) // 2
}
