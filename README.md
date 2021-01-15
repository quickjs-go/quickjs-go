# quickjs

[![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](LICENSE)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/quickjs-go/quickjs-go)

Go bindings to [QuickJS](https://bellard.org/quickjs/): a fast, small, and embeddable [ES2020](https://tc39.github.io/ecma262/) JavaScript interpreter.

These bindings are a WIP and do not match full parity with QuickJS' API, though expose just enough features to be usable.

These bindings have been tested to cross-compile and run successfully on Linux, Windows, and Mac using gcc-7 and mingw32 without any addtional compiler or linker flags.

## Usage

```
$ go get github.com/quickjs-go/quickjs-go
```

## Guidelines

1. Free `quickjs.Runtime` and `quickjs.Context` once you are done using them.
2. Free `quickjs.Value`'s returned by `Eval()` and `EvalFile()`. All other values do not need to be freed, as they get garbage-collected.
3. You may access the stacktrace of an error returned by `Eval()` or `EvalFile()` or `Call()` by casting it to a `*quickjs.Error`.
4. Make new copies of arguments should you want to return them in functions you created.
5. Make sure to call `runtime.LockOSThread()` to ensure that QuickJS always operates in the exact same thread.
6. Add JsInterface and JsThread for run javascript in golang goroutine

## Example

The full example code below may be found by clicking [here](examples/normal/main.go). Find more API examples [here](quickjs_test.go) and advance use found by clicking [here](examples/advance/main.go).

```go
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/quickjs-go/quickjs-go"
	"strings"
)

func check(err error) {
	if err != nil {
		var evalErr *quickjs.Error
		if errors.As(err, &evalErr) {
			fmt.Println(evalErr.Cause)
			fmt.Println(evalErr.Stack)
		}
		panic(err)
	}
}

func main() {
	runtime := quickjs.NewRuntime()
	defer runtime.Free()

	context := runtime.NewContext()
	defer context.Free()

	globals := context.Globals()

	// Test evaluating template strings.

	result, err := context.Eval("`Hello world! 2 ** 8 = ${2 ** 8}.`")
	check(err)
	defer result.Free()

	fmt.Println(result.String())
	fmt.Println()

	// Test evaluating numeric expressions.

	result, err = context.Eval(`1 + 2 * 100 - 3 + Math.sin(10)`)
	check(err)
	defer result.Free()

	fmt.Println(result.Int64())
	fmt.Println()

	// Test evaluating big integer expressions.

	result, err = context.Eval(`128n ** 16n`)
	check(err)
	defer result.Free()

	fmt.Println(result.BigInt())
	fmt.Println()

	// Test evaluating big decimal expressions.

	result, err = context.Eval(`128l ** 12l`)
	check(err)
	defer result.Free()

	fmt.Println(result.BigFloat())
	fmt.Println()

	// Test evaluating boolean expressions.

	result, err = context.Eval(`false && true`)
	check(err)
	defer result.Free()

	fmt.Println(result.Bool())
	fmt.Println()

	// Test setting and calling functions.

	A := func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		fmt.Println("A got called!")
		return ctx.Null()
	}

	B := func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		fmt.Println("B got called!")
		return ctx.Null()
	}

	globals.Set("A", context.Function(A))
	globals.Set("B", context.Function(B))

	_, err = context.Eval(`for (let i = 0; i < 10; i++) { if (i % 2 === 0) A(); else B(); }`)
	check(err)

	fmt.Println()

	// Test setting global variables.

	_, err = context.Eval(`HELLO = "world"; TEST = false;`)
	check(err)

	names, err := globals.PropertyNames()
	check(err)

	fmt.Println("Globals:")
	for _, name := range names {
		val := globals.GetByAtom(name.Atom)
		defer val.Free()

		fmt.Printf("'%s': %s\n", name, val)
	}
	fmt.Println()

	// Test evaluating arbitrary expressions from flag arguments.

	flag.Parse()
	if flag.NArg() == 0 {
		return
	}

	result, err = context.Eval(strings.Join(flag.Args(), " "))
	check(err)
	defer result.Free()

	if result.IsObject() {
		names, err := result.PropertyNames()
		check(err)

		fmt.Println("Object:")
		for _, name := range names {
			val := result.GetByAtom(name.Atom)
			defer val.Free()

			fmt.Printf("'%s': %s\n", name, val)
		}
	} else {
		fmt.Println(result.String())
	}
}
```

## License

QuickJS is released under the MIT license.
