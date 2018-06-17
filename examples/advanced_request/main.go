package main

import (
	"fmt"

	// external
	storage "github.com/sniperkit/colly-storage/pkg"
	bck_badger "github.com/sniperkit/colly-storage/plugin/backend/badger"
	bck_boltdb "github.com/sniperkit/colly-storage/plugin/backend/boltdb"
	bck_bboltdb "github.com/sniperkit/colly-storage/plugin/backend/boltdb_bbolt"

	// internal
	request "github.com/sniperkit/colly-request/pkg"
	http_backoff "github.com/sniperkit/colly-request/plugin/backoff"
	http_cache "github.com/sniperkit/colly-request/plugin/cache"
	http_stats "github.com/sniperkit/colly-request/plugin/stats"
)

const wrapperIdentifier = "xcolly-advanced-request"

func main() {
	fmt.Printf("Running '%s'...\n", wrapperIdentifier)
}
