package db

import (
	"fmt"
	"github.com/lin1heart/spider/go/src/util"
)

func UpdateOssCrawlUrl(crawlUrl string, novelId int) {
	stmt, err := Mysql.Prepare(`UPDATE oss SET crawl_url=? WHERE id=?`)
	util.CheckError(err)
	res, err := stmt.Exec(crawlUrl, novelId)
	util.CheckError(err)
	num, err := res.RowsAffected()
	fmt.Printf("updateNextCrawlUrl id %d affect row %d", num)
	util.CheckError(err)
}
