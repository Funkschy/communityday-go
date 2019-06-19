package main

import (
	"github.com/funkschy/communityday/pkg/url"
	"log"
)

func main() {
	raw, err := url.ReadURLFromCommandLine()
	if err != nil {
		log.Fatal(err)
	}

	crawler := NewCrawler()
	crawler.Crawl(raw)
}
