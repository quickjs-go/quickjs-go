package quickjs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObject(t *testing.T) {
	runtime := NewRuntime()
	defer runtime.Free()

	context := runtime.NewContext()
	defer context.Free()

	test := context.Object()
	test.Set("A", context.String("String A"))
	test.Set("B", context.String("String B"))
	test.Set("C", context.String("String C"))
	context.Globals().Set("test", test)

	result, err := context.Eval(`Object.keys(test).map(key => test[key]).join(" ")`, EVAL_GLOBAL)
	require.NoError(t, err)
	defer result.Free()

	require.EqualValues(t, "String A String B String C", result.String())
}

func TestArray(t *testing.T) {
	runtime := NewRuntime()
	defer runtime.Free()

	context := runtime.NewContext()
	defer context.Free()

	test := context.Array()
	for i := int64(0); i < 3; i++ {
		test.SetByInt64(i, context.String(fmt.Sprintf("test %d", i)))
	}
	for i := int64(0); i < test.Len(); i++ {
		require.EqualValues(t, fmt.Sprintf("test %d", i), test.GetByUint32(uint32(i)).String())
	}

	context.Globals().Set("test", test)

	result, err := context.Eval(`test.map(v => v.toUpperCase())`, EVAL_GLOBAL)
	require.NoError(t, err)
	defer result.Free()

	require.EqualValues(t, `TEST 0,TEST 1,TEST 2`, result.String())
}

func TestBadSyntax(t *testing.T) {
	runtime := NewRuntime()
	defer runtime.Free()

	context := runtime.NewContext()
	defer context.Free()

	_, err := context.Eval(`"bad syntax'`, EVAL_MODULE)
	require.Error(t, err)
}

func TestFunctionThrowError(t *testing.T) {
	expected := errors.New("expected error")

	runtime := NewRuntime()
	defer runtime.Free()

	context := runtime.NewContext()
	defer context.Free()

	context.Globals().SetFunction("A", func(ctx *Context, this Value, args []Value) Value {
		return ctx.ThrowError(expected)
	})

	_, actual := context.Eval("A()", EVAL_GLOBAL)
	require.Error(t, actual)
	require.EqualValues(t, "Error: "+expected.Error(), actual.Error())
}

func TestFunction(t *testing.T) {
	runtime := NewRuntime()
	defer runtime.Free()

	context := runtime.NewContext()
	defer context.Free()

	A := make(chan struct{})
	B := make(chan struct{})

	context.Globals().SetFunction("A", func(ctx *Context, this Value, args []Value) Value {
		require.Len(t, args, 4)
		require.True(t, args[0].IsString() && args[0].String() == "hello world!")
		require.True(t, args[1].IsNumber() && args[1].Int32() == 1)
		require.True(t, args[2].IsNumber() && args[2].Int64() == 8)
		require.True(t, args[3].IsNull())

		close(A)

		return ctx.String("A says hello")
	})

	context.Globals().SetFunction("B", func(ctx *Context, this Value, args []Value) Value {
		require.Len(t, args, 0)

		close(B)

		return ctx.Float64(256)
	})

	result, err := context.Eval(`A("hello world!", 1, 2 ** 3, null)`, EVAL_GLOBAL)
	require.NoError(t, err)
	defer result.Free()

	require.True(t, result.IsString() && result.String() == "A says hello")
	<-A

	result, err = context.Eval(`B()`, EVAL_GLOBAL)
	require.NoError(t, err)
	defer result.Free()

	require.True(t, result.IsNumber() && result.Uint32() == 256)
	<-B
}

func TestJsFunction(t *testing.T) {
	runtime := NewRuntime()
	defer runtime.Free()

	context := runtime.NewContext()
	defer context.Free()

	context.Globals().SetFunction("Callback", func(ctx *Context, this Value, args []Value) Value {
		require.Len(t, args, 1)
		require.True(t, args[0].IsFunction())

		return context.JsFunction(context.Null(), args[0], []Value{context.String("args test")})
	})

	result, err := context.Eval(`Callback(function(args){return args})`, EVAL_GLOBAL)
	require.NoError(t, err)
	defer result.Free()

	require.True(t, result.IsString() && result.String() == "args test")
}

func TestMemoryLimit(t *testing.T) {

	runtime := NewRuntime()
	defer runtime.Free()

	const kB = 1 << 10
	runtime.SetMemoryLimit(32 * kB)

	context := runtime.NewContext()
	defer context.Free()

	result, err := context.Eval(`var array = []; while (true) { array.push(null) }`, EVAL_GLOBAL)

	if assert.Error(t, err, "expected a memory limit violation") {
		require.Equal(t, "InternalError: out of memory", err.Error())
	}
	result.Free()
}
