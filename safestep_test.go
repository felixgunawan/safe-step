package safestep

import (
	"context"
	"testing"
	"time"
)

func TestOperationBasic(t *testing.T) {
	step := New()
	step = step.AddInput("id", 1)
	f1 := func(input map[string]interface{}) (interface{}, error) {
		return 1, nil
	}
	f2 := func(input map[string]interface{}) (interface{}, error) {
		step.AddInput("id2", 2)
		return 1.5, nil
	}
	f3 := func(input map[string]interface{}) (interface{}, error) {
		return 3, nil
	}
	f4 := func(input map[string]interface{}) (interface{}, error) {
		return "abcde", nil
	}
	f5 := func(input map[string]interface{}) (interface{}, error) {
		return 5, nil
	}
	_, err := step.
		AddFunction("f1", f1).
		AddFunction("f2", f2).
		AddFunction("f3", f3).
		Step().
		AddFunction("f4", f4).
		AddFunction("f5", f5).
		Do()
	if err != nil {
		t.Errorf("error on basic operation : %v", err)
	}
}

func TestOperationPanic(t *testing.T) {
	step := New()
	step = step.AddInput("id", 1)
	f1 := func(input map[string]interface{}) (interface{}, error) {
		panic("test")
		return 1, nil
	}
	f2 := func(input map[string]interface{}) (interface{}, error) {
		step.AddInput("id2", 2)
		return 2, nil
	}
	f3 := func(input map[string]interface{}) (interface{}, error) {
		return 3, nil
	}
	f4 := func(input map[string]interface{}) (interface{}, error) {
		return 4, nil
	}
	f5 := func(input map[string]interface{}) (interface{}, error) {
		return 5, nil
	}
	_, err := step.
		AddFunction("f1", f1).
		AddFunction("f2", f2).
		AddFunction("f3", f3).
		Step().
		AddFunction("f4", f4).
		AddFunction("f5", f5).
		Do()
	if err == nil {
		t.Errorf("no error triggered on panic test")
	}
}

func TestOperationTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	step := NewWithContext(ctx)
	step = step.AddInput("id", 1)
	f1 := func(input map[string]interface{}) (interface{}, error) {
		time.Sleep(20 * time.Millisecond)
		return 1, nil
	}
	f2 := func(input map[string]interface{}) (interface{}, error) {
		step.AddInput("id2", 2)
		return 1.5, nil
	}
	f3 := func(input map[string]interface{}) (interface{}, error) {
		return 3, nil
	}
	f4 := func(input map[string]interface{}) (interface{}, error) {
		return "abcde", nil
	}
	f5 := func(input map[string]interface{}) (interface{}, error) {
		return 5, nil
	}
	_, err := step.
		AddFunction("f1", f1).
		AddFunction("f2", f2).
		AddFunction("f3", f3).
		Step().
		AddFunction("f4", f4).
		AddFunction("f5", f5).
		Do()
	if err != context.DeadlineExceeded {
		t.Errorf("no timeout error triggered")
	}
}

func TestOperationBasicWithMaxConcurrency(t *testing.T) {
	step := New()
	step = step.AddInput("id", 1)
	f1 := func(input map[string]interface{}) (interface{}, error) {
		return 1, nil
	}
	f2 := func(input map[string]interface{}) (interface{}, error) {
		step.AddInput("id2", 2)
		return 1.5, nil
	}
	f3 := func(input map[string]interface{}) (interface{}, error) {
		return 3, nil
	}
	f4 := func(input map[string]interface{}) (interface{}, error) {
		return "abcde", nil
	}
	f5 := func(input map[string]interface{}) (interface{}, error) {
		return 5, nil
	}
	_, err := step.
		AddFunction("f1", f1).
		AddFunction("f2", f2).
		AddFunction("f3", f3).
		Step().
		AddFunction("f4", f4).
		AddFunction("f5", f5).
		DoWithMaxConcurrency(2)
	if err != nil {
		t.Errorf("error on basic operation : %v", err)
	}
}

func TestOperationPanicWithMaxConcurrency(t *testing.T) {
	step := New()
	step = step.AddInput("id", 1)
	f1 := func(input map[string]interface{}) (interface{}, error) {
		panic("test")
		return 1, nil
	}
	f2 := func(input map[string]interface{}) (interface{}, error) {
		step.AddInput("id2", 2)
		return 2, nil
	}
	f3 := func(input map[string]interface{}) (interface{}, error) {
		return 3, nil
	}
	f4 := func(input map[string]interface{}) (interface{}, error) {
		return 4, nil
	}
	f5 := func(input map[string]interface{}) (interface{}, error) {
		return 5, nil
	}
	_, err := step.
		AddFunction("f1", f1).
		AddFunction("f2", f2).
		AddFunction("f3", f3).
		Step().
		AddFunction("f4", f4).
		AddFunction("f5", f5).
		DoWithMaxConcurrency(3)
	if err == nil {
		t.Errorf("no error triggered on panic test")
	}
}

func TestOperationTimeoutWithMaxConcurrency(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	step := NewWithContext(ctx)
	step = step.AddInput("id", 1)
	f1 := func(input map[string]interface{}) (interface{}, error) {
		time.Sleep(20 * time.Millisecond)
		return 1, nil
	}
	f2 := func(input map[string]interface{}) (interface{}, error) {
		step.AddInput("id2", 2)
		return 1.5, nil
	}
	f3 := func(input map[string]interface{}) (interface{}, error) {
		return 3, nil
	}
	f4 := func(input map[string]interface{}) (interface{}, error) {
		return "abcde", nil
	}
	f5 := func(input map[string]interface{}) (interface{}, error) {
		return 5, nil
	}
	_, err := step.
		AddFunction("f1", f1).
		AddFunction("f2", f2).
		AddFunction("f3", f3).
		Step().
		AddFunction("f4", f4).
		AddFunction("f5", f5).
		DoWithMaxConcurrency(5)
	if err != context.DeadlineExceeded {
		t.Errorf("no timeout error triggered %v", err)
	}
}
