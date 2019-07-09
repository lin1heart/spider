package crwal

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"github.com/lin1heart/spider/go/src/util"
	"math/rand"
	"regexp"
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

	rp, err1 := proxy.RoundRobinProxySwitcher(util.RandomProxy())
	if err1 != nil {
		fmt.Println("roundRobin ", err1)
	}
	c.SetProxyFunc(rp)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String(), r.Headers)
	})

	// On every a element which has href attribute call callback
	c.OnHTML(".content_read", func(e *colly.HTMLElement) {
		content := e.ChildText("#content")

		var re = regexp.MustCompile(`\s\schaptererror\(\);`)
		s := re.ReplaceAllString(content, ``)

		match, _ := regexp.MatchString("^正在手打中，客官请稍等片刻，内容更新后，需要重新刷新页面，才能获取最新更新！", content)
		fmt.Println("match", match)

		fmt.Println("content", content, s)

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

	err = c.Visit("https://util.online/headers")
	err = c.Visit("https://util.online/headers")
	util.CheckError(err)

}
