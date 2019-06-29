package novel

import (
	"fmt"
	"github.com/lin1heart/spider/go/src/db"
	"github.com/lin1heart/spider/go/src/util"
	"strconv"
)

type NovelRow struct {
	Title    string
	Content  string
	CrawlUrl string
	OssId    int
}
type OssRow struct {
	Id       int
	Name     string
	CrawlUrl string
	Type     string
	Url      string
}

const NOVEL__TYPE = "NOVEL"

var mysql = db.Mysql

func PrepareNovel(name string, crawlUrl string) OssRow {
	sqlString := fmt.Sprintf("SELECT * FROM oss WHERE name = '%s' ", name)
	rows, err := db.Mysql.Query(sqlString)
	util.CheckError(err)

	ossResults, err := db.ConvertToRows(rows)
	ossId, err := strconv.Atoi(ossResults[0]["id"])

	var oss OssRow
	oss.Id = ossId
	oss.Name = ossResults[0]["name"]

	if len(ossResults) == 0 {
		return insertNewNovel(name, crawlUrl)
	} else {
		novelRows := prepareChapter(ossId)
		if len(novelRows) == 0 {
			oss.CrawlUrl = ossResults[0]["crawl_url"]
		} else {
			oss.CrawlUrl = novelRows[0]["crawl_url"]
		}
		return oss
	}
}

func insertNewNovel(name string, crawlUrl string) OssRow {
	fmt.Println("insertNewNovel", name)
	stmt, err := mysql.Prepare(`INSERT oss (name, type, url, crawl_url) values (?,?,?,?)`)
	util.CheckError(err)

	_, err = stmt.Exec(name, NOVEL__TYPE, nil, crawlUrl)
	util.CheckError(err)
	return PrepareNovel(name, crawlUrl)
}

func prepareChapter(ossId int) map[int]map[string]string {
	sqlString := fmt.Sprintf("SELECT * FROM novel WHERE oss_id = %d ORDER BY chapter_index DESC LIMIT 1", ossId)
	rows, err := db.Mysql.Query(sqlString)
	util.CheckError(err)

	results, err := db.ConvertToRows(rows)
	return results
}

func HandleNovelRow(novel NovelRow) {
	fmt.Println("ready HandleNovelRow", novel)
	exist := checkExist(novel.CrawlUrl)

	if !exist {
		fmt.Println("ignore due to exist", novel)
		return
	}

	novelRows := prepareChapter(novel.OssId)
	var latestChapterIndex int = 0
	var err error
	if len(novelRows[0]) != 0 {
		latestChapterIndex, err = strconv.Atoi(novelRows[0]["chapter_index"])
	}

	stmt, err := mysql.Prepare(`INSERT novel (chapter_index, chapter_title, oss_id, content, crawl_url) values (?,?,?,?,?)`)
	util.CheckError(err)

	res, err := stmt.Exec(latestChapterIndex+1, novel.Title, novel.OssId, novel.Content, novel.CrawlUrl)
	util.CheckError(err)
	id, err := res.LastInsertId()
	util.CheckError(err)
	fmt.Println(id)
}

func checkExist(crawUrl string) bool {
	sqlString := fmt.Sprintf(`SELECT * FROM novel WHERE crawl_url = "%s" `, crawUrl)
	rows, err := db.Mysql.Query(sqlString)
	util.CheckError(err)

	results, err := db.ConvertToRows(rows)
	if len(results) == 0 {
		return true
	}
	return false
}
