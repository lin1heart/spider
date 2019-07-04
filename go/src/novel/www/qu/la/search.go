package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/lin1heart/spider/go/src/util"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var SEARCH_URL ="https://sou.xanbhx.com/search?siteid=qula&q="

func crawlSearchResult(name string) {

	//var ossId int
	var counter = 0
	var minHtmlIndex = 99999999999

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

		text := strings.TrimSpace(e.Text)
		href := e.Attr("href")
		splitsArr := strings.Split(href,".html")
		htmlIndexStr := splitsArr[0]
		htmlIndex, err := strconv.Atoi(htmlIndexStr)
		util.CheckError(err)

		counter += 1
		if htmlIndex < minHtmlIndex {
			minHtmlIndex = htmlIndex
		}
		if counter == 20 {
			fmt.Println("list text ",text)
			fmt.Println("list href ",href)
		}
		fmt.Println("counter", counter)
	})


	c.OnHTML("#search-main .search-list li .s2 a", func(e *colly.HTMLElement) {
		text := strings.TrimSpace( e.Text)
		href := e.Attr("href")

		if name == text {
			fmt.Println("search text ", text)
			fmt.Println("search href ", href)
			c.Visit(href)
		}


		//novelRow := novel.NovelRow{
		//	Title:        title,
		//	Content:      cleanContent,
		//	CrawlUrl:     crawlUrl,
		//	NextCrawlUrl: nextAbsoluteUrl,
		//	OssId:        ossId,
		//}
		//novel.HandleNovelRow(novelRow)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL, r.StatusCode)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	escapedTitle := url.QueryEscape(name)
	fmt.Println("escapedTitle", escapedTitle)
	fmt.Println("2", SEARCH_URL + escapedTitle)
	err := c.Visit(SEARCH_URL + escapedTitle)
	util.CheckError(err)

}

func main() {
	crawlSearchResult("遮天")
}


