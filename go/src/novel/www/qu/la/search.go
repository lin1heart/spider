package la

import (
	"crypto/tls"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/lin1heart/spider/go/src/db"
	"github.com/lin1heart/spider/go/src/util"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var SEARCH_URL = "https://sou.xanbhx.com/search?siteid=qula&q="

var mysql = db.Mysql
var searchedIds []string

func crawlNovelChapters(id int, basehref string) {

	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			log.Printf("crawlNovelChapters defer e", err) // 这里的err其实就是panic传入的内容，55
		}
	}()


	var counter = 0
	var minHtmlIndex = 99999999999
	var htmlhref = ""

	c := colly.NewCollector(
		colly.DisallowedDomains("https://sccdn.002lzj.com"),
		colly.UserAgent(util.RandomString()),
	)
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String(), r.Headers)
	})

	c.OnHTML("#list dd a", func(e *colly.HTMLElement) {
		counter += 1

		href := e.Attr("href")
		splitsArr := strings.Split(href, "/")

		//text := strings.TrimSpace(e.Text)
		//fmt.Println("list text ", text)
		//fmt.Println("list href ", href)

		if len(splitsArr) == 4 {
			htmlIndexStr := splitsArr[3]

			splits2 := strings.Split(htmlIndexStr, ".html")
			htmlIndex, err := strconv.Atoi(splits2[0])
			if err == nil {
				if htmlIndex < minHtmlIndex {
					minHtmlIndex = htmlIndex
					htmlhref = href
				}
			} else {
				fmt.Print(err)
			}
		}
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL, r.StatusCode)
	})

	c.OnScraped(func(r *colly.Response) {

		updateOssCrawlUrl(basehref+htmlhref, id)

	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})
	err := c.Visit(basehref)
	util.CheckError(err)

}

func crawlSearchResult(name string, id int) {

	c := colly.NewCollector(
		colly.DisallowedDomains("https://sccdn.002lzj.com"),
		colly.UserAgent(util.RandomString()),
	)
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String(), r.Headers)
	})

	c.OnHTML("#search-main .search-list li .s2 a", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		href := e.Attr("href")

		if name == text {
			fmt.Println("search text ", text)
			fmt.Println("search href ", href)
			crawlNovelChapters(id, href)
		}
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	escapedTitle := url.QueryEscape(name)
	err := c.Visit(SEARCH_URL + escapedTitle)
	if err != nil {
		log.Printf("Visit %s e ", SEARCH_URL + escapedTitle, err)
	}

}
func updateOssCrawlUrl(crawlUrl string, novelId int) {
	stmt, err := mysql.Prepare(`UPDATE oss SET crawl_url=? WHERE id=?`)
	util.CheckError(err)
	res, err := stmt.Exec(crawlUrl, novelId)
	util.CheckError(err)
	num, err := res.RowsAffected()
	fmt.Printf("updateNextCrawlUrl id %d affect row %d", num)
	util.CheckError(err)
}

func queryNullCrawlUrlOss() {
	sqlString := fmt.Sprintf("SELECT * FROM oss WHERE (crawl_url IS NULL or crawl_url = '') AND type ='NOVEL' %s ", generateSqlSuffix(searchedIds))
	rows, err := mysql.Query(sqlString)

	ossResults, err := db.ConvertToRows(rows)
	util.CheckError(err)

	for _, ossRow := range ossResults {
		id, err := strconv.Atoi(ossRow["id"])
		util.CheckError(err)
		name := ossRow["name"]
		crawlSearchResult(name, id)
		searchedIds = append(searchedIds, ossRow["id"])
		fmt.Printf("add ossId %v to searchedIds \n", id)
	}
}
func LoopQueryNullCrawlUrlOss() {
	for true {

		queryNullCrawlUrlOss()
		time.Sleep(10 * time.Second)

	}
}
