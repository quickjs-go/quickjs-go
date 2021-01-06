package quickjs

import (
	osruntime "runtime"
)

type JsInterface interface {
	Register(runtime Runtime, context *Context) bool
	Unregister(runtime Runtime, context *Context)
}

type JsThread struct {
	*Runtime
	*Context
	Interface JsInterface
	eval      chan jsEval
	call      chan jsCall
	close     chan struct{}
}

type jsResult struct {
	Result Value
	Error  error
}

type jsEval struct {
	Code   string
	Result chan jsResult
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
		context.InitStdModule()
		context.InitOsModule()
		context.StdHelper()

		eval := make(chan jsEval)
		call := make(chan jsCall)
		close := make(chan struct{})
		jsthread <- &JsThread{&runtime, context, Interface, eval, call, close}

		if Interface.Register(runtime, context) {
			defer Interface.Unregister(runtime, context)
		} else {
			return
		}

		for {
			select {
			case _eval, ok := <-eval:
				if ok {
					_result, _error := context.Eval(_eval.Code)
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

	return <-jsthread
}

func (j *JsThread) Close() {
	close(j.close)

	if j.Context != nil {
		j.Context.Free()
		j.Context = nil
	}

	if j.Runtime != nil {
		j.Runtime.Free()
		j.Runtime = nil
	}
}

func (j *JsThread) Eval(code string) (result Value, err error) {
	temp := jsEval{code, make(chan jsResult)}
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
