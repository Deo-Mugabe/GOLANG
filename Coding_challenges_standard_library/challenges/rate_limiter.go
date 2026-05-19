// We have a system that receives 10,000 Job per second via a Kafka topic. 
// Each job requires calling a 3rd-party API that has a strict rate limit of 500 concurrent requests. 
// If we exceed this, we get banned. If we go too slow, our Kafka consumer lag grows too large.
// Task: Write a Go implementation for a worker pool that:

// Consumes from a channel
// Limits active executions to exactly 500
// Handles a context timeout of 2 seconds per API call
// Logs the total number of successful vs. failed jobs after the channel is closed

package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

type Job struct {
    ID int
}

func callAPI(ctx context.Context, job Job) error {
    // simulate API call
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(100 * time.Millisecond):
        return nil
    }
}

func main() {
    jobs := make(chan Job, 1000)
    sem := make(chan struct{}, 500) // limit to 500 concurrent

    var (
        wg      sync.WaitGroup
        mu      sync.Mutex
        success int
        failed  int
    )

    // simulate Kafka consumer feeding jobs
    go func() {
        for i := 0; i < 10000; i++ {
            jobs <- Job{ID: i}
        }
        close(jobs)
    }()

    for job := range jobs {
        sem <- struct{}{} // acquire slot
        wg.Add(1)

        go func(j Job) {
            defer wg.Done()
            defer func() { <-sem }() // release slot

            ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
            defer cancel()

            err := callAPI(ctx, j)

            mu.Lock()
            if err != nil {
                failed++
            } else {
                success++
            }
            mu.Unlock()
        }(job)
    }

    wg.Wait()
    fmt.Printf("Success: %d, Failed: %d\n", success, failed)
}