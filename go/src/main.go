package main

import (
	"fmt"
	"github.com/lin1heart/spider/go/src/novel/www/luoxia/com"
	"github.com/lin1heart/spider/go/src/novel/www/qu/la"
	"github.com/lin1heart/spider/go/src/photo/miaomi"
	"github.com/lin1heart/spider/go/src/util"
)

func main() {

	fmt.Println("ENV", util.ENV)
	fmt.Println("ENTRY", util.ENTRY)

	if util.ENTRY != "photo" {
		go la.Main()
		com.Main()
	} else {
		miaomi.Main()
	}

}
