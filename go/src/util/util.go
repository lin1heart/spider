package util

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

func CheckError(errMasg error) {
	if errMasg != nil {
		//fmt.Println("error %s", errMasg)
		log.Fatal(errMasg)
		panic(errMasg)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}


func GenerateSqlIdsSuffix(ids []string) string {
	idsStr := strings.Join(ids, ",")

	if idsStr == "" {
		return ""
	}
	sqlSuffix := fmt.Sprintf(" AND id NOT IN ( %s )", idsStr)
	return sqlSuffix
}