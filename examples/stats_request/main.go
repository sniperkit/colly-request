package main

import (
	"fmt"

	request "github.com/sniperkit/colly-request/pkg"
	http_stats "github.com/sniperkit/colly-request/plugin/stats"
)

const wrapperIdentifier = "xcolly-stats-request"

func main() {
	fmt.Printf("Running '%s' example...\n", wrapperIdentifier)
}
