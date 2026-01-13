package main

import (
	"log"

	"github.com/jansaidl/meilisearch-scraper/src/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
