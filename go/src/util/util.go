package util

import "fmt"

func CheckError(errMasg error) {
	if errMasg != nil {
		fmt.Println("error %s", errMasg)
		panic(errMasg)
	}
}
