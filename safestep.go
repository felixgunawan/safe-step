package safestep

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// SafeStep used for handling multiple layers goroutine execution
type SafeStep struct {
	ctx       context.Context
	cancel    context.CancelFunc
	input     map[string]interface{}
	step      []map[string]asyncFunc
	tempFuncs map[string]asyncFunc
	result    map[string]interface{}
	err       error
}

type goRoutineResp struct {
	code   string
	result interface{}
	err    error
}

type asyncFunc func(input map[string]interface{}) (interface{}, error)

// ErrTimeout timeout error
var ErrTimeout = errors.New("timeout occurred while executing async function")

// New initialization
func New() *SafeStep {
	return &SafeStep{
		input:     make(map[string]interface{}),
		result:    make(map[string]interface{}),
		tempFuncs: make(map[string]asyncFunc),
		step:      make([]map[string]asyncFunc, 0),
	}
}

// AddCtx adding context to safestep struct
func (step *SafeStep) AddCtx(ctx context.Context) *SafeStep {
	step.ctx = ctx
	return step
}

// SetTimeout set max timeout for all steps to finish (context timeout)
func (step *SafeStep) SetTimeout(timeout time.Duration) *SafeStep {
	if step.ctx == nil {
		step.ctx = context.Background()
	}
	step.ctx, step.cancel = context.WithTimeout(step.ctx, timeout)
	return step
}

// AddInput add input which can be used by asyncFunc parameter, note that it can also be used in previous step so it can acts like dependency
func (step *SafeStep) AddInput(code string, input interface{}) *SafeStep {
	step.input[code] = input
	return step
}

// AddFunction adding asyncFunc with function code (must be unique, otherwise previous function result will be overwritten)
func (step *SafeStep) AddFunction(code string, function asyncFunc) *SafeStep {
	step.tempFuncs[code] = function
	return step
}

// Step appends async func step
func (step *SafeStep) Step() *SafeStep {
	step.step = append(step.step, step.tempFuncs)
	step.tempFuncs = make(map[string]asyncFunc)
	return step
}

// Do execute all async functions according to their order
func (step *SafeStep) Do() (map[string]interface{}, error) {
	// check just in case there are still some function not appended
	if len(step.tempFuncs) > 0 {
		step.Step()
	}
	// execute async funcs in their respective order
	for _, s := range step.step {
		chGo := make(chan goRoutineResp, len(s)) // to get asyncFunc result
		var wg sync.WaitGroup                    // to wait all goroutine/timeout finish (whichever first)
		var mu sync.RWMutex                      // prevent race condition
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
		if waitTimeout(step.ctx, &wg) { // check timeout
			step.err = ErrTimeout
			break
		}
		for range s { // put all asyncFunc response to result based on function code
			mu.RLock()
			resp := <-chGo
			mu.RUnlock()
			if resp.err != nil {
				step.err = resp.err
				break
			}
			step.result[resp.code] = resp.result
		}
		close(chGo)
	}
	return step.result, step.err
}

// waitTimeout waits for the waitgroup for the specified max timeout, return true if timeout occurs
func waitTimeout(ctx context.Context, wg *sync.WaitGroup) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-ctx.Done():
		return true // timed out
	}
}
