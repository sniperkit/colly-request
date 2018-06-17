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
	http_cache "github.com/sniperkit/colly-request/plugin/cache"
)

const wrapperIdentifier = "xcolly-cached-request"

func main() {
	fmt.Printf("Running '%s' example...\n", wrapperIdentifier)
}
