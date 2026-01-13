package main

// Task represents a unit of work to be processed
type Task struct {
	ID   int
	Data string
}

// Result contains the processed output for a task
type Result struct {
	ID   int
	Hash string
}
