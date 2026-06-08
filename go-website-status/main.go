package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	targetURL = "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN"
)

type websites struct {
	Name     string `json:"name"`
	Interval int    `json:"interval"`
}

type config struct {
	Urls []websites `json:"websites"`
}

type WebhookPayload struct {
	Contents string `json:"content"`
}

func sendDiscordNotifications(ctx context.Context, alerts chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	// Create the client ONCE outside the loop to reuse TCP connections efficiently
	client := &http.Client{}

	for {
		select {
		case message := <-alerts:
			// 1. Prepare Payload
			payload := WebhookPayload{Contents: message}
			jsonBytes, err := json.Marshal(payload)
			if err != nil {
				fmt.Println("failed to marshal payload", err)
				continue // Use continue instead of return so the worker doesn't die!
			}

			// 2. Build Request
			req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(jsonBytes))
			if err != nil {
				fmt.Println("failed to make a http request", err)
				continue
			}
			req.Header.Set("Content-Type", "application/json")

			// 3. Apply Timeout Context (Notice we use a new variable name: reqCtx)
			reqCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			req = req.WithContext(reqCtx)

			// 4. Execute Request
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("failed to send webhook: %v\n", err)
				cancel() // Clean up the timeout context if it fails early
				continue
			}

			// 5. Evaluate Response
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				fmt.Printf("received non-success status code from Discord: %d\n", resp.StatusCode)
			} else {
				fmt.Println("Webhook sent successfully!")
			}

			// CRITICAL: Explicitly close and cancel here! 
			// Do not use `defer` inside an infinite loop or they will leak memory.
			resp.Body.Close()
			cancel()

		case <-ctx.Done():
			fmt.Println("Shutting down the notification worker...")
			return
		}
	}
}
func ping(url string) bool {

	resp, err := http.Head(url)
	if err != nil {
		fmt.Println("failed to get a response for the url", url, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		fmt.Printf("bad status code from website %s: %d\n", url, resp.StatusCode)
		return false
	}
	fmt.Println("received success response for the url", url)
	return true
}

func checkStatus(ctx context.Context, url string, interval int, wg *sync.WaitGroup, alerts chan string) {

	defer wg.Done()
	duration := time.Duration(interval) * time.Second
	tick := time.NewTicker(duration)
	defer tick.Stop()
	stateDown := false
	pingCount := 0
	for {

		if !ping(url) {
			fmt.Println("failed to get a ping reply for the site", url)
			pingCount += 1
		} else {
			stateDown = false
			fmt.Println("got a ping reply for the site", url)
			pingCount = 0
		}

		if pingCount > 3 && !stateDown {
			alerts <- fmt.Sprintf("failed to get response for ping from url %s", url)
			stateDown = true
		}
		select {
		case <-ctx.Done():
			fmt.Println("Worker stopping...")
			return
		case <-tick.C:
			break
		}
	}
}
func main() {

	configFile := "./configs/config.json"
	fileContents, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("failed to read the contents of the file", err)
		os.Exit(1)
	}
	configObj := config{}
	err = json.Unmarshal(fileContents, &configObj)
	if err != nil {
		fmt.Println("failed to parse the config file", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	alerts := make(chan string, 100)
	wg.Add(1)
	go sendDiscordNotifications(ctx, alerts, &wg)
	
	for _, website := range configObj.Urls {
		interval := 5
		if website.Interval != 0 {
			fmt.Println("checking website status every discord", interval, " seconds")
			interval = website.Interval
		}
		fmt.Println("fist url to process", website.Name)
		wg.Add(1)
		go checkStatus(ctx, website.Name, interval, &wg, alerts)
	}

	time.Sleep(30 * time.Second)
	fmt.Println("stopping the server")
	cancel()

	fmt.Println("waiting for all the go-routines to gracefully finish")
	wg.Wait()
	fmt.Println("all-go-routines gracefully ended")
}
