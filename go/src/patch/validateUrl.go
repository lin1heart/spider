package patch

import (
	"crypto/tls"
	"fmt"
	"github.com/lin1heart/spider/go/src/db"
	"github.com/lin1heart/spider/go/src/util"
	"net/http"
	"strconv"
	"time"
)

func init() {
	fmt.Println("patch init")
}
func main() {
	startId := 4600
	endId := 244008

	inteval := 1000

	for index := startId; index < endId; {
		fmt.Println("index ", index)

		results := db.QueryPhotoAfterRow(startId, inteval)

		for _, row := range results {
			id, err := strconv.Atoi(row["id"])
			util.CheckError(err)
			valid := validate(row["url"], id)
			if valid == false {
				db.UpdateUrlToEmpty(id)
			}
			fmt.Printf("%d valid %v \n", id, valid)
			if id >= startId {
				startId = id
				fmt.Println("new startid ", startId)
			}

		}
		index += inteval
		fmt.Println("delay 10 s with id: ", startId)
		time.Sleep(time.Second * 10)
	}
	fmt.Println("succeed")
}

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//DisableKeepAlives: true,
}
var client = &http.Client{Transport: tr}

func validate(url string, id int) bool {
	resp, err := client.Get(url)
	util.CheckError(err)

	defer func () {
		err := resp.Body.Close()
		util.CheckError(err)
	} ()

	if 404 == resp.StatusCode {
		fmt.Printf("id %d 404 url %s \n", id, url)
		return false
	}
	return true
}
