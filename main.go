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

	// Create tasks using range over integer (Go 1.22+)
	tasks := make([]Task, numTasks)
	for i := range numTasks {
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

	fmt.Println("\n--- Mutex Demo ---")
	demoMutex()

	fmt.Println("\n--- Semaphore Demo ---")
	demoSemaphore()
}

// BankAccount represents a simple bank account with thread-safe operations
type BankAccount struct {
	balance int
	mu      sync.Mutex
}

func (b *BankAccount) Deposit(amount int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.balance += amount
}

func (b *BankAccount) Withdraw(amount int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.balance >= amount {
		b.balance -= amount
		return true
	}
	return false
}

func (b *BankAccount) Balance() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.balance
}

// demoMutex simulates concurrent bank transactions.
// Without mutex protection, balance would become incorrect due to race conditions.
func demoMutex() {
	account := &BankAccount{balance: 1000}
	var wg sync.WaitGroup

	// 10 goroutines each deposit 100
	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			account.Deposit(100)
			fmt.Printf("Depositor %d: +100\n", id+1)
		}(i)
	}

	// 5 goroutines each withdraw 200
	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if account.Withdraw(200) {
				fmt.Printf("Withdrawer %d: -200\n", id+1)
			} else {
				fmt.Printf("Withdrawer %d: insufficient funds\n", id+1)
			}
		}(i)
	}

	wg.Wait()
	// Expected: 1000 + (10 * 100) - (5 * 200) = 1000
	fmt.Printf("Final balance: $%d (expected $1000)\n", account.Balance())
}

// demoSemaphore simulates rate-limited API calls.
// Many APIs limit concurrent requests (e.g., 3 at a time) to prevent overload.
func demoSemaphore() {
	maxConcurrent := 3
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	// Simulated user IDs to fetch from an "API"
	userIDs := []int{101, 102, 103, 104, 105, 106, 107, 108}
	results := make(chan string, len(userIDs))

	fmt.Printf("Fetching %d users (max %d concurrent requests)\n", len(userIDs), maxConcurrent)

	for _, userID := range userIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			sem <- struct{}{} // Acquire slot (blocks if 3 already running)

			// Simulate API call
			time.Sleep(80 * time.Millisecond)
			userData := fmt.Sprintf("User %d: name=user_%d, active=true", id, id)
			results <- userData

			fmt.Printf("  Fetched user %d\n", id)
			<-sem // Release slot
		}(userID)
	}

	wg.Wait()
	close(results)

	// Collect results
	fmt.Println("Results:")
	for result := range results {
		fmt.Printf("  %s\n", result)
	}
}
