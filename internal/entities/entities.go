package entities

import (
	"fmt"
	"time"
)

type Report struct {
	TotalRequests          int
	SuccessfulRequests     int
	StatusCodeDistribution map[int]int
	TotalTime              time.Duration
}

func (r *Report) Render() {
	fmt.Println("Load Test Report:")
	fmt.Println("-----------------")
	fmt.Printf("Total Requests: %d\n", r.TotalRequests)
	fmt.Printf("Successful Requests (HTTP 200): %d\n", r.SuccessfulRequests)
	fmt.Printf("Total Time Taken: %s\n", r.TotalTime)
	fmt.Println("Status Code Distribution:")
	for code, count := range r.StatusCodeDistribution {
		fmt.Printf("  HTTP %d: %d\n", code, count)
	}
}
