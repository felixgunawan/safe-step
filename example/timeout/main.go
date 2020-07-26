package main

import (
	"context"
	"fmt"
	safestep "github.com/felixgunawan/safe-step"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	step := NewWithContext(ctx)
	step = safestep.AddInput("id", 1)
	f1 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function 1 started")
		fmt.Printf("id = %d\n", input["id"])
		time.Sleep(time.Millisecond * 500)
		fmt.Println("function 1 ended")
		return 1, nil
	}
	f2 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function 2 started")
		time.Sleep(time.Millisecond * 750)
		fmt.Println("function 2 ended")
		step.AddInput("id2", 2)
		return 1.5, nil
	}
	f3 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function 3 started")
		time.Sleep(time.Millisecond * 1000)
		fmt.Println("function 3 ended")
		return 3, nil
	}
	f4 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function 4 started")
		time.Sleep(time.Millisecond * 100)
		fmt.Println("function 4 ended")
		return "abcde", nil
	}
	f5 := func(input map[string]interface{}) (interface{}, error) {
		fmt.Println("function 5 started")
		fmt.Printf("id2 = %d\n", input["id2"])
		time.Sleep(time.Millisecond * 1000)
		fmt.Println("function 5 ended")
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
		fmt.Printf("err = %v\n", err)
	}
	fmt.Printf("result = %v\n", res)
}
