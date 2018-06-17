package main

import (
	"fmt"
	"log"

	httpclient "github.com/sniperkit/colly-request/plugin/client/unmarshal/v1"
)

func main() {
	content, err := httpclient.Bytes("http://www.example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("%#v", content)
}
