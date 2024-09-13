package main

import (
	"fmt"
	"syscall/js"
)

func JsFunc(fn func()) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn()
		return nil
	})
}
func JsFuncOut(fn func() any) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return fn()
	})
}
func JsFuncIn(fn func(args []js.Value)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn(args)
		return nil
	})
}

func Method(methodName string, args []js.Value) struct {
	StringArg func(int, string) (string, error)
	FloatArg  func(int, string) (float64, error)
	IntArg    func(int, string) (int64, error)
} {
	return struct {
		StringArg func(int, string) (string, error)
		FloatArg  func(int, string) (float64, error)
		IntArg    func(int, string) (int64, error)
	}{
		StringArg: func(index int, argName string) (string, error) {
			return StringArg(args, index, argName, methodName)
		},
		FloatArg: func(index int, argName string) (float64, error) {
			return NumberArg[float64](args, index, argName, methodName)
		},
		IntArg: func(index int, argName string) (int64, error) {
			return NumberArg[int64](args, index, argName, methodName)
		},
	}
}

func StringArg(args []js.Value, index int, argName string, methodName string) (string, error) {
	if len(args) <= index {
		return "", fmt.Errorf("%s() missing argument: %s", methodName, argName)
	}
	if args[index].Type() != js.TypeString {
		return "", fmt.Errorf("%s() invalid argument %s value: %v, expected number", methodName, argName, args[index])
	}
	return args[index].String(), nil
}

func NumberArg[V int64 | float64](args []js.Value, index int, argName string, methodName string) (V, error) {
	if len(args) <= index {
		return 0, fmt.Errorf("%s() missing argument: %s", methodName, argName)
	}
	if args[index].Type() != js.TypeNumber {
		return 0, fmt.Errorf("%s() invalid argument %s value: %v, expected string", methodName, argName, args[index])
	}
	return V(args[index].Float()), nil
}
