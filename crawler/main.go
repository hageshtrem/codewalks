package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type Crawler struct {
	fetcher Fetcher
	pages   map[string]struct{}
	sync.RWMutex
}

func NewCrawler(fetcher Fetcher) *Crawler {
	return &Crawler{
		fetcher: fetcher,
		pages:   map[string]struct{}{},
		RWMutex: sync.RWMutex{},
	}
}

func (c *Crawler) isURLChecked(url string) bool {
	c.RLock()
	defer c.RUnlock()
	_, ok := c.pages[url]
	return ok
}

func (c *Crawler) addURL(url string) {
	c.Lock()
	defer c.Unlock()
	c.pages[url] = struct{}{}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func (c *Crawler) Crawl(url string, depth int, wg *sync.WaitGroup) {
	if depth <= 0 {
		return
	}
	c.addURL(url)
	body, urls, err := c.fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		if !c.isURLChecked(u) {
			wg.Add(1)
			go func(u string) {
				defer wg.Done()
				c.Crawl(u, depth-1, wg)
			}(u)
		}
	}
	return
}

func main() {
	c := NewCrawler(fetcher)
	wg := &sync.WaitGroup{}
	c.Crawl("https://golang.org/", 4, wg)
	wg.Wait()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
