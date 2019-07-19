package photo

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/lin1heart/spider/go/src/db"
	"github.com/lin1heart/spider/go/src/util"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

func Sync() {
	arr1 := []interface{}{util.PhotoType["PHOTO_WEST"], util.PhotoType["PHOTO_EAST"]}
	arr2 := []interface{}{util.PhotoType["PHOTO_UNIFORM"], util.PhotoType["RANK"]}
	arr3 := []interface{}{util.PhotoType["PHOTO_PURE"], util.PhotoType["PHOTO_SELF"]}
	arr4 := []interface{}{util.PhotoType["PHOTO_COMIC"], util.PhotoType["PHOTO_RAPE"]}
	go loopTask(arr1)
	go loopTask(arr2)
	go loopTask(arr3)
	loopTask(arr4)

}
func loopTask(imageTypes []interface{}) {
	for true {
		go task(imageTypes)
		time.Sleep(10 * time.Second)
	}
}

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
var client = &http.Client{Transport: tr}

func task(imageTypes []interface{}) {

	rows := db.QueryEmptyPhotos(imageTypes)

	for _, row := range rows {

		url := row["crawl_url"]
		title := row["title"]
		index := row["index"]
		ossId := row["oss_id"]
		idString := row["id"]

		id, err := strconv.Atoi(idString)
		if err != nil {
			fmt.Println("sync task strconv.Atoi err", idString, err, row)
			continue
		}

		fileName := title + "-" + index

		fmt.Println("download ", ossId, url, fileName)

		resp, err := client.Get(url)
		defer resp.Body.Close()

		if err != nil {
			fmt.Println("client get  err", err)
			time.Sleep(5 * time.Second)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		util.CheckError(err)
		//out, _ := os.Create(fileName)
		//io.Copy(out, bytes.NewReader(body))

		uploadUrl := fmt.Sprintf(`%s/%s/%s`, util.UPLOAD_BASE, ossId, index)

		err = upload(uploadUrl, body)
		if err != nil {
			fmt.Println("upload err", err)
			time.Sleep(5 * time.Second)
			continue
		}

		db.UpdatePhotoUrl(id, uploadUrl)
		time.Sleep(2 * time.Second)
	}

	fmt.Println("rows", rows)
}

func upload(url string, file []byte) error {

	fmt.Println("upload ", url, len(file))

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile("uploadfile", "test.jpg")
	util.CheckError(err)

	_, err = io.Copy(formFile, bytes.NewReader(file))
	util.CheckError(err)

	// 发送表单
	contentType := writer.FormDataContentType()
	writer.Close() // 发送之前必须调用Close()以写入结尾行
	resp, err := http.Post(url, contentType, buf)
	defer resp.Body.Close()

	return err
}
