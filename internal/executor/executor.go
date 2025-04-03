package executor

import (
	"net/http"
	"sync"
	"time"

	"github.com/lmtani/learning-go-loadtest/internal/entities"
)

func ExecuteLoadTest(url string, totalRequests, concurrency int) (*entities.Report, error) {
	var wg sync.WaitGroup
	startTime := time.Now()
	result := &entities.Report{
		TotalRequests:          totalRequests,
		SuccessfulRequests:     0,
		StatusCodeDistribution: make(map[int]int),
	}

	sem := make(chan struct{}, concurrency)
	var resultMutex sync.Mutex

	for range totalRequests {
		sem <- struct{}{}
		wg.Add(1)

		go func() {
			defer func() {
				<-sem
			}()

			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			resultMutex.Lock()
			defer resultMutex.Unlock()

			if resp.StatusCode == http.StatusOK {
				result.SuccessfulRequests++
			}
			result.StatusCodeDistribution[resp.StatusCode]++
			wg.Done()
		}()
	}

	wg.Wait()
	endTime := time.Now()
	result.TotalTime = endTime.Sub(startTime)
	return result, nil
}
