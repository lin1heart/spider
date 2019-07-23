package patch

import (
	"crypto/tls"
	"fmt"
	"github.com/lin1heart/spider/go/src/db"
	"github.com/lin1heart/spider/go/src/util"
	"net/http"
	"strconv"
)

func init() {
	fmt.Println("patch init")
}
func main() {
	startId := 118154
	endId := 233648

	inteval := 100

	for index := startId; index < endId; {
		fmt.Println("index ", index)

		results := db.QueryPhotoAfterRow(index, inteval)

		for _, row := range results {
			id, err := strconv.Atoi(row["id"])
			util.CheckError(err)
			valid := validate(row["url"], id)
			if valid == false {
				db.UpdateUrlToEmpty(id)
			}
		}
		index += inteval
	}
}

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//DisableKeepAlives: true,
}
var client = &http.Client{Transport: tr}

func validate(url string, id int) bool {

	resp, err := client.Get(url)
	util.CheckError(err)

	//defer func () {
	//	err := resp.Body.Close()
	//	util.CheckError(err)
	//} ()
	fmt.Printf("stauts code %d \n", resp.StatusCode)
	if 404 == resp.StatusCode {
		fmt.Printf("id %d 404 url %s \n", id, url)
		return false
	}
	return true
}
