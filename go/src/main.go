package main

import (
	"github.com/lin1heart/spider/go/src/novel/www/luoxia/com"
	"github.com/lin1heart/spider/go/src/novel/www/qu/la"
	"github.com/lin1heart/spider/go/src/photo/miaomi"
	"github.com/lin1heart/spider/go/src/util"
)

func imageMain() {
	// 猫咪
	miaomi.Main()
}
func novelMain() {
	// 笔趣阁
	go la.Main()
	// 落霞
	com.Main()
}
func main() {
	if util.ENTRY == "photo" {
		imageMain()
	} else {
		novelMain()
	}

	//imageMain()
}
