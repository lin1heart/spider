package main

import (
	"github.com/lin1heart/spider/go/src/novel/dushixianzun"
	"github.com/lin1heart/spider/go/src/util"
)

func main() {
	dushixianzun.Crawl()
	var e error
	util.CheckError(e)
}

