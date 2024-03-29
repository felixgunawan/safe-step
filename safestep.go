package safestep

import (
	"context"
	"fmt"
	"sync"
)

// SafeStep used for handling multiple layers goroutine execution
type SafeStep interface {
	// AddInput add input which can be used by asyncFunc parameter
	// it can also be called in previous function call step so it can acts like dependency
	AddInput(code string, input interface{}) SafeStep
	// GetInput will get input thar can be used inside asyncFunc
	GetInput(code string) interface{}
	// AddFunction adding asyncFunc with function code
	// code must be unique, otherwise previous function result will be overwritten
	AddFunction(code string, function asyncFunc) SafeStep
	// Step appends async function step
	Step() SafeStep
	// Do execute all async functions according to their order
	Do() (map[string]interface{}, error)
	// Do execute all async functions according to their order(with maximum concurrency)
	DoWithMaxConcurrency(maxConcurrency int) (map[string]interface{}, error)
}

// SafeStepStruct used for handling multiple layers goroutine execution
type SafeStepStruct struct {
	ctx       context.Context
	input     map[string]interface{}
	step      []map[string]asyncFunc
	tempFuncs map[string]asyncFunc
	result    map[string]interface{}
	mapLock   sync.RWMutex
	err       error
}

type goRoutineResp struct {
	code   string
	result interface{}
	err    error
}

type asyncFunc func() (interface{}, error)

// New initialization
func New() SafeStep {
	return &SafeStepStruct{
		ctx:       context.Background(),
		input:     make(map[string]interface{}),
		result:    make(map[string]interface{}),
		tempFuncs: make(map[string]asyncFunc),
		step:      make([]map[string]asyncFunc, 0),
	}
}

// New initialization
func NewWithContext(ctx context.Context) SafeStep {
	return &SafeStepStruct{
		ctx:       ctx,
		input:     make(map[string]interface{}),
		result:    make(map[string]interface{}),
		tempFuncs: make(map[string]asyncFunc),
		step:      make([]map[string]asyncFunc, 0),
	}
}

// GetInput will get input thar can be used inside asyncFunc
func (step *SafeStepStruct) GetInput(code string) interface{} {
	step.mapLock.Lock()
	input := step.input[code]
	step.mapLock.Unlock()
	return input
}

// AddInput add input which can be used by asyncFunc parameter, note that it can also be used in previous step so it can acts like dependency
func (step *SafeStepStruct) AddInput(code string, input interface{}) SafeStep {
	step.mapLock.Lock()
	step.input[code] = input
	step.mapLock.Unlock()
	return step
}

// AddFunction adding asyncFunc with function code (must be unique, otherwise previous function result will be overwritten)
func (step *SafeStepStruct) AddFunction(code string, function asyncFunc) SafeStep {
	step.mapLock.Lock()
	step.tempFuncs[code] = function
	step.mapLock.Unlock()
	return step
}

// Step appends async function step
func (step *SafeStepStruct) Step() SafeStep {
	step.step = append(step.step, step.tempFuncs)
	step.tempFuncs = make(map[string]asyncFunc)
	return step
}

// Do execute all async functions according to their order
func (step *SafeStepStruct) Do() (map[string]interface{}, error) {
	// check just in case there are still some function not appended
	if len(step.tempFuncs) > 0 {
		step.Step()
	}
	// execute async funcs in their respective order
	for _, s := range step.step {
		chGo := make(chan goRoutineResp, len(s)) // to get asyncFunc result
		defer close(chGo)
		var wg sync.WaitGroup // to wait all goroutine/timeout finish (whichever first)
		var mu sync.RWMutex   // prevent race condition
		wg.Add(len(s))
		for code, f := range s {
			go func(code string, f asyncFunc) {
				defer func() { // recover go routine in case of panic
					if r := recover(); r != nil {
						mu.Lock()
						chGo <- goRoutineResp{
							code: code,
							err:  fmt.Errorf("%v", r), // convert panic to err
						}
						mu.Unlock()
						wg.Done()
					}
				}()
				res, err := f()
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
		step.err = waitTimeout(step.ctx, &wg)
		if step.err != nil {
			return step.result, step.err
		}
		for range s { // put all asyncFunc response to result based on function code
			mu.RLock()
			resp := <-chGo
			mu.RUnlock()
			if resp.err != nil {
				step.err = resp.err
				return step.result, step.err
			}
			step.result[resp.code] = resp.result
		}
	}
	return step.result, step.err
}

// Do execute all async functions according to their order(with maximum concurrency)
func (step *SafeStepStruct) DoWithMaxConcurrency(maxConcurrency int) (map[string]interface{}, error) {
	// check just in case there are still some function not appended
	if len(step.tempFuncs) > 0 {
		step.Step()
	}
	// execute async funcs in their respective order
	for _, s := range step.step {
		chGo := make(chan goRoutineResp, len(s)) // to get asyncFunc result
		defer close(chGo)
		var wg sync.WaitGroup // to wait all goroutine/timeout finish (whichever first)
		var mu sync.RWMutex   // prevent race condition
		var counter = 0       // count total data has been executed
		for code, f := range s {
			wg.Add(1)
			go func(code string, f asyncFunc) {
				defer func() { // recover go routine in case of panic
					if r := recover(); r != nil {
						mu.Lock()
						chGo <- goRoutineResp{
							code: code,
							err:  fmt.Errorf("%v", r), // convert panic to err
						}
						mu.Unlock()
						wg.Done()
					}
				}()
				res, err := f()
				mu.Lock()
				chGo <- goRoutineResp{
					code:   code,
					result: res,
					err:    err,
				}
				mu.Unlock()
				wg.Done()
			}(code, f)

			// limit maximum concurrency that will be executed in this function
			if ((counter + 1) == len(s)) ||
				((counter+1)%maxConcurrency == 0) {
				step.err = waitTimeout(step.ctx, &wg)
				if step.err != nil {
					return step.result, step.err
				}
			}
			counter++ // add counter
		}
		for range s { // put all asyncFunc response to result based on function code
			mu.RLock()
			resp := <-chGo
			mu.RUnlock()
			if resp.err != nil {
				step.err = resp.err
				return step.result, step.err
			}
			step.result[resp.code] = resp.result
		}
	}
	return step.result, step.err
}

// waitTimeout waits for the waitgroup for the specified max timeout, return error if context done
func waitTimeout(ctx context.Context, wg *sync.WaitGroup) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
