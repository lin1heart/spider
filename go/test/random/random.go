package random

import (
	"fmt"
	"github.com/lin1heart/spider/go/src/db"
	"math/rand"
)

func main() {

	var list = []string{"11", "22", "33", "44"}

	fmt.Println(rand.Intn(len(list)))
	fmt.Println(rand.Intn(len(list)))

	imageTypes := []interface{}{"PHOTO_PURE", "PHOTO_COMIC"}
	rows := db.QueryEmptyPhotos(imageTypes)
	fmt.Println("rows", rows)
}
