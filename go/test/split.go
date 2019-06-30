package main

import (
	"fmt"
	"strings"
)

func main() {
	var str = "https://www.qu.la/book/2639/1460031.html https://www.qu.la/book/2639/1460032.html"

	var next = "123.html"

	splitStr := strings.Split(str, "/")
	fmt.Println("splitStr", splitStr)
	splitStr[len(splitStr)-1] = next
	fmt.Println("splitStr2", splitStr)

	newStr := strings.Join(splitStr, "/")
	fmt.Println("newStr", newStr)

}
