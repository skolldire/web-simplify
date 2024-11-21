package task_runner

import (
	"context"
	"errors"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

// TestService_Execute
func TestService_Execute_NoError(t *testing.T) {
	// Arrange
	tasks := map[string]Tasker{
		"task1": Task[int, int]{
			Func: func(ctx context.Context, i int) (int, error) {
				return i + 1, nil
			},
			Args: 1,
		},
	}

	// Act
	res := Execute(context.TODO(), tasks)

	// Assert
	assert.Equal(t, res["task1"].Res, 2)
}

// TestService_Execute_Error
func TestService_Execute_Error(t *testing.T) {
	// Arrange
	tasks := map[string]Tasker{
		"task1": Task[int, int]{
			Func: func(ctx context.Context, i int) (int, error) {
				return 0, errors.New("error")
			},
			Args: 1,
		},
	}

	// Act
	res := Execute(context.TODO(), tasks)

	// Assert
	assert.Equal(t, res["task1"].Err, errors.New("error"))
}

// TestService_Execute_MultipleTasks
func TestService_Execute_MultipleTasks(t *testing.T) {
	// Arrange
	tasks := map[string]Tasker{
		"task1": Task[int, int]{
			Func: func(ctx context.Context, i int) (int, error) {
				return i + 1, nil
			},
			Args: 1,
		},
		"task2": Task[int, int]{
			Func: func(ctx context.Context, i int) (int, error) {
				return i + 2, nil
			},
			Args: 1,
		},
	}

	// Act
	res := Execute(context.TODO(), tasks)

	// Assert
	assert.Equal(t, res["task1"].Res, 2)
	assert.Equal(t, res["task2"].Res, 3)
}

// TestService_Execute_MultipleTasksWithError
func TestService_Execute_MultipleTasksWithError(t *testing.T) {
	// Arrange
	tasks := map[string]Tasker{
		"task1": Task[int, int]{
			Func: func(ctx context.Context, i int) (int, error) {
				return i + 1, nil
			},
			Args: 1,
		},
		"task2": Task[int, int]{
			Func: func(ctx context.Context, i int) (int, error) {
				return 0, errors.New("error")
			},
			Args: 1,
		},
	}

	// Act
	res := Execute(context.TODO(), tasks)

	// Assert
	assert.Equal(t, res["task1"].Res, 2)
	assert.Equal(t, res["task2"].Err, errors.New("error"))
}

// TestService_Execute_MultipleTasksWithMultipleErrors
func TestService_Execute_MultipleTasksWithMultipleErrors(t *testing.T) {
	// Arrange
	tasks := map[string]Tasker{
		"task1": Task[int, int]{
			Func: func(ctx context.Context, i int) (int, error) {
				return 0, errors.New("error")
			},
			Args: 1,
		},
		"task2": Task[int, int]{
			Func: func(ctx context.Context, i int) (int, error) {
				return 0, errors.New("error")
			},
			Args: 1,
		},
	}

	// Act
	res := Execute(context.TODO(), tasks)

	// Assert
	assert.Equal(t, res["task1"].Err, errors.New("error"))
	assert.Equal(t, res["task2"].Err, errors.New("error"))
}

// TestService_Execute benchmark
func BenchmarkService_Execute(b *testing.B) {
	// Arrange
	tasks := map[string]Tasker{
		"task1": Task[int, int]{
			Func: func(ctx context.Context, i int) (int, error) {
				return i + 1, nil
			},
			Args: 1,
		},
		"task2": Task[int, int]{
			Func: func(ctx context.Context, i int) (int, error) {
				return i + 2, nil
			},
			Args: 1,
		},
	}

	// Act
	for i := 0; i < b.N; i++ {
		Execute(context.TODO(), tasks)
	}
}
