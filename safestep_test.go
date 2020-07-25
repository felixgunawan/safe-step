package safestep

import (
	"fmt"
	"testing"
	"time"
)

func TestOperationBasic(t *testing.T) {
	step := New()
	step = step.AddInput("id", 1)
	f1 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function1 started")
		fmt.Printf("id = %d\n", input["id"])
		time.Sleep(time.Millisecond * 1000)
		fmt.Println("function1 ended")
		return 1, nil
	}
	f2 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function2 started")
		step.AddInput("id2", 2)
		return 2, nil
	}
	f3 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function3 started")
		return 3, nil
	}
	f4 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function4 started")
		return 4, nil
	}
	f5 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function5 started")
		fmt.Printf("id2 = %d\n", input["id2"])
		return 5, nil
	}
	res, err := step.
		AddFunction("f1", f1).
		AddFunction("f2", f2).
		AddFunction("f3", f3).
		Step().
		AddFunction("f4", f4).
		AddFunction("f5", f5).
		Do()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("result = %v", res)
}

func TestOperationPanic(t *testing.T) {
	step := New()
	step = step.AddInput("id", 1)
	f1 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function1 started")
		fmt.Printf("id = %d\n", input["id"])
		time.Sleep(time.Millisecond * 1000)
		panic("test")
		return 1, nil
	}
	f2 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function2 started")
		step.AddInput("id2", 2)
		return 2, nil
	}
	f3 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function3 started")
		return 3, nil
	}
	f4 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function4 started")
		return 4, nil
	}
	f5 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function5 started")
		fmt.Printf("id2 = %d\n", input["id2"])
		return 5, nil
	}
	res, err := step.
		AddFunction("f1", f1).
		AddFunction("f2", f2).
		AddFunction("f3", f3).
		Step().
		AddFunction("f4", f4).
		AddFunction("f5", f5).
		Do()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("result = %v", res)
}
