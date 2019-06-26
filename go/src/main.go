package main

import "fmt"
import "unsafe"

func main() {
	var a int = 4

	fmt.Println("111")
	println(unsafe.Sizeof(a))
}