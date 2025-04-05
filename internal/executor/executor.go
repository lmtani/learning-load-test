package executor

import (
	"io"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/lmtani/learning-go-loadtest/internal/entities"
)

// ProgressUpdate contains information about the current progress of the load test
type ProgressUpdate struct {
	CompletedRequests int
	TotalRequests     int
	ElapsedTime       time.Duration
}

// ExecuteLoadTest performs a load test against the specified URL
// If progressCh is not nil, it will send progress updates through the channel
func ExecuteLoadTest(url string, totalRequests, concurrency int, progressCh chan<- ProgressUpdate) (*entities.Report, error) {
	var wg sync.WaitGroup
	startTime := time.Now()
	var completedRequests int

	// Initialize the report with enhanced fields
	result := &entities.Report{
		TotalRequests:          totalRequests,
		SuccessfulRequests:     0,
		StatusCodeDistribution: make(map[int]int),
		ResponseTimes:          make([]time.Duration, 0, totalRequests),
		MinResponseTime:        time.Hour, // Start with a large value
		MaxResponseTime:        0,
		AvgResponseTime:        0,
	}

	// Channel for limiting concurrency
	sem := make(chan struct{}, concurrency)
	var resultMutex sync.Mutex

	// For tracking progress
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	// Start a goroutine to report progress
	if progressCh != nil {
		go func() {
			for range ticker.C {
				resultMutex.Lock()
				completed := completedRequests
				resultMutex.Unlock()

				select {
				case progressCh <- ProgressUpdate{
					CompletedRequests: completed,
					TotalRequests:     totalRequests,
					ElapsedTime:       time.Since(startTime),
				}:
				default:
					// Non-blocking send to prevent slowdowns if channel is not being read
				}

				if completed >= totalRequests {
					return
				}
			}
		}()
	}

	// Execute requests
	for i := 0; i < totalRequests; i++ {
		sem <- struct{}{}
		wg.Add(1)

		go func() {
			defer func() {
				<-sem
				wg.Done()
			}()

			// Measure individual request time
			requestStart := time.Now()

			// Execute the request
			resp, err := http.Get(url)
			if err != nil {
				resultMutex.Lock()
				result.FailedRequests++
				resultMutex.Unlock()
				return
			}

			// Read and discard the body to ensure connection is properly closed
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			// Calculate response time
			responseTime := time.Since(requestStart)

			resultMutex.Lock()
			defer resultMutex.Unlock()

			// Increment completed requests counter for progress tracking
			completedRequests++

			// Update status code statistics
			if resp.StatusCode == http.StatusOK {
				result.SuccessfulRequests++
			}
			result.StatusCodeDistribution[resp.StatusCode]++

			// Update response time statistics
			result.ResponseTimes = append(result.ResponseTimes, responseTime)

			// Update min/max response times
			if responseTime < result.MinResponseTime {
				result.MinResponseTime = responseTime
			}
			if responseTime > result.MaxResponseTime {
				result.MaxResponseTime = responseTime
			}
		}()
	}

	wg.Wait()
	endTime := time.Now()
	result.TotalTime = endTime.Sub(startTime)

	// Calculate average response time
	if len(result.ResponseTimes) > 0 {
		var totalResponseTime time.Duration
		for _, t := range result.ResponseTimes {
			totalResponseTime += t
		}
		result.AvgResponseTime = totalResponseTime / time.Duration(len(result.ResponseTimes))
	}

	// Calculate percentiles
	if len(result.ResponseTimes) > 0 {
		// Sort response times
		times := make([]time.Duration, len(result.ResponseTimes))
		copy(times, result.ResponseTimes)

		// Simple bubble sort (for small datasets)
		for i := 0; i < len(times); i++ {
			for j := i + 1; j < len(times); j++ {
				if times[i] > times[j] {
					times[i], times[j] = times[j], times[i]
				}
			}
		}

		// Calculate percentiles
		p50Index := int(math.Floor(float64(len(times)) * 0.5))
		p95Index := int(math.Floor(float64(len(times)) * 0.95))
		p99Index := int(math.Floor(float64(len(times)) * 0.99))

		result.P50ResponseTime = times[p50Index]
		result.P95ResponseTime = times[p95Index]
		result.P99ResponseTime = times[p99Index]
	}

	// Calculate failed requests
	result.FailedRequests = totalRequests - result.SuccessfulRequests

	return result, nil
}
