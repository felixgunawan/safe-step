# safe-step

![alt text](https://github.com/felixgunawan/safe-step/blob/master/img/img.jpg?raw=true)

A simple golang package to safely handle your multiple layers goroutine execution.

## Installation

```bash
go get github.com/felixgunawan/safe-step
```

## Example

```golang
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
- Recover from panic (will convert panic to error)
- Ability to set max timeout for all async functions execution

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
