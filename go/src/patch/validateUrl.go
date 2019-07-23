package patch

import (
	"fmt"
	"github.com/lin1heart/spider/go/src/db"
)

func init() {
	fmt.Println("patch init")
}
func main() {
	startId := 0
	//endId := 233648
	endId := 150

	for index := startId; index < endId; {
		fmt.Println("index ", index)
		index += 50
		results := db.QueryPhotoAfterRow(index)

		fmt.Println("results", results)
	}
}
