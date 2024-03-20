package mathxf

import (
	"fmt"
	"math"
	"reflect"
)

var DefFunc = map[string]reflect.Value{
	"sum":   reflect.ValueOf(defSum),
	"avg":   reflect.ValueOf(defAvg),
	"max":   reflect.ValueOf(defMax),
	"min":   reflect.ValueOf(defMin),
	"cbrt":  reflect.ValueOf(defCbrt),
	"sqrt":  reflect.ValueOf(defSqrt),
	"round": reflect.ValueOf(defRound),
	"floor": reflect.ValueOf(defFloor),
	"ceil":  reflect.ValueOf(defCeil),
	"abs":   reflect.ValueOf(defAbs),
	"sin":   reflect.ValueOf(defSin),
	"cos":   reflect.ValueOf(defCos),
	"tan":   reflect.ValueOf(defTan),
	"asin":  reflect.ValueOf(defAsin),
	"acos":  reflect.ValueOf(defAcos),
	"atan":  reflect.ValueOf(defAtan),
	"atan2": reflect.ValueOf(defAtan2),
	"sinh":  reflect.ValueOf(defSinh),
	"cosh":  reflect.ValueOf(defCosh),
	"tanh":  reflect.ValueOf(defTanh),
	"asinh": reflect.ValueOf(defAsinh),
}

func defSum(args ...*Value) (*Value, error) {
	var sumV float64
	for _, item := range args {
		if item.IsDecimal() {
			sumV += item.Float()
		} else {
			return nil, fmt.Errorf("sum: argument must be a number")
		}
	}
	return AsValue(sumV), nil
}
func defAvg(args ...*Value) (*Value, error) {
	var sumV float64
	for _, item := range args {
		if item.IsDecimal() {
			sumV += item.Float()
		} else {
			return nil, fmt.Errorf("sum: argument must be a number")
		}
	}
	return AsValue(sumV / float64(len(args))), nil
}

func defMax(args ...*Value) (*Value, error) {
	var maxV float64
	for _, item := range args {
		if item.IsDecimal() {
			maxV = math.Max(maxV, item.Float())
		}
	}
	return AsValue(maxV), nil
}

func defMin(args ...*Value) (*Value, error) {
	var minV float64
	for _, item := range args {
		if item.IsDecimal() {
			minV = math.Min(minV, item.Float())
		}
	}
	return AsValue(minV), nil
}
func defCbrt(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("cbrt: argument must be a number")
	}
	return AsValue(math.Cbrt(arg.Float())), nil
}
func defSqrt(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("cbrt: argument must be a number")
	}
	return AsValue(math.Sqrt(arg.Float())), nil
}

func defRound(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("round: argument must be a number")
	}
	return AsValue(math.Round(arg.Float())), nil
}
func defFloor(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("floor: argument must be a number")
	}
	return AsValue(math.Floor(arg.Float())), nil
}
func defCeil(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("ceil: argument must be a number")
	}
	return AsValue(math.Ceil(arg.Float())), nil
}
func defAbs(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("ceil: argument must be a number")
	}
	return AsValue(math.Abs(arg.Float())), nil
}
func defSin(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("sin: argument must be a number")
	}
	return AsValue(math.Sin(arg.Float())), nil
}
func defCos(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("cos: argument must be a number")
	}
	return AsValue(math.Cos(arg.Float())), nil
}
func defTan(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("tan: argument must be a number")
	}
	return AsValue(math.Tan(arg.Float())), nil
}
func defAsin(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("asin: argument must be a number")
	}
	return AsValue(math.Asin(arg.Float())), nil
}
func defAcos(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("acos: argument must be a number")
	}
	return AsValue(math.Acos(arg.Float())), nil
}
func defAtan(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("atan: argument must be a number")
	}
	return AsValue(math.Atan(arg.Float())), nil
}
func defAtan2(arg1, arg2 *Value) (*Value, error) {
	if !arg1.IsDecimal() || !arg2.IsDecimal() {
		return nil, fmt.Errorf("atan2: arguments must be numbers")
	}
	return AsValue(math.Atan2(arg1.Float(), arg2.Float())), nil
}
func defSinh(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("sinh: argument must be a number")
	}
	return AsValue(math.Sinh(arg.Float())), nil
}
func defCosh(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("cosh: argument must be a number")
	}
	return AsValue(math.Cosh(arg.Float())), nil
}
func defTanh(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("tanh: argument must be a number")
	}
	return AsValue(math.Tanh(arg.Float())), nil
}
func defAsinh(arg *Value) (*Value, error) {
	if !arg.IsDecimal() {
		return nil, fmt.Errorf("asinh: argument must be a number")
	}
	return AsValue(math.Asinh(arg.Float())), nil
}
