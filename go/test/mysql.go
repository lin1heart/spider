package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

import (
	"fmt"
)

func checkErr(errMasg error) {
	if errMasg != nil {
		fmt.Println("error %s", errMasg)
		panic(errMasg)
	}
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(39.104.226.149:3306)/spider")
	checkErr(err)
	fmt.Println("sql succeed  %s", db)

}
