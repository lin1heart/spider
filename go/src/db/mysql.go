package db

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/lin1heart/spider/go/src/ssh"
	"github.com/lin1heart/spider/go/src/util"
	"net"
)
import _ "github.com/go-sql-driver/mysql"

var Mysql *sql.DB

func init() {
	fmt.Println("mysql init", util.ENV)

	if util.ENV == "prod" {
		var err error
		Mysql, err = sql.Open("mysql", "root:root@tcp(39.104.226.149:3306)/spider?charset=utf8")
		util.CheckError(err)

		Mysql.SetMaxOpenConns(20)
		Mysql.SetMaxIdleConns(10)
		Mysql.Ping()
		return
	}
	dbUser := "root"           // DB username
	dbPass := "root"           // DB Password
	dbHost := "127.0.0.1:3306" // DB Hostname/IP
	dbName := "spider"         // Database name

	mysql.RegisterDial("tcpchannel", func(addr string) (net.Conn, error) {
		return ssh.SshClient.Dial("tcp", addr)
	})

	var err error
	Mysql, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcpchannel(%s)/%s", dbUser, dbPass, dbHost, dbName))
	util.CheckError(err)
	fmt.Printf("Successfully connected to the db\n")
	//Mysql.SetMaxOpenConns(20)
	Mysql.SetMaxIdleConns(10)
	Mysql.Ping()

}

func PrintResult(query *sql.Rows) {
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

func ConvertToRows(query *sql.Rows) (map[int]map[string]string, error) {
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
		}
		row := make(map[string]string) //每行数据
		for k, v := range values {     //每行数据是放在values里面，现在把它挪到row里
			key := column[k]
			row[key] = string(v)
		}
		results[i] = row //装入结果集中
		i++
	}
	//for k, v := range results { //查询出来的数组
	//	fmt.Println(k, v)
	//}
	return results, nil
}
