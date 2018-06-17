package main

import (
	"fmt"
	"log"

	httpclient "github.com/sniperkit/colly-request/plugin/client/unmarshal/v1"
)

func main() {
	content, err := httpclient.String("http://www.google.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("%s", content)
}
