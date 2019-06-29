package novel

import (
	"fmt"
	"github.com/lin1heart/spider/go/src/db"
	"github.com/lin1heart/spider/go/src/util"
)

type Novel struct {
	Title   string
	Content string
	Url		string
}
type Oss struct {
	Name string
	Type string
	url	string
}

const NOVEL__TYPE = "NOVEL"
var mysql = db.Mysql


func PrepareNovel (name string, crawlUrl string) map[string]string {
	sqlString := fmt.Sprintf("SELECT* FROM oss WHERE name = '%s' ", name)
	rows, err := db.Mysql.Query(sqlString)
	util.CheckError(err)

	results, err := db.ConvertToRows(rows)

	if len(results)==0 {
		return insertNewNovel(name, crawlUrl)
	} else {
		return results[0]
	}
}

func insertNewNovel (name string, crawlUrl string) map[string]string {
	fmt.Println("insertNewNovel", name)
	stmt, err := mysql.Prepare(`INSERT oss (name, type, url, crawl_url) values (?,?,?,?)`)
	util.CheckError(err)

	_, err = stmt.Exec(name, NOVEL__TYPE, nil, crawlUrl)
	util.CheckError(err)
	return PrepareNovel(name, crawlUrl)
}

func HandleNovelRow(novel Novel) {
	fmt.Println("ready HandleNovelRow", novel)
	stmt, err := mysql.Prepare(`INSERT novel (chapter_index,chapter_title, oss_id, content, url) values (?,?,?)`)
	//stmt, err := mysql.Prepare(fmt.Sprintf(""))
	util.CheckError(err)

	res, err := stmt.Exec("tony", 20, 1)
	util.CheckError(err)
	id, err := res.LastInsertId()
	util.CheckError(err)
	fmt.Println(id)
}

