package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

func GetPageSizes(urls []string) (map[string]int64, error) {
	results := make(map[string]int64)
	var wg sync.WaitGroup
	var requestErr error
	var mutex sync.Mutex

	for _, url := range urls {
		wg.Add(1)
		go func(urlToDownload string, mutex *sync.Mutex) {
			defer wg.Done()
			resp, err := http.Get(urlToDownload)
			defer mutex.Unlock()
			mutex.Lock()
			if err != nil {
				requestErr = err
				return
			}
			defer resp.Body.Close()
			_, err = io.ReadAll(resp.Body)
			if err != nil {
				requestErr = err
				return
			}
			fmt.Println(resp, resp.ContentLength)
			results[urlToDownload] = resp.ContentLength
		}(url, &mutex)
	}

	wg.Wait()
	return results, requestErr
}

func main() {
	inputUrls := []string{"https://v2.jokeapi.dev/joke/Any"}
	urlContentLen, err := GetPageSizes(inputUrls)
	if err != nil {
		fmt.Println("failed to download all of the urls", err)

	}
	fmt.Printf("urls and their lengths %v \n", urlContentLen)
}
