package novel

import (
	"fmt"
)

type Novel struct {
	Title   string
	Content string
	Url		string
}

func HandleNovelRow(novel Novel) {
	fmt.Println("ready HandleNovelRow", novel)
}

