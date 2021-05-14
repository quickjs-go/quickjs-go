package quickjs

import (
	osruntime "runtime"
)

type JsInterface interface {
	Register(runtime Runtime, context *Context, thread *JsThread) bool
	Unregister(runtime Runtime, context *Context, thread *JsThread)
}

type JsThread struct {
	*Runtime
	*Context
	init  chan JsInterface
	eval  chan jsEval
	call  chan jsCall
	close chan struct{}
}

type jsResult struct {
	Result Value
	Error  error
}

type jsEval struct {
	Code   string
	Result chan jsResult
	EvalType int
}

type jsCall struct {
	obj    Value
	fn     Value
	args   []Value
	Result chan jsResult
}

func NewJsThread(Interface JsInterface) *JsThread {
	if Interface == nil {
		return nil
	}

	jsthread := make(chan *JsThread)

	go func() {
		osruntime.LockOSThread()

		runtime := NewRuntime()
		context := runtime.NewContext()

		init := make(chan JsInterface)
		eval := make(chan jsEval)
		call := make(chan jsCall)
		close := make(chan struct{})
		thread := &JsThread{&runtime, context, init, eval, call, close}
		jsthread <- thread

		for {
			select {
			case _init, ok := <-init:
				if ok {
					if _init.Register(runtime, context, thread) {
						defer _init.Unregister(runtime, context, thread)
					}
				}
			case _eval, ok := <-eval:
				if ok {
					_result, _error := context.Eval(_eval.Code,_eval.EvalType)
					_eval.Result <- jsResult{_result, _error}
				}
			case _call, ok := <-call:
				if ok {
					_result, _error := context.Call(_call.obj, _call.fn, _call.args)
					_call.Result <- jsResult{_result, _error}
				}
			case <-close:
				break
			}
		}
	}()

	temp := <-jsthread
	temp.init <- Interface
	return temp
}

func (j *JsThread) Close() {
	close(j.close)

	if j.Runtime != nil {
		j.Runtime.StdFreeHandlers()
	}

	if j.Context != nil {
		j.Context.Free()
		j.Context = nil
	}

	if j.Runtime != nil {
		j.Runtime.Free()
		j.Runtime = nil
	}
}

func (j *JsThread) Eval(code string,evaltype int) (result Value, err error) {
	temp := jsEval{code, make(chan jsResult),evaltype}
	j.eval <- temp

	select {
	case m, ok := <-temp.Result:
		if ok {
			result = m.Result
			err = m.Error
		}
	}

	return
}

func (j *JsThread) Call(obj, fn Value, args []Value) (result Value, err error) {
	temp := jsCall{obj, fn, args, make(chan jsResult)}
	j.call <- temp

	select {
	case m, ok := <-temp.Result:
		if ok {
			result = m.Result
			err = m.Error
		}
	}

	return
}
