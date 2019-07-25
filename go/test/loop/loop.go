package loop

import (
	"crypto/tls"
	"fmt"
	"github.com/lin1heart/spider/go/src/util"
	"net/http"
	"os"
)


func main() {
	length := 10000

	fmt.Println("pid is ", os.Getpid())
	for index:=1; index < length; index ++{
		valid := validate("http://localhost:8888/1024px-Bitcoin.svg.png", index)
		fmt.Printf("%d is valid %v;", index, valid)
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

	defer func () {
		err := resp.Body.Close()
		util.CheckError(err)
	} ()

	if 200 != resp.StatusCode {
		fmt.Printf("id %d 404 url %s \n", id, url)
		return false
	}
	return true
}