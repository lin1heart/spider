package miaomi

import (
	"crypto/tls"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"github.com/lin1heart/spider/go/src/db"
	"github.com/lin1heart/spider/go/src/photo"
	"github.com/lin1heart/spider/go/src/util"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Crawl(crawlUrl string) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			log.Printf("Crawl defer e", err) // 这里的err其实就是panic传入的内容，55
		}
	}()

	c := colly.NewCollector(
		colly.DisallowedDomains("www.google-analytics.com", "tpc.googlesyndication.com", "www.ftpd188.com", "img.alicdn.com", "sc02.alicdn.com"),
		colly.UserAgent(util.RandomString()),
	)

	randomProxy := util.RandomProxy()
	fmt.Println("random proxy ", randomProxy)

	rp, err1 := proxy.RoundRobinProxySwitcher(randomProxy)
	if err1 != nil {
		fmt.Println("roundRobin ", err1)
	}
	c.SetProxyFunc(colly.ProxyFunc(rp))

	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String(), r.Headers)
	})

	c.OnHTML("main", func(e *colly.HTMLElement) {
		mediumType := e.ChildText(".cat_pos_l a:nth-child(2)")
		//imageType := e.ChildText(".cat_pos_l a:nth-child(3)")
		imageTitle := e.ChildText(".cat_pos_l a:nth-child(4)")

		webUrl := e.Request.URL.String()

		nextRelativeUrl := e.ChildAttr(".next-page .content-next2 a", "href")
		nextSplits := strings.Split(webUrl, "tupian")
		nextAbsoluteUrl := nextSplits[0] + nextRelativeUrl

		images := []string{}

		photos := []db.PhotoRow{}

		e.ForEach(".content img", func(_ int, elem *colly.HTMLElement) {
			photoUrl := elem.Attr("data-original")
			photoTitle := elem.Attr("title")

			splits := strings.Split(photoUrl, "/")

			if len(splits) == 7 {
				splits2 := strings.Split(splits[6], "-")

				splits3 := strings.Split(splits2[1], ".")
				index, err := strconv.Atoi(splits3[0])
				util.CheckError(err)

				images = append(images, elem.Attr("data-original"))

				photo := db.PhotoRow{
					Title:    photoTitle,
					Url:      "",
					CrawlUrl: photoUrl,
					OssId:    -1,
					Index:    index,
				}
				photos = append(photos, photo)
			}
		})

		//fmt.Println("mediumType",mediumType)
		//fmt.Println("imageType",imageType)
		//fmt.Println("webUrl",webUrl)
		//fmt.Println("crawlUrl",crawlUrl)
		//fmt.Println("images",images)

		if mediumType == "福利图片" {

			oss := db.OssRow{
				Id:       -1,
				Name:     imageTitle,
				Url:      "",
				Type:     "PHOTO_PURE",
				CrawlUrl: webUrl,
			}
			photo.HandlePhotoRows(oss, photos, nextAbsoluteUrl)
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

	err := c.Visit(crawlUrl)
	if err != nil {
		log.Printf("Visit %s e ", crawlUrl, err)
	}

}

func PreparePhoto(key string, value string) string {
	dbValue := db.QueryKeyValue(key)
	if dbValue == "" {
		db.InsertKeyValue(key, value)
		return value
	}
	return dbValue
}

func Main() {
	crawlUrl := PreparePhoto(db.MAOMI_KEY, db.MAOMI)
	Crawl(crawlUrl)
}
