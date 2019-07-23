package http

import (
	"crypto/tls"
	"fmt"
	"github.com/lin1heart/spider/go/src/util"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func main() {

	rand.Seed(time.Now().Unix())
	fmt.Println("proxy is ", rand.Intn(5))

	str := util.RandomProxy()
	fmt.Println("proxyUrl:", str)
	proxyUrl, err := url.Parse(str)
	util.CheckError(err)

	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(proxyUrl),
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get("http://util.online/headers")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println("body", string(body))
}
