package main

import (
	"fmt"
	"log"

	httpclient "github.com/sniperkit/colly-request/plugin/client/unmarshal/v1"
)

func main() {
	urls := []string{
		"http://www.golang.org",
		"http://www.clojure.org",
		"http://www.haskell.org",
	}
	var files []httpclient.File
	err := httpclient.Download(urls, &files)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Download files success!")
}
