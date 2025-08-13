package main

import "sync"

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}
type fetchedURLs struct {
	bodies         []string
	checkUniqueURL map[string]bool
	sync.RWMutex
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) ([]string, error) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		return nil, nil
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		return nil, err
	}
	result := &fetchedURLs{
		bodies:         []string{},
		checkUniqueURL: make(map[string]bool),
	}
	result.bodies = append(result.bodies, body)
	result.checkUniqueURL[url] = true
	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		// if res, err := Crawl(u, depth-1, fetcher); err == nil {
		// 	result = append(result, res...)
		// }
		go crawl(u, depth-1, fetcher, &wg, result)
	}
	wg.Wait()
	return result.bodies, nil
}
func crawl(url string, depth int, fetcher Fetcher, wg *sync.WaitGroup, result *fetchedURLs) {
	defer wg.Done()
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}
	result.Lock()
	if isFetched := result.checkUniqueURL[url]; isFetched {
		result.Unlock()
		return
	} else {
		result.checkUniqueURL[url] = true
		result.bodies = append(result.bodies, body)
		result.Unlock()
	}

	for _, u := range urls {
		wg.Add(1)
		// if res, err := Crawl(u, depth-1, fetcher); err == nil {
		// 	result = append(result, res...)
		// }
		go crawl(u, depth-1, fetcher, wg, result)
	}
}
