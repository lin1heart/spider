package db

import (
	"fmt"
	"github.com/lin1heart/spider/go/src/util"
	"strconv"
	"strings"
)

const MAOMI_KEY = "maomi.photo.pure"
const MAOMI = "https://www.968uy.com/tupian/12457.html"

func UpdateOssCrawlUrl(crawlUrl string, novelId int) {
	stmt, err := Mysql.Prepare(`UPDATE oss SET crawl_url=? WHERE id=?`)
	util.CheckError(err)
	res, err := stmt.Exec(crawlUrl, novelId)
	util.CheckError(err)
	num, err := res.RowsAffected()
	fmt.Printf("updateNextCrawlUrl id %d affect row %d \n", num)
	util.CheckError(err)
}

func QueryKeyValue(key string) string {
	sqlString := fmt.Sprintf("SELECT * FROM key_value WHERE `key` = '%s' ORDER BY id DESC LIMIT 1", key)
	rows, err := Mysql.Query(sqlString)
	util.CheckError(err)

	results, err := ConvertToRows(rows)
	if len(results) != 0 {
		return results[0]["value"]
	}
	return ""
}

func InsertKeyValue(key string, value string) {
	stmt, err := Mysql.Prepare("INSERT key_value (`key`, value) values (?,?)")
	util.CheckError(err)

	_, err = stmt.Exec(key, value)
	util.CheckError(err)
}
func UpdateKeyValue(key string, value string) {
	stmt, err := Mysql.Prepare("UPDATE key_value SET value = ? WHERE `key`= ? ")
	util.CheckError(err)
	res, err := stmt.Exec(value, key)
	util.CheckError(err)
	num, err := res.RowsAffected()
	fmt.Printf("updateKeyValue key %s value %s, affect rows %d \n", key, value, num)
	util.CheckError(err)
}

func CheckOssNameExist(name string) (bool, int) {
	sqlString := fmt.Sprintf(`SELECT * FROM oss WHERE name = "%s" `, name)
	rows, err := Mysql.Query(sqlString)
	util.CheckError(err)

	results, err := ConvertToRows(rows)
	if len(results) != 0 {
		id, err := strconv.Atoi(results[0]["id"])
		util.CheckError(err)
		return true, id
	}
	return false, -1
}

func CheckPhotoExist(crawl_url string) (bool, int) {
	sqlString := fmt.Sprintf(`SELECT * FROM photo WHERE crawl_url = "%s" `, crawl_url)
	rows, err := Mysql.Query(sqlString)
	util.CheckError(err)

	results, err := ConvertToRows(rows)
	if len(results) != 0 {
		id, err := strconv.Atoi(results[0]["id"])
		util.CheckError(err)
		return true, id
	}
	return false, -1
}

func InsertOss(row OssRow) int {
	stmt, err := Mysql.Prepare("INSERT oss (name, type, crawl_url) values (?,?,?)")
	util.CheckError(err)

	res, err := stmt.Exec(row.Name, row.Type, row.CrawlUrl)
	util.CheckError(err)
	id, err := res.LastInsertId()

	return int(id)
}

func InsertPhoto(row PhotoRow) int {
	stmt, err := Mysql.Prepare("INSERT photo (title, url, crawl_url, oss_id, `index`) values (?,?,?,?,?)")
	util.CheckError(err)

	res, err := stmt.Exec(row.Title, row.Url, row.CrawlUrl, row.OssId, row.Index)
	util.CheckError(err)
	id, err := res.LastInsertId()

	return int(id)
}

func QueryEmptyPhotos(imageTypes []interface{}) map[int]map[string]string {

	sqlString := fmt.Sprintf(`select photo.id, photo.crawl_url, oss.id as oss_id, photo.title, photo.index 
		from oss 
		inner join photo on photo.oss_id = oss.id 
		and (photo.url = "" or photo.url is null)
		and oss.type in ( ? %s )
		limit 50
	`, strings.Repeat(",?", len(imageTypes)-1))
	fmt.Println("sqlString", sqlString)
	args := []interface{}{}
	args = append(args, imageTypes...)

	rows, err := Mysql.Query(sqlString, args...)
	util.CheckError(err)

	results, err := ConvertToRows(rows)
	return results
}
func UpdatePhotoUrl(id int, url string) {
	stmt, err := Mysql.Prepare(`UPDATE photo SET url=? WHERE id=?`)
	util.CheckError(err)
	res, err := stmt.Exec(url, id)
	util.CheckError(err)
	num, err := res.RowsAffected()
	fmt.Printf("UpdatePhotoUrl id %d affect row %d \n", id, num)
	util.CheckError(err)
}
