package com

import (
	"crypto/tls"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/lin1heart/spider/go/src/db"
	"github.com/lin1heart/spider/go/src/novel"
	"github.com/lin1heart/spider/go/src/util"
	"log"
	"net/http"
	"strings"
	"time"
)

const WWW_QU_LA_PREFIX = "https://www.luoxia.com/%"

var crawlingIds []string

func Crawl(name string, crawlUrl string) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			log.Printf("Crawl defer e", err) // 这里的err其实就是panic传入的内容，55
		}
	}()

	var ossId int
	c := colly.NewCollector(
		colly.DisallowedDomains("www.google-analytics.com", "tpc.googlesyndication.com", "dt.adsafeprotected.com"),
		colly.UserAgent(util.RandomString()),
	)
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String(), r.Headers)
	})

	c.OnHTML(".post", func(e *colly.HTMLElement) {
		title := e.ChildText("#nr_title")
		content := e.ChildText("#nr1 p ")
		newContent := util.FilterEmoji(content)
		fmt.Println("content", content)
		fmt.Println("newContent", newContent)
		nextRelativeUrl := e.ChildAttr(".nav2 .next a", "href")
		nextSpiltUrls := strings.Split(nextRelativeUrl, "/")

		nextSpiltUrl := ""
		if len(nextSpiltUrls) == 5 {

			nextSpiltUrl = nextSpiltUrls[4]
		}
		crawlUrl := e.Request.URL.String()

		splitUrl := strings.Split(e.Request.URL.String(), "/")
		splitUrl[len(splitUrl)-1] = nextSpiltUrl

		nextAbsoluteUrl := strings.Join(splitUrl, "/")

		if nextRelativeUrl == "" {
			nextAbsoluteUrl = ""
		}

		novelRow := novel.NovelRow{
			Title:        title,
			Content:      newContent,
			CrawlUrl:     crawlUrl,
			NextCrawlUrl: nextAbsoluteUrl,
			OssId:        ossId,
		}
		novel.HandleNovelRow(novelRow)

		if nextAbsoluteUrl == "" {
			fmt.Printf("%s will sleep 10 min due to latest \n", name)
			time.Sleep(10 * time.Minute)
		} else {
			time.Sleep(5 * time.Second)
		}
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL, r.StatusCode)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	row := novel.PrepareNovel(name, crawlUrl)

	fmt.Println("PrepareNovel row ", row)

	ossId = row.Id
	err := c.Visit(row.CrawlUrl)
	if err != nil {
		log.Printf("Visit %s e ", row.CrawlUrl, err)
	}

}

func loopCrawl(name string, crawlUrl string) {
	for true {
		Crawl(name, crawlUrl)
	}
}

func queryTodoOss() {

	sqlString := fmt.Sprintf("SELECT * FROM oss WHERE crawl_url LIKE '%s' AND complete = 0 %s ", WWW_QU_LA_PREFIX, util.GenerateSqlIdsSuffix(crawlingIds))
	rows, err := db.Mysql.Query(sqlString)

	ossResults, err := db.ConvertToRows(rows)
	util.CheckError(err)

	for _, ossRow := range ossResults {
		id := ossRow["id"]
		name := ossRow["name"]
		crawlUrl := ossRow["crawl_url"]

		fmt.Printf("add ossId %s to CrawlingIds \n", id)

		crawlingIds = append(crawlingIds, id)
		go loopCrawl(name, crawlUrl)
	}
}

func Main() {
	go LoopQueryNullCrawlUrlOss()
	for true {
		queryTodoOss()
		time.Sleep(10 * time.Second)
	}
}
