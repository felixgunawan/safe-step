# safe-step
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://raw.githubusercontent.com/felixgunawan/safe-step/master/LICENSE)
[![GoReport](https://goreportcard.com/badge/github.com/felixgunawan/safe-step)](https://goreportcard.com/report/github.com/felixgunawan/safe-step)

A simple golang library to safely handle your multiple layers goroutine execution.
<p align="center">
  <img src="https://github.com/felixgunawan/safe-step/blob/master/img/img.jpg?raw=true">
</p>

## Installation

```bash
go get github.com/felixgunawan/safe-step
```

## Example 1

```golang
package main

import (
	"fmt"
	safestep "github.com/felixgunawan/safe-step"
	"time"
)

func main() {
	step := safestep.New()
	step.AddInput("id", 1)
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
```

## Example 2

```golang
package main

import (
	"fmt"
	safestep "github.com/felixgunawan/safe-step"
	"time"
)

func main() {
	step := safestep.New()
	step.AddInput("id", 1)
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
		AddFunction("f4", f4).
		AddFunction("f5", f5).
		DoWithMaxConcurrency(3)
	if err != nil {
		fmt.Printf("err = %v", err)
	}
	fmt.Printf("result = %v", res)
}
```

Output : 
```golang
function 1 started
id = 1
function 2 started
function 3 started
function 1 ended
function 2 ended
function 3 ended
function 5 started
id2 = 2
function 4 started
function 4 ended
function 5 ended
result = map[f1:1 f2:1.5 f3:3 f4:abcde f5:5]
```

## Features
- Safe goroutine, will recover from panic inside goroutine execution (will convert panic to error)
- Context-aware (see /example/timeout for implementation example using context.WithTimeout)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
