# Goroutines Demo

A simple demonstration of how Go's goroutines can speed up CPU-intensive operations by utilizing multiple CPU cores.

## Overview

This project compares sequential vs concurrent processing of tasks using goroutines. It demonstrates the worker pool pattern for parallel task processing.

## Quick Start

```bash
go run .
```

## Output

```
Goroutines Demo - Processing 100 tasks
CPU Cores: 16

Sequential:  8.502ms  (11762 items/sec)
2 workers:   4.498ms  (22230 items/sec) - 1.89x faster
4 workers:   2.498ms  (40019 items/sec) - 3.40x faster
16 workers:  1.002ms  (99800 items/sec) - 8.49x faster

--- Mutex Demo ---
Depositor 1: +100
Withdrawer 2: -200
...
Final balance: $1000 (expected $1000)

--- Semaphore Demo ---
Fetching 8 users (max 3 concurrent requests)
  Fetched user 101
  Fetched user 102
  ...
Results:
  User 101: name=user_101, active=true
  ...
```

## Project Structure

```
.
├── main.go      # Main processing logic and benchmark
├── structs.go   # Task and Result data structures
├── go.mod       # Go module definition
└── README.md    # Project documentation
```

## How It Works

### Sequential Processing
Processes tasks one at a time in order. Simple but only uses one CPU core.

### Concurrent Processing
Uses a worker pool pattern where multiple goroutines process tasks in parallel from a shared job queue.

Key components:
- **Channels**: For sending tasks to workers and collecting results
- **WaitGroups**: To wait for all workers to complete
- **Goroutines**: Lightweight threads that process tasks concurrently

### Mutex (Mutual Exclusion)
A mutex protects shared state so only one goroutine can access it at a time. Without it, concurrent writes cause race conditions and lost updates.

The demo simulates a bank account with concurrent deposits and withdrawals. The mutex ensures the balance is always correct.

### Semaphore
A semaphore limits how many goroutines can perform an action simultaneously. In Go, a buffered channel acts as a semaphore.

The demo simulates rate-limited API calls, fetching user data with a max of 3 concurrent requests (common API restriction).

## Requirements

- Go 1.22+ (uses range over integers)
- No external dependencies
