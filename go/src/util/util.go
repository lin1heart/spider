package util

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

var args = os.Args
var ENV = os.Getenv("ENV")
var ENTRY = ""
var IMG_ENDPOINT = os.Getenv("IMG_ENDPOINT")
var DB_HOST = os.Getenv("DB_HOST")
var DB_PORT = os.Getenv("DB_PORT")
var Pid = os.Getpid()

func init() {
	fmt.Println("util init")
	if len(args) >= 2 {
		ENTRY = args[1]
	}
	osUpload := os.Getenv("IMG_ENDPOINT")
	if osUpload != "" {
		IMG_ENDPOINT = osUpload
	}
	if DB_HOST == "" {
		DB_HOST = "218.168.168.105"
	}
	if DB_PORT == "" {
		DB_PORT = "3306"
	}
	if IMG_ENDPOINT == "" {
		IMG_ENDPOINT = "http://218.168.168.105:8888"
	}

	fmt.Println("ProxyList:", ProxyList)
	fmt.Println("ENV:", ENV)
	fmt.Println("ENTRY:", ENTRY)
	fmt.Println("Pid:", Pid)
}

func CheckError(errMasg error) {
	if errMasg != nil {
		fmt.Println("error %s", errMasg)
		//log.Fatal(errMasg)
		panic(errMasg)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenerateSqlIdsSuffix(ids []string) string {
	idsStr := strings.Join(ids, ",")

	if idsStr == "" {
		return ""
	}
	sqlSuffix := fmt.Sprintf(" AND id NOT IN ( %s )", idsStr)
	return sqlSuffix
}

func FilterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			new_content += string(value)
		}
	}
	return new_content
}

var ProxyList = []string{
	//"socks5://127.0.0.1:1080",
	//"socks5://47.96.123.41:1080",   // lin
	//"socks5://47.101.60.252:1080",  // pawl
	//"socks5://39.104.226.149:1080", // jj
	//"socks5://47.96.123.41:1080",   // lin
	//"socks5://3.0.176.116:1080",    // aws jj
	//"socks5://3.113.16.157:1080",   // aws lin

}

func RandomProxy() string {
	rand.Seed(time.Now().Unix())
	return ProxyList[rand.Intn(len(ProxyList))]
}

var PhotoType = map[string]string{
	"PHOTO_EAST":    "PHOTO_EAST",
	"PHOTO_WEST":    "PHOTO_WEST",
	"PHOTO_PURE":    "PHOTO_PURE",
	"PHOTO_SELF":    "PHOTO_SELF",
	"PHOTO_UNIFORM": "PHOTO_UNIFORM",
	"PHOTO_RAPE":    "PHOTO_RAPE",
	"PHOTO_COMIC":   "PHOTO_COMIC",
	"PHOTO_RANK":    "PHOTO_RANK",
}
