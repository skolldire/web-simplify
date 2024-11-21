package task_runner

import "context"

// Tasker defines an interface for an executable task.
type Tasker interface {
	Execute(ctx context.Context) (result interface{}, duration int, err error)
}

// Task implements the Tasker interface for tasks with input and output of any type.
type Task[I, O any] struct {
	Func func(context.Context, I) (O, error)
	Args I
}

// Result represents the result of an executed task.
type Result struct {
	ID   string
	Err  error
	Res  interface{}
	Time int
}
