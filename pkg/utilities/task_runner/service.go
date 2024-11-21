package task_runner

import (
	"context"
	"sync"
	"time"
)

// Execute all tasks in parallel and collect the results.
func Execute(ctx context.Context, tasks map[string]Tasker) map[string]Result {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	results := make(map[string]Result)

	for id, task := range tasks {
		wg.Add(1)
		go func(id string, task Tasker) {
			defer wg.Done()
			res, t, err := task.Execute(ctx)
			mutex.Lock()
			results[id] = Result{ID: id, Err: err, Res: res, Time: t}
			mutex.Unlock()
		}(id, task)
	}
	wg.Wait()

	return results
}

// Execute the task and returns the result and execution time.
func (t Task[I, O]) Execute(ctx context.Context) (interface{}, int, error) {
	start := time.Now()
	out, err := t.Func(ctx, t.Args)
	duration := time.Since(start)
	return out, int(duration.Milliseconds()), err
}
