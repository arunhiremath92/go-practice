package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type websites struct {
	Name string `json:"name"`
}

type config struct {
	Urls []websites `json:"websites"`
}

func ping(url string) bool {
	return true
}

func checkStatus(url string, interval int, done chan struct{}, wg *sync.WaitGroup) {

	defer wg.Done()
	for {

		if !ping(url) {
			fmt.Println("failed to get a ping reply for the site", url)
		} else {
			fmt.Println("got a ping reply for the site", url)
		}

		duration := time.Duration(interval) * time.Second
		tick := time.Tick(duration)
		select {
		case <-done:
			fmt.Println("Worker stopping...")
			return
		case <-tick:
			break
		}
	}
}
func main() {

	configFile := "/Users/arunhiremath/workspace/go-website-status/configs/config.json"
	fileContents, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("failed to read the contents of the file", err)
	}
	configObj := config{}
	err = json.Unmarshal(fileContents, &configObj)
	if err != nil {
		fmt.Println("failed to parse the config file", err)
	}
	done := make(chan struct{})
	var wg sync.WaitGroup
	for _, website := range configObj.Urls {
		fmt.Println("fist url to process", website.Name)
		wg.Add(1)
		go checkStatus(website.Name, 5, done, &wg)
	}

	time.Sleep(30 * time.Second)
	fmt.Println("stopping the server")
	close(done)

	fmt.Println("waiting for all the go-routines to gracefully finish")
	wg.Wait()
	fmt.Println("all-go-routines gracefully ended")
}
