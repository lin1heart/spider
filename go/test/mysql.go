package main

import (
	"database/sql"
	"strconv"
)
import _ "github.com/go-sql-driver/mysql"

import (
	"fmt"
)

func checkErr(errMasg error) {
	if errMasg != nil {
		//fmt.Println("error %s", errMasg)
		//panic(errMasg)
		//log.Fatal(errMasg)
	}
}

func printResult(query *sql.Rows) {
	column, _ := query.Columns()              //读出查询出的列字段名
	values := make([][]byte, len(column))     //values是每个列的值，这里获取到byte里
	scans := make([]interface{}, len(column)) //因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
	for i := range values {                   //让每一行数据都填充到[][]byte里面
		scans[i] = &values[i]
	}
	results := make(map[int]map[string]string) //最后得到的map
	i := 0
	for query.Next() { //循环，让游标往下移动
		if err := query.Scan(scans...); err != nil { //query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			fmt.Println(err)
			return
		}
		row := make(map[string]string) //每行数据
		for k, v := range values {     //每行数据是放在values里面，现在把它挪到row里
			key := column[k]
			row[key] = string(v)
		}
		results[i] = row //装入结果集中
		i++
	}
	for k, v := range results { //查询出来的数组
		fmt.Println(k, v)
	}
}

func main() {

	db, err := sql.Open("mysql", "root:root@tcp(39.104.226.149:3306)/spider?charset=utf8")
	checkErr(err)
	query, err1 := db.Query("select * from novel order by id desc limit 2")
	checkErr(err1)
	printResult(query)

	fmt.Println(strconv.ParseBool(""))
}
