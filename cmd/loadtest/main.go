package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/lmtani/learning-go-loadtest/internal/executor"
)

func main() {
	// Parse command line flags
	url := flag.String("url", "", "URL of the service to test")
	requests := flag.Int("requests", 100, "Total number of requests to perform")
	concurrency := flag.Int("concurrency", 1, "Number of concurrent requests")
	flag.Parse()

	if *url == "" {
		log.Fatal("URL must be provided")
	}

	// Display a welcome message
	fmt.Println("Starting load test...")
	fmt.Printf("URL: %s\n", *url)
	fmt.Printf("Requests: %d\n", *requests)
	fmt.Printf("Concurrency: %d\n", *concurrency)
	fmt.Println("Please wait...")

	// Execute the load test
	report, err := executor.ExecuteLoadTest(*url, *requests, *concurrency)
	if err != nil {
		log.Fatalf("Error executing load test: %v", err)
	}

	// Render the report
	report.Render()
}
