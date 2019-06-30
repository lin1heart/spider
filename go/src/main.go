package main

import (
	"github.com/lin1heart/spider/go/src/novel/www.qu.la/dushixianzun"
	"github.com/lin1heart/spider/go/src/novel/www.qu.la/gaoshoujimo"
	"github.com/lin1heart/spider/go/src/novel/www.qu.la/gaoshoujimo2"
)

func main() {
	go dushixianzun.Main()
	go gaoshoujimo.Main()
	gaoshoujimo2.Main()
}
