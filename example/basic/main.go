package main

import (
	"fmt"
	safestep "github.com/felixgunawan/safe-step"
	"time"
)

func main() {
	step := safestep.New()
	step = step.AddInput("id", 1)
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
	// this will :
	// 1. run f1,f2,f3 in goroutine and wait for all of them to finish
	// 2. run f4,f5 in goroutine and wait again
	// 3. return result of all function execution in map
	res, err := step.
		AddFunction("f1", f1).
		AddFunction("f2", f2).
		AddFunction("f3", f3).
		Step().
		AddFunction("f4", f4).
		AddFunction("f5", f5).
		Do()
	if err != nil {
		fmt.Printf("err = %v", err)
	}
	fmt.Printf("result = %v", res)
}
