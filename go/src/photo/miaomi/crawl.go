package miaomi

import (
	"crypto/tls"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/lin1heart/spider/go/src/db"
	"github.com/lin1heart/spider/go/src/photo"
	"github.com/lin1heart/spider/go/src/util"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	//randomProxy := util.RandomProxy()
	//fmt.Println("random proxy ", randomProxy)
	//
	//rp, err1 := proxy.RoundRobinProxySwitcher(randomProxy)
	//if err1 != nil {
	//	fmt.Println("roundRobin ", err1)
	//}
	//c.SetProxyFunc(colly.ProxyFunc(rp))

	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String(), r.Headers)
	})

	c.OnHTML("main", func(e *colly.HTMLElement) {
		mediumType := e.ChildText(".cat_pos_l a:nth-child(2)")
		imageType := e.ChildText(".cat_pos_l a:nth-child(3)")
		imageTitle := e.ChildText(".cat_pos_l a:nth-child(4)")

		webUrl := e.Request.URL.String()

		nextRelativeUrl := e.ChildAttr(".next-page .content-next2 a", "href")
		nextSplits := strings.Split(webUrl, "/tupian")
		nextAbsoluteUrl := nextSplits[0] + nextRelativeUrl

		images := []string{}

		photos := []db.PhotoRow{}

		photoType := "PHOTO_ERROR"
		switch imageType {
		case "清纯唯美":
			photoType = "PHOTO_PURE"
			break
		case "自拍偷拍":
			photoType = "PHOTO_SELF"
			break
		case "亚洲色图":
			photoType = "PHOTO_EAST"
			break
		case "欧美色图":
			photoType = "PHOTO_WEST"
			break
		case "美腿丝袜":
			photoType = "PHOTO_UNIFORM"
			break
		case "乱伦熟女":
			photoType = "PHOTO_RAPE"
			break
		case "卡通动漫":
			photoType = "PHOTO_COMIC"
			break
		}

		e.ForEach(".content img", func(_ int, elem *colly.HTMLElement) {
			photoUrl := elem.Attr("data-original")
			photoTitle := elem.Attr("title")

			splits := strings.Split(photoUrl, "/")

			if len(splits) == 7 {
				//splits2 := strings.Split(splits[6], "-")
				var indexString string
				splits2 := strings.Split(splits[6], "_tmb.jpg")
				if len(splits2) == 2 {
					indexString = splits2[0]
				} else {
					splits3 := strings.Split(splits[6], ".")
					indexString = splits3[0]
				}

				index, err := strconv.Atoi(indexString)
				util.CheckError(err)

				images = append(images, elem.Attr("data-original"))

				photo := db.PhotoRow{
					Title:    photoTitle,
					Url:      "",
					CrawlUrl: strings.TrimSpace(photoUrl),
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

		if nextRelativeUrl == "" {
			fmt.Println("miaomi return due to nextRelativeUrl ", nextRelativeUrl)
			time.Sleep(1 * time.Hour)
			return
		}
		if mediumType == "福利图片" {

			oss := db.OssRow{
				Id:       -1,
				Name:     imageTitle,
				Url:      "",
				Type:     photoType,
				CrawlUrl: webUrl,
			}
			photo.HandlePhotoRows(oss, photos, nextAbsoluteUrl)
			//time.Sleep(30 * time.Second)

		}

	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL, r.StatusCode)
		splits := strings.Split(fmt.Sprintf("%s", r.Request.URL), "/")
		if len(splits) == 3 {
			fmt.Println("new site", r.Request.URL, r.StatusCode)

			originSplits := strings.Split(crawlUrl, "/")

			newUrl := fmt.Sprintf("%s/%s/%s", r.Request.URL, originSplits[3], originSplits[4])
			fmt.Println("new Url", newUrl)
			db.InsertKeyValue(db.MAOMI_KEY, newUrl)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	err := c.Visit(crawlUrl)
	if err != nil {
		log.Printf("Visit %s e ", crawlUrl, err)
		time.Sleep(5 * time.Second)
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

func loopCrawl() {
	for true {
		crawlUrl := PreparePhoto(db.MAOMI_KEY, db.MAOMI)
		Crawl(crawlUrl)
	}
}

func Main() {
	loopCrawl()
	//go loopCrawl()
	//photo.Sync()
}
