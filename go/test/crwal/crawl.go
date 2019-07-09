package crwal

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"github.com/lin1heart/spider/go/src/util"
	"math/rand"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.DisallowedDomains("https://sccdn.002lzj.com"),
		colly.UserAgent(RandomString()),
	)

	rp, err1 := proxy.RoundRobinProxySwitcher(util.ProxyList...)
	if err1 != nil {
		fmt.Println("roundRobin ", err1)
	}
	c.SetProxyFunc(rp)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String(), r.Headers)
	})

	// On every a element which has href attribute call callback
	c.OnHTML(".post", func(e *colly.HTMLElement) {
		title := e.ChildText("#nr_title")
		content := e.ChildText("#nr1 ")
		nextAbsoluteUrl := e.ChildAttr(".nav2 .next a", "href")

		crawlUrl := e.Request.URL.String()

		fmt.Println("title", title)
		fmt.Println("content", content)
		if nextAbsoluteUrl == "" {
			fmt.Println("nextAbsoluteUrl is null")
			nextAbsoluteUrl = ""
		}
		nextUrlSplits := strings.Split(nextAbsoluteUrl, "/")
		currentUrlSplits := strings.Split(crawlUrl, "/")
		if len(nextUrlSplits) == 5 && nextUrlSplits[4] != currentUrlSplits[4] {
			fmt.Println("invalid next url")
			nextAbsoluteUrl = ""
		}

	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL, r.StatusCode, string(r.Body))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
	c.AllowURLRevisit = true
	var err error

	err = c.Visit("https://www.luoxia.com/baiye/12453.htm")
	util.CheckError(err)

}
