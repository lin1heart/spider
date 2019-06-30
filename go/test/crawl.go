package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/lin1heart/spider/go/src/util"
	"math/rand"
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

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String(), r.Headers)
	})

	// On every a element which has href attribute call callback
	c.OnHTML(".content_read", func(e *colly.HTMLElement) {

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

	err := c.Visit("https://util.online/headers")
	util.CheckError(err)

}
