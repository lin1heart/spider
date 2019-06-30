package gaoshoujimo

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/lin1heart/spider/go/src/novel"
	"github.com/lin1heart/spider/go/src/util"
	"math/rand"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Main() {
	for true {
		Crawl()
	}
	//Crawl()
}
func Crawl() {

	var ossId int
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.DisallowedDomains("https://sccdn.002lzj.com"),
		colly.UserAgent(RandomString()),
	)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String(), r.Headers)
	})

	// On every a element which has href attribute call callback
	c.OnHTML(".content_read", func(e *colly.HTMLElement) {
		title := e.ChildText(".bookname h1")
		content := e.ChildText("#content")
		nextRelativeUrl := e.ChildAttr(".bottem2  .next", "href")
		crawlUrl := e.Request.URL.String()

		splitUrl:= strings.Split(e.Request.URL.String(),"/")
		splitUrl[len(splitUrl) - 1] = nextRelativeUrl

		nextAbsoluteUrl := strings.Join(splitUrl, "/")

		if nextRelativeUrl == "./" {
			nextAbsoluteUrl = ""
		}

		novelRow := novel.NovelRow{
			Title:        title,
			Content:      content,
			CrawlUrl:     crawlUrl,
			NextCrawlUrl: nextAbsoluteUrl,
			OssId:        ossId,
		}
		novel.HandleNovelRow(novelRow)

		if nextRelativeUrl == "./" {
			nextAbsoluteUrl = ""
			fmt.Println("will sleep 10 min due to latest")
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

	row := novel.PrepareNovel("高手寂寞", "https://www.qu.la/book/2639/1459927.html")

	fmt.Println("PrepareNovel row ", row)

	ossId = row.Id
	err := c.Visit(row.CrawlUrl)
	util.CheckError(err)

}
