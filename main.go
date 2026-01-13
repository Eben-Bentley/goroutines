package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// processTask performs CPU-intensive hashing on the task data.
// This simulates real-world operations like file hashing, encryption, etc.
func processTask(task Task) Result {
	data := []byte(task.Data)
	for i := 0; i < 1000; i++ {
		hash := sha256.Sum256(data)
		data = hash[:]
	}
	return Result{
		ID:   task.ID,
		Hash: hex.EncodeToString(data),
	}
}

// processSequential processes tasks one at a time in order.
// Simple and predictable, but doesn't use multiple CPU cores.
func processSequential(tasks []Task) []Result {
	results := make([]Result, len(tasks))
	for i, task := range tasks {
		results[i] = processTask(task)
	}
	return results
}

// processWithGoroutines processes tasks concurrently using a worker pool pattern.
// Distributes work across multiple goroutines to utilize multiple CPU cores.
func processWithGoroutines(tasks []Task, numWorkers int) []Result {
	jobs := make(chan Task, len(tasks))
	results := make(chan Result, len(tasks))

	// Start worker pool
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range jobs {
				results <- processTask(task)
			}
		}()
	}

	// Send all tasks to workers
	for _, task := range tasks {
		jobs <- task
	}
	close(jobs)

	// Wait for all workers to finish, then close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results
	resultSlice := make([]Result, 0, len(tasks))
	for result := range results {
		resultSlice = append(resultSlice, result)
	}

	return resultSlice
}

func main() {
	numCPU := runtime.NumCPU()
	numTasks := 100

	// Create tasks
	tasks := make([]Task, numTasks)
	for i := 0; i < numTasks; i++ {
		tasks[i] = Task{
			ID:   i + 1,
			Data: fmt.Sprintf("Data item %d", i+1),
		}
	}

	fmt.Printf("Goroutines Demo - Processing %d tasks\n", numTasks)
	fmt.Printf("CPU Cores: %d\n\n", numCPU)

	// Sequential processing
	start := time.Now()
	resultsSeq := processSequential(tasks)
	durationSeq := time.Since(start)
	fmt.Printf("Sequential:  %v  (%.0f items/sec)\n", durationSeq, float64(len(resultsSeq))/durationSeq.Seconds())

	// Concurrent processing with different worker counts
	workerCounts := []int{2, 4, numCPU}
	for _, workers := range workerCounts {
		start := time.Now()
		resultsConc := processWithGoroutines(tasks, workers)
		durationConc := time.Since(start)
		speedup := float64(durationSeq) / float64(durationConc)
		fmt.Printf("%d workers:   %v  (%.0f items/sec) - %.2fx faster\n",
			workers, durationConc, float64(len(resultsConc))/durationConc.Seconds(), speedup)
	}
}
