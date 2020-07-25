package safestep

import (
	"fmt"
	"sync"
)

type SafeStep struct {
	input     map[string]interface{}
	step      []map[string]asyncFunc
	result    map[string]interface{}
	error     error
	tempFuncs map[string]asyncFunc
}

type goRoutineResp struct {
	code   string
	result interface{}
	err    error
}

type asyncFunc func(input map[string]interface{}) (interface{}, error)

func New() *SafeStep {
	return &SafeStep{
		input:     make(map[string]interface{}),
		result:    make(map[string]interface{}),
		tempFuncs: make(map[string]asyncFunc),
		step:      make([]map[string]asyncFunc, 0),
	}
}

func (step *SafeStep) AddInput(code string, input interface{}) *SafeStep {
	step.input[code] = input
	return step
}

func (step *SafeStep) AddFunction(code string, function asyncFunc) *SafeStep {
	step.tempFuncs[code] = function
	return step
}

func (step *SafeStep) Step() *SafeStep {
	step.step = append(step.step, step.tempFuncs)
	step.tempFuncs = make(map[string]asyncFunc)
	return step
}

func (step *SafeStep) Do() (map[string]interface{}, error) {
	if len(step.tempFuncs) > 0 {
		step.Step()
	}
	for _, s := range step.step {
		chGo := make(chan goRoutineResp, len(s))
		var wg sync.WaitGroup
		var mu sync.RWMutex
		wg.Add(len(s))
		for code, f := range s {
			go func(code string, f asyncFunc) {
				defer func() {
					if r := recover(); r != nil {
						mu.Lock()
						chGo <- goRoutineResp{
							code: code,
							err:  fmt.Errorf("%v", r),
						}
						mu.Unlock()
						wg.Done()
					}
				}()
				res, err := f(step.input)
				mu.Lock()
				chGo <- goRoutineResp{
					code:   code,
					result: res,
					err:    err,
				}
				mu.Unlock()
				wg.Done()
			}(code, f)
		}
		wg.Wait()
		for range s {
			mu.RLock()
			resp := <-chGo
			mu.RUnlock()
			if resp.err != nil {
				return step.result, resp.err
			}
			step.result[resp.code] = resp.result
		}
		close(chGo)
	}
	return step.result, nil
}
