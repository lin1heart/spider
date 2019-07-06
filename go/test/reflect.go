package test

import (
	"fmt"
	"reflect"
)

type User struct {
	id   int
	name string
	age  int
	tip  string
}

func main() {

	var value interface{} = &User{1, "Tom", 12, "nan"}
	v := reflect.ValueOf(value)
	fmt.Println("aaa", v)
	fmt.Println("bbb", value)
}
