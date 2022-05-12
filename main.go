package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func GetCountOccurrences(url string) (int64, error) {
	client := http.Client{
		Timeout: time.Duration(time.Second *  10),
	}

	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	count := strings.Count(string(body), "Go")
	return int64(count), nil
}

func main() {

	var totalCount int64 = 0

	urls := make([]string, 0)
	urls = append(urls, "https://golang.org")
	urls = append(urls, "https://golang.org")
	urls = append(urls, "https://golang.org")
	urls = append(urls, "https://golang.org")
	urls = append(urls, "https://golang.org")
	urls = append(urls, "https://golang.org")

	wg := sync.WaitGroup{}
	limit := make(chan struct{}, 5)

	log.Println("Start counting")

	for _, url := range urls {
		wg.Add(1)
		limit <- struct{}{}
		url := url
		go func() {
			defer func() {
				wg.Done()
				<-limit
			}()
			count, err := GetCountOccurrences(url)
			if err != nil {
				log.Printf("Getting count occurrences string  %v failed: %v\n", url, err)
				return
			}
			log.Printf("Count for %v: %v \n", url, count)
			atomic.AddInt64(&totalCount, count)
		}()
	}
	wg.Wait()

	log.Printf("Total %v", totalCount)
}
