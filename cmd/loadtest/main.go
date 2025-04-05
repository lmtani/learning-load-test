package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/lmtani/learning-go-loadtest/internal/executor"
	"github.com/lmtani/learning-go-loadtest/internal/ui"
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

	// Create a channel for progress updates
	progressCh := make(chan executor.ProgressUpdate)

	// Start a goroutine to display progress
	go func() {
		for progress := range progressCh {
			ui.ClearProgressBar()

			// Render and print the progress bar
			fmt.Println(ui.RenderProgressBar(progress))
		}
	}()

	// Execute the load test with progress reporting
	report, err := executor.ExecuteLoadTest(*url, *requests, *concurrency, progressCh)
	if err != nil {
		log.Fatalf("Error executing load test: %v", err)
	}

	// Close the progress channel
	close(progressCh)

	// Sleep briefly to ensure the progress display is complete
	time.Sleep(200 * time.Millisecond)

	// Clear the progress bar before showing the final report
	ui.ClearProgressBar()

	// Render the final report
	fmt.Println(ui.RenderReport(report))
}
