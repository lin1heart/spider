package photo

import (
	"fmt"
	"github.com/lin1heart/spider/go/src/db"
)

func HandlePhotoRows(oss db.OssRow, photos []db.PhotoRow, nextUrl string) {
	fmt.Println("ready HandlePhotoRow photo id", oss, nextUrl, len(photos))
	exist, _ := db.CheckOssNameExist(oss.Name)
	if exist {
		fmt.Println("oss exist", oss)
		return
	}
	id := db.InsertOss(oss)
	for _, photo := range photos {
		handlePhotoRow(id, photo)
	}
	db.UpdateKeyValue(db.MAOMI_KEY, nextUrl)
}

func handlePhotoRow(ossid int, photo db.PhotoRow) {
	exist, _ := db.CheckPhotoExist(photo.CrawlUrl)
	if exist {
		fmt.Println("photo exist", photo)
		return
	}
	photo.OssId = ossid
	db.InsertPhoto(photo)
}
