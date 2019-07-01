package la

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/lin1heart/spider/go/src/db"
	"github.com/lin1heart/spider/go/src/novel"
	"github.com/lin1heart/spider/go/src/util"
	"strings"
	"time"
)

const WWW_QU_LA_PREFIX = "https://www.qu.la/%"

var crawlingIds []string

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

func loopCrawl(name string, crawlUrl string) {
	for true {
		Crawl(name, crawlUrl)
	}
}

func generateSqlSuffix(crawlingIds []string) string {
	crawlingIdsStr := strings.Join(crawlingIds, ",")

	if crawlingIdsStr == "" {
		return ""
	}
	sqlSuffix := fmt.Sprintf(" AND id NOT IN ( %s )", crawlingIdsStr)
	return sqlSuffix
}

func queryTodoOss() {

	sqlString := fmt.Sprintf("SELECT * FROM oss WHERE crawl_url LIKE '%s' AND complete = 0 %s ", WWW_QU_LA_PREFIX, generateSqlSuffix(crawlingIds))
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
	for true {
		queryTodoOss()
		time.Sleep(5 * time.Second)
	}
}
