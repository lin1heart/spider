package la

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/lin1heart/spider/go/src/novel"
	"github.com/lin1heart/spider/go/src/util"
	"strings"
	"time"
)

func Main() {
	//go gaoshoujimo() // 已完结
	//go gaoshoujimo2() // 已完结
	//go chongshengzhiwonengshengji() // 已完结
	go nitianxieshen()
	dushixianzun()
}
func gaoshoujimo() {
	for true {
		Crawl("高手寂寞", "https://www.qu.la/book/2639/1459927.html")
	}
}
func gaoshoujimo2() {
	for true {
		Crawl("高手寂寞2", "https://www.qu.la/book/620/410112.html")
	}
}
func dushixianzun() {
	for true {
		Crawl("都市仙尊", "https://www.qu.la/book/85467/4563618.html")
	}
}
func nitianxieshen() {
	for true {
		Crawl("逆天邪神", "https://www.qu.la/book/3648/10567026.html")
	}
}
func chongshengzhiwonengshengji() {
	for true {
		Crawl("重生之我能升级", "https://www.qu.la/book/182726/9302392.html")
	}
}



func Crawl(name string, crawlUrl string) {

	var ossId int
	c := colly.NewCollector(
		colly.DisallowedDomains("https://sccdn.002lzj.com"),
		colly.UserAgent(util.RandomString()),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String(), r.Headers)
	})

	c.OnHTML(".content_read", func(e *colly.HTMLElement) {
		title := e.ChildText(".bookname h1")
		content := e.ChildText("#content")
		nextRelativeUrl := e.ChildAttr(".bottem2  .next", "href")
		crawlUrl := e.Request.URL.String()

		splitUrl := strings.Split(e.Request.URL.String(), "/")
		splitUrl[len(splitUrl)-1] = nextRelativeUrl

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

		fmt.Printf("%s nextRelativeUrl %s nextAbsoluteUrl %s \n", name, nextRelativeUrl, nextAbsoluteUrl)


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
	util.CheckError(err)

}
