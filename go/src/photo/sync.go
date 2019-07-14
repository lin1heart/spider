package photo

import (
	"bytes"
	"fmt"
	"github.com/lin1heart/spider/go/src/db"
	"github.com/lin1heart/spider/go/src/util"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

func Sync() {
	for true {
		task()
		time.Sleep(5 * time.Second)
	}

}

var UPLOAD_BASE = "http://localhost:8888"

func task() {

	rows := db.QueryEmptyPhotos("PHOTO_PURE")

	for _, row := range rows {

		url := row["crawl_url"]
		title := row["title"]
		index := row["index"]
		ossId := row["oss_id"]
		idString := row["id"]

		id, err := strconv.Atoi(idString)
		util.CheckError(err)

		fileName := title + "-" + index

		fmt.Println("download ", ossId, url, fileName)
		resp, _ := http.Get(url)
		body, _ := ioutil.ReadAll(resp.Body)
		//out, _ := os.Create(fileName)
		//io.Copy(out, bytes.NewReader(body))

		uploadUrl := fmt.Sprintf(`%s/%s/%s`, UPLOAD_BASE, ossId, index)

		upload(uploadUrl, body)
		db.UpdatePhotoUrl(id, uploadUrl)
		time.Sleep(5 * time.Second)
	}

	fmt.Println("rows", rows)
}

func upload(url string, file []byte) {

	fmt.Println("upload ", url, len(file))

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile("uploadfile", "test.jpg")
	if err != nil {
		log.Fatalf("Create form file failed: %s\n", err)
	}

	_, err = io.Copy(formFile, bytes.NewReader(file))
	util.CheckError(err)

	// 发送表单
	contentType := writer.FormDataContentType()
	writer.Close() // 发送之前必须调用Close()以写入结尾行
	_, err = http.Post(url, contentType, buf)
	util.CheckError(err)

}
