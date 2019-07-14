package main

import (
	"github.com/lin1heart/spider/go/src/novel/www/luoxia/com"
	"github.com/lin1heart/spider/go/src/novel/www/qu/la"
	"github.com/lin1heart/spider/go/src/photo/miaomi"
	"github.com/lin1heart/spider/go/src/util"
)

func main() {

	if util.ENTRY != "photo" {
		go la.Main()
		com.Main()
	} else {
		miaomi.Main()
	}

}
