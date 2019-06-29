package dushixianzun

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/lin1heart/spider/go/src/novel"
	"github.com/lin1heart/spider/go/src/util"
	"time"
)

func Crawl() {

	var ossId int
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.DisallowedDomains("https://sccdn.002lzj.com"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML(".content_read", func(e *colly.HTMLElement) {
		title := e.ChildText(".bookname h1")
		content := e.ChildText("#content")
		nextRelativeUrl := e.ChildAttr(".bottem2  .next", "href")
		crawlUrl := e.Request.URL.String()
		nextAbsoluteUrl := "https://www.qu.la/book/85467/" + nextRelativeUrl

		novelRow := novel.NovelRow{
			Title:    title,
			Content:  content,
			CrawlUrl: crawlUrl,
			OssId:    ossId,
		}
		novel.HandleNovelRow(novelRow)

		fmt.Println("nextAbsoluteUrl", nextAbsoluteUrl)
		if nextRelativeUrl == "./" {
			fmt.Println("will sleep 10 min due to latest")
			time.Sleep(10 * time.Minute)
		}
		time.Sleep(30 * time.Second)

		c.Visit(nextAbsoluteUrl)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
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

	row := novel.PrepareNovel("都市仙尊", "https://www.qu.la/book/85467/4563618.html")

	fmt.Println("PrepareNovel CrawlUrl", row.CrawlUrl)
	fmt.Println("PrepareNovel Name", row.Name)
	fmt.Println("PrepareNovel Id", row.Id)
	fmt.Println("PrepareNovel row ", row)

	ossId = row.Id
	err := c.Visit(row.CrawlUrl)
	util.CheckError(err)

}
