package entities

import (
	"fmt"
	"time"
)

type Report struct {
	// Request statistics
	TotalRequests      int
	SuccessfulRequests int
	FailedRequests     int
	TotalTime          time.Duration

	// Status code information
	StatusCodeDistribution map[int]int

	// Response time statistics
	ResponseTimes   []time.Duration
	MinResponseTime time.Duration
	MaxResponseTime time.Duration
	AvgResponseTime time.Duration
	P50ResponseTime time.Duration
	P95ResponseTime time.Duration
	P99ResponseTime time.Duration
}

func (r *Report) Render() {
	// The actual rendering is delegated to the UI package
	// This is just a wrapper that prints the rendered output
	fmt.Println(r.RenderString())
}

// RenderString returns the report as a formatted string
// This allows the UI package to use it without circular imports
func (r *Report) RenderString() string {
	// This will be replaced by the UI package's implementation
	// We keep a simple version here for fallback or testing
	var result string
	result += "Load Test Report:\n"
	result += "-----------------\n"
	result += fmt.Sprintf("Total Requests: %d\n", r.TotalRequests)
	result += fmt.Sprintf("Successful Requests (HTTP 200): %d\n", r.SuccessfulRequests)
	result += fmt.Sprintf("Total Time Taken: %s\n", r.TotalTime)
	result += "Status Code Distribution:\n"
	for code, count := range r.StatusCodeDistribution {
		result += fmt.Sprintf("  HTTP %d: %d\n", code, count)
	}
	return result
}
