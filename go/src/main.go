package main

import (
	"github.com/lin1heart/spider/go/src/novel/www/luoxia/com"
	"github.com/lin1heart/spider/go/src/novel/www/qu/la"
	"github.com/lin1heart/spider/go/src/photo/miaomi"
)

func main() {
	go la.Main()
	go com.Main()
	miaomi.Main()
}
