package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"github.com/mateuszdyminski/bloom-filter/bloom"
)

// Create the Extender implementation, based on the gocrawl-provided DefaultExtender,
// because we don't want/need to override all methods.
type ExampleExtender struct {
	gocrawl.DefaultExtender
	bloom   *bloom.BloomFilter
	visited uint64
}

// Visit .
func (e *ExampleExtender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {

	// add entry
	e.bloom.AddString(ctx.URL().String())

	// increment counter
	atomic.AddUint64(&e.visited, 1)

	// find all images on the page
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			fmt.Printf("Page %s - img %s\n", ctx.URL().String(), src)
		}
	})

	fmt.Printf("Visited %s %d \n", ctx.URL().String(), e.visited)

	return nil, true
}

// Filter check in bloom filter if page was already seen or not.
func (e *ExampleExtender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	return !e.bloom.TestString(ctx.URL().String())
}

func main() {
	// Set custom options
	ext := new(ExampleExtender)
	ext.bloom = bloom.NewWithEstimates(10000000, 0.01)

	opts := gocrawl.NewOptions(ext)
	opts.CrawlDelay = 0 * time.Millisecond
	opts.LogFlags = gocrawl.LogNone
	opts.SameHostOnly = false
	opts.MaxVisits = 20000

	// Create crawler and start at root of duckduckgo
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run([]string{"http://gazeta.pl/", "http://onet.pl/", "http://wp.pl/"})

}
