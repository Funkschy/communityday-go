package main

import (
	"github.com/funkschy/communityday/pkg/url"
	"log"
	"net/http"
	"sync"
)

type urlPair struct {
	base string
	link string
}

// Crawler besucht das gesamte Internet
type Crawler struct {
	sync.Map
	http.Client
	// Die URLs, die noch besucht werden sollen
	toVisit chan urlPair
}

// NewCrawler erstellt einen neuen Crawler
func NewCrawler() Crawler {
	return Crawler{
		sync.Map{},
		url.Client(),
		make(chan urlPair, 1)}
}

// Crawl startet von einer URL und durchsucht das gesamte Internet
func (c *Crawler) Crawl(startURL string) {
	go func() { c.toVisit <- urlPair{"", startURL} }()

	for linkPair := range c.toVisit {
		go c.visit(linkPair.base, linkPair.link)
	}
}

func (c *Crawler) visit(base string, link string) {
	uri, err := url.ParseURL(base, link)
	if err != nil {
		log.Printf("Could not parse link %s of base %s", link, base)
		return
	}

	uriString := uri.String()
	if _, ok := c.Load(uriString); ok {
		// url ist schon besucht
		return
	}
	c.Store(uriString, true)

	resp, err := c.Get(uriString)
	if err != nil {
		log.Printf("Could not get: %s\n", err)
		return
	}

	links, err := url.FetchLinks(resp.Body)
	resp.Body.Close()

	if err != nil {
		log.Printf("Html of %s is malformed\n", uriString)
		return
	}

	log.Printf("Visited %s", uriString)

	for _, nextLink := range links {
		c.toVisit <- urlPair{uriString, nextLink}
	}
}
