package la

import (
	"crypto/tls"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/lin1heart/spider/go/src/db"
	"github.com/lin1heart/spider/go/src/novel"
	"github.com/lin1heart/spider/go/src/util"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const WWW_QU_LA_PREFIX = "https://www.qu.la/%"

var crawlingIds []string

func Crawl(name string, crawlUrl string) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			log.Printf("Crawl defer e", err) // 这里的err其实就是panic传入的内容，55
		}
	}()

	var ossId int
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

		var re = regexp.MustCompile(`\s\schaptererror\(\);`)
		cleanContent := re.ReplaceAllString(content, ``)

		match, _ := regexp.MatchString("^正在手打中，客官请稍等片刻，内容更新后，需要重新刷新页面，才能获取最新更新", content)

		if match == true && nextAbsoluteUrl == "" {
			fmt.Printf("%s will sleep 10 min due to invalid content \n", name, content)
			time.Sleep(10 * time.Minute)
			return
		}

		novelRow := novel.NovelRow{
			Title:        title,
			Content:      cleanContent,
			CrawlUrl:     crawlUrl,
			NextCrawlUrl: nextAbsoluteUrl,
			OssId:        ossId,
		}
		novel.HandleNovelRow(novelRow)

		//fmt.Printf("%s nextRelativeUrl %s nextAbsoluteUrl %s \n", name, nextRelativeUrl, nextAbsoluteUrl)

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
