package main

import (
	"github.com/amanoooo/spider/go/src/novel/dushixianzun"
	"github.com/amanoooo/spider/go/src/util"
)

func main() {
	dushixianzun.Crawl()
	var e error
	util.CheckError(e)
}

