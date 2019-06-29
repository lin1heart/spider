package dushixianzun

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/lin1heart/spider/go/src/novel"
)


func Crawl() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.DisallowedDomains("https://sccdn.002lzj.com"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML(".content_read", func(e *colly.HTMLElement) {
		title := e.ChildText(".bookname h1")
		content := e.ChildText("#content")
		url := e.Request.URL.String()

		novelRow := novel.Novel{
			Title:   title,
			Content: content,
			Url:	 url,
		}
		//fmt.Println(novel)
		//HandleNovelRow(novel)
		novel.HandleNovelRow(novelRow)

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


	novel.PrepareNovel("都市仙尊", "https://www.qu.la/book/85467/4563618.html")
	//err := c.Visit("https://www.qu.la/book/85467/4563618.html")
	//util.CheckError(err)


}
