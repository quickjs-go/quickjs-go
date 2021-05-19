package main

import (
	"errors"
	"log"

	"github.com/quickjs-go/quickjs-go"
)

type callback func(string) string

type Instance struct {
	callbacks map[string]callback
}

func (i *Instance) Add(str string, fn callback) {
	i.callbacks[str] = fn
}

func (i *Instance) Run(str string, arg string) string {
	if fn, ok := i.callbacks[str]; ok {
		return fn(arg)
	} else {
		return ""
	}
}

type Demo struct {
	Instance
}

func NewDemo() *Demo {
	return &Demo{}
}

func (d *Demo) Register(runtime quickjs.Runtime, context *quickjs.Context, thread *quickjs.JsThread) bool {
	test := context.Object()
	test.SetFunction("AddCallback", func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		log.Println(args)
		if len(args) != 2 {
			return ctx.Error(errors.New("args must be 2"))
		}

		if args[0].IsString() && args[1].IsFunction() {
			var _instance *Instance
			value := test.Get("instance")
			if value.IsNumber() {
				if temp, ok := quickjs.ObjectId(value.Int64()).Get(); ok {
					if _instance, ok = temp.(*Instance); !ok {
						return ctx.Error(errors.New("instance type error"))
					}
				} else {
					return ctx.Error(errors.New("instance not exist"))
				}
			} else {
				return ctx.Error(errors.New("instance error"))
			}

			if _instance == nil {
				return ctx.Error(errors.New("error:instance is nil"))
			}

			fn := ctx.DupValue(args[1])
			_instance.Add(args[0].String(), func(str string) string {
				var result string
				if fn.IsFunction() {
					r, err := thread.Call(context.Null(), fn, []quickjs.Value{context.String(str)})
					log.Println(r, err)
					result = r.String()
				} else {
					log.Println("fn is not function")
				}

				return result
			})

			return ctx.Bool(true)
		} else {
			return ctx.Error(errors.New("args type error"))
		}
	})

	_instance := &Instance{callbacks: make(map[string]callback)}
	d.Instance = *_instance
	obj := quickjs.NewObjectId(_instance)
	test.Set("instance", context.Int64(int64(obj)))
	context.Globals().Set("test", test)

	return true
}

func (d *Demo) Unregister(runtime quickjs.Runtime, context *quickjs.Context, thread *quickjs.JsThread) {

}

func main() {
	d := &Demo{}
	js := quickjs.NewJsThread(d)
	if js != nil {
		defer js.Close()
	}

	result, err := js.Eval(`
		import * as os from 'os'
		console.log("run 2000")
		test.AddCallback("key",function(info) {
			os.sleep(2000)
			return "in js callback:" + info
		})
	`, quickjs.EVAL_MODULE)

	log.Println(result, err)

	d.Run("key", "i am in golang")

	go func() {
		d.Run("key", "i am in golang goroutine")
	}()

	select {}
}
