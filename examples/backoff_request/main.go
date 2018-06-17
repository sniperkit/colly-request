package main

import (
	"fmt"

	request "github.com/sniperkit/colly-request/pkg"
	http_backoff "github.com/sniperkit/colly-request/plugin/backoff"
)

const wrapperIdentifier = "xcolly-backoff-request"

func main() {
	fmt.Println("Running '%s' example...", wrapperIdentifier)
}
