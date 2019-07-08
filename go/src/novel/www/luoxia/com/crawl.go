package com

import "time"

func Main() {
	go LoopQueryNullCrawlUrlOss()
	for true {
		//queryTodoOss()
		time.Sleep(10 * time.Second)
	}
}