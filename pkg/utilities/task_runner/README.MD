# Task Runner

Gorun is a package designed to facilitate the execution of tasks in parallel in Go. Using the power of goroutines and effective concurrency patterns, Gorun allows developers to implement and handle asynchronous tasks efficiently and safely.

### Features

- **Efficient Parallel Execution:** Leverages goroutines to execute tasks concurrently, maximizing the use of system resources.
- **Simplified Task Handling:** Defines tasks with inputs and outputs of any type, simplifying the execution of concurrent functions.
- **Support for Runtime and Results:** Automatically measures the runtime of each task and captures the results for easy management.
- **Goroutines Pooling:** Reduces the overhead of constantly creating and destroying goroutines through an optional pooling system.
- **Efficient Memory Use:** Optimizations to reduce unnecessary allocations and improve overall performance.

### Basic Implementation

Below is a basic example of how to define and execute tasks using Gorun:

```go
package main

import (
	"context"
	"errors"
	"fmt"
	tr "github.com/skolldire/web-simplify/pkg/task_runner/task_runner"
)

func main() {
	// Define your tasks
	taskOk := tr.Task[int, int]{
		Args: 1,
		Func: func(ctx context.Context, i int) (int, error) {
			return i + 1, nil
		},
	}

	taskError := tr.Task[int, int]{
		Args: 0,
		Func: func(ctx context.Context, i int) (int, error) {
			if i == 0 {
				return 0, errors.New("cannot divide by zero")
			}
			return 2 / i, nil
		},
	}

	// Execute all tasks
	results := tr.Execute(context.TODO(), map[string]tr.Tasker{
		"task-ok":         taskOk,
		"task-with-error": taskError,
	})

	// Show the results
	for id, result := range results {
		fmt.Printf("Result of %s: %v, duration: %d ms error: %v\n",
			id,
			result.Res,
			result.Time,
			result.Err,
		)
	}
}
```
