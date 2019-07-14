package main

import (
	"github.com/lin1heart/spider/go/src/novel/www/luoxia/com"
	"github.com/lin1heart/spider/go/src/novel/www/qu/la"
)

func main() {
	go la.Main()
	com.Main()
	//miaomi.Main()
}
