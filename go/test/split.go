package test

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

	fmt.Sprint(" this is %v ", "123")
	var name interface{} = "yinzhengjie"
	fmt.Printf("My name is %v !\n", name)
	var age interface{} = 18
	fmt.Printf("I am [%d] years old。", age)

	test()
}

func test() {
	a := "123.html"
	b := "book/"

}
