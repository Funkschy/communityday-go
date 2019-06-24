package main

import (
	"github.com/funkschy/communityday/pkg/url"
	"log"
	"net/http"
	u "net/url"
	"sync"
)

const numWorker = 16

type urlPair struct {
	base string
	link string
}

// Crawler besucht das gesamte Internet
type Crawler struct {
	sync.Map
	http.Client
	results chan *u.URL
	jobs    chan *u.URL
}

// NewCrawler erstellt einen neuen Crawler
func NewCrawler() Crawler {
	results := make(chan *u.URL)
	jobs := make(chan *u.URL, 16)

	return Crawler{
		sync.Map{},
		url.Client(),
		results,
		jobs,
	}
}

// Crawl startet von einer URL und durchsucht das gesamte Internet
func (c *Crawler) Crawl(startURL string) {
	start, _ := url.ParseURL("", startURL)
	c.jobs <- start

	c.initThreadPool(numWorker)

	for url := range c.results {
		urlString := url.String()
		if c.isVisited(urlString) {
			continue
		}
		c.setVisited(urlString)
		log.Printf("Visited %s", url.String())

		go func(url *u.URL) { c.jobs <- url }(url)
	}
}

func (c *Crawler) isVisited(url string) bool {
	_, exists := c.Load(url)
	return exists
}

func (c *Crawler) setVisited(url string) {
	c.Store(url, true)
}

func (c *Crawler) initThreadPool(numWorkers int) {
	for i := 0; i < numWorker; i++ {
		go func() {
			for job := range c.jobs {
				c.work(job)
			}
		}()
	}
}

func (c *Crawler) work(uri *u.URL) {
	uriString := uri.String()

	for _, nextLink := range c.visit(uri) {
		uri, err := url.ParseURL(uriString, nextLink.String())
		if err != nil {
			log.Printf("Could not parse link %s of base %s", nextLink, uriString)
			return
		}

		c.results <- uri
	}
}

func (c *Crawler) visit(uri *u.URL) []*u.URL {
	uriString := uri.String()
	emptyList := make([]*u.URL, 0, 0)

	resp, err := c.Get(uriString)
	if err != nil {
		log.Printf("Could not get: %s\n", err)
		return emptyList
	}

	links, err := url.FetchLinks(resp.Body)
	resp.Body.Close()

	if err != nil {
		log.Printf("Html of %s is malformed\n", uriString)
		return emptyList
	}

	numLinks := len(links)
	list := make([]*u.URL, numLinks, numLinks)

	for i, link := range links {
		uri, err := url.ParseURL(uriString, link)
		if err != nil {
			log.Printf("Could not parse link %s of base %s", link, uriString)
			return emptyList
		}
		list[i] = uri
	}

	return list
}
