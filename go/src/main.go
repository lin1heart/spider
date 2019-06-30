package main

import (
	"github.com/lin1heart/spider/go/src/novel/www.qu.la/dushixianzun"
	"github.com/lin1heart/spider/go/src/novel/www.qu.la/gaoshoujimo"
)

func main() {
	go dushixianzun.Main()
	gaoshoujimo.Main()
}
