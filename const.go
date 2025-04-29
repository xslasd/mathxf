package mathxf

import (
	"math"

	"github.com/shopspring/decimal"
)

var DefConst = map[string]*ValElement{
	"e":  NewConstValElement(decimal.NewFromFloat(math.E), false),
	"pi": NewConstValElement(decimal.NewFromFloat(math.Pi), false),

	"sum":   NewConstValElement(defSum, true),
	"avg":   NewConstValElement(defAvg, true),
	"max":   NewConstValElement(defMax, true),
	"min":   NewConstValElement(defMin, true),
	"cbrt":  NewConstValElement(defCbrt, true),
	"sqrt":  NewConstValElement(defSqrt, true),
	"round": NewConstValElement(defRound, true),
	"floor": NewConstValElement(defFloor, true),
	"ceil":  NewConstValElement(defCeil, true),
	"abs":   NewConstValElement(defAbs, true),
	"sin":   NewConstValElement(defSin, true),
	"cos":   NewConstValElement(defCos, true),
	"tan":   NewConstValElement(defTan, true),
	"asin":  NewConstValElement(defAsin, true),
	"acos":  NewConstValElement(defAcos, true),
	"atan":  NewConstValElement(defAtan, true),
	"atan2": NewConstValElement(defAtan2, true),
	"sinh":  NewConstValElement(defSinh, true),
	"cosh":  NewConstValElement(defCosh, true),
	"tanh":  NewConstValElement(defTanh, true),
	"asinh": NewConstValElement(defAsinh, true),
}

func defSum(ctx *EvaluatorContext, args ...*Value) (*Value, error) {
	alen := len(args)
	if ctx.IsHighPrecision {
		var sumV decimal.Decimal
		for _, item := range args {
			if !item.IsNumber() {
				return nil, ArgumentNotNumberErr.SetMessagef("sum")
			}
			if alen == 1 {
				return AsValue(item.Decimal()), nil
			}
			sumV = sumV.Add(item.Decimal())
		}
		return AsValue(sumV), nil
	}
	var sumV float64
	for _, item := range args {
		if !item.IsNumber() {
			return nil, ArgumentNotNumberErr.SetMessagef("sum")
		}
		if alen == 1 {
			return AsValue(item.Float()), nil
		}
		sumV += item.Float()
	}

	return AsValue(sumV), nil
}
func defAvg(ctx *EvaluatorContext, args ...*Value) (*Value, error) {
	alen := len(args)
	if alen == 0 {
		return nil, ArgumentNotNumberErr.SetMessagef("avg")
	}
	if ctx.IsHighPrecision {
		var rest []decimal.Decimal
		for ind, item := range args {
			if !item.IsNumber() {
				return nil, ArgumentNotNumberErr.SetMessagef("avg")
			}
			if ind == 0 {
				if alen == 1 {
					return AsValue(item.Decimal()), nil
				}
				continue
			}
			rest = append(rest, item.Decimal())
		}
		maxV := decimal.Avg(args[0].Decimal(), rest...)
		return AsValue(maxV), nil
	}
	var sumV float64
	for _, item := range args {
		if !item.IsNumber() {
			return nil, ArgumentNotNumberErr.SetMessagef("avg")
		}
		if alen == 1 {
			return AsValue(item.Float()), nil
		}
		sumV += item.Float()
	}
	return AsValue(sumV / float64(len(args))), nil
}

func defMax(ctx *EvaluatorContext, args ...*Value) (*Value, error) {
	alen := len(args)
	if alen == 0 {
		return nil, ArgumentNotEnoughErr.SetMessagef("max", ">=1", 0)
	}
	if ctx.IsHighPrecision {
		var rest []decimal.Decimal
		for ind, item := range args {
			if !item.IsNumber() {
				return nil, ArgumentNotNumberErr.SetMessagef("max", item.Val).SetCol(ind)
			}
			if ind == 0 {
				if alen == 1 {
					return AsValue(item.Decimal()), nil
				}
				continue
			}
			rest = append(rest, item.Decimal())
		}
		maxV := decimal.Max(args[0].Decimal(), rest...)
		return AsValue(maxV), nil
	}
	var maxV float64
	for _, item := range args {
		if !item.IsNumber() {
			return nil, ArgumentNotNumberErr.SetMessagef("max")
		}
		if alen == 1 {
			return AsValue(item.Float()), nil
		}
		maxV = math.Max(maxV, item.Float())
	}
	return AsValue(maxV), nil
}

func defMin(ctx *EvaluatorContext, args ...*Value) (*Value, error) {
	alen := len(args)
	if alen == 0 {
		return nil, ArgumentNotNumberErr.SetMessagef("min")
	}
	if ctx.IsHighPrecision {
		var rest []decimal.Decimal
		for ind, item := range args {
			if !item.IsNumber() {
				return nil, ArgumentNotNumberErr.SetMessagef("min")
			}
			if ind == 0 {
				if alen == 1 {
					return AsValue(item.Decimal()), nil
				}
				continue
			}
			rest = append(rest, item.Decimal())
		}
		minV := decimal.Min(args[0].Decimal(), rest...)
		return AsValue(minV), nil
	}
	var minV float64
	for _, item := range args {
		if !item.IsNumber() {
			return nil, ArgumentNotNumberErr.SetMessagef("min")
		}
		if alen == 1 {
			return item, nil
		}
		minV = math.Min(minV, item.Float())
	}
	return AsValue(minV), nil
}
func defCbrt(arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("cbrt")
	}
	return AsValue(math.Cbrt(arg.Float())), nil
}
func defSqrt(arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("cbrt")
	}
	return AsValue(math.Sqrt(arg.Float())), nil
}

func defRound(ctx *EvaluatorContext, arg *Value, n *Value) (*Value, error) {
	if !arg.IsNumber() || !n.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("round")
	}
	if ctx.IsHighPrecision {
		return AsValue(arg.Decimal().Round(int32(n.Integer()))), nil
	}
	p := math.Pow10(n.Integer())
	return AsValue(math.Round(arg.Float()*p) / p), nil
}
func defFloor(ctx *EvaluatorContext, arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("floor")
	}
	if ctx.IsHighPrecision {
		return AsValue(arg.Decimal().Floor()), nil
	}
	return AsValue(math.Floor(arg.Float())), nil
}
func defCeil(ctx *EvaluatorContext, arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("ceil")
	}
	if ctx.IsHighPrecision {
		return AsValue(arg.Decimal().Ceil()), nil
	}
	return AsValue(math.Ceil(arg.Float())), nil
}
func defAbs(ctx *EvaluatorContext, arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("abs")
	}
	if ctx.IsHighPrecision {
		return AsValue(arg.Decimal().Abs()), nil
	}
	return AsValue(math.Abs(arg.Float())), nil
}
func defSin(ctx *EvaluatorContext, arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("sin")
	}
	if ctx.IsHighPrecision {
		return AsValue(arg.Decimal().Sin()), nil
	}
	return AsValue(math.Sin(arg.Float())), nil
}
func defCos(ctx *EvaluatorContext, arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("cos")
	}
	if ctx.IsHighPrecision {
		return AsValue(arg.Decimal().Cos()), nil
	}
	return AsValue(math.Cos(arg.Float())), nil
}
func defTan(ctx *EvaluatorContext, arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("tan")
	}
	if ctx.IsHighPrecision {
		return AsValue(arg.Decimal().Tan()), nil
	}
	return AsValue(math.Tan(arg.Float())), nil
}
func defAsin(arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("asin")
	}
	return AsValue(math.Asin(arg.Float())), nil
}
func defAcos(arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("acos")
	}
	return AsValue(math.Acos(arg.Float())), nil
}
func defAtan(arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("atan")
	}
	return AsValue(math.Atan(arg.Float())), nil
}
func defAtan2(arg1, arg2 *Value) (*Value, error) {
	if !arg1.IsNumber() || !arg2.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("atan2")
	}
	return AsValue(math.Atan2(arg1.Float(), arg2.Float())), nil
}
func defSinh(arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("sinh")
	}
	return AsValue(math.Sinh(arg.Float())), nil
}
func defCosh(arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("cosh")
	}
	return AsValue(math.Cosh(arg.Float())), nil
}
func defTanh(arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("tanh")
	}
	return AsValue(math.Tanh(arg.Float())), nil
}
func defAsinh(arg *Value) (*Value, error) {
	if !arg.IsNumber() {
		return nil, ArgumentNotNumberErr.SetMessagef("asinh")
	}
	return AsValue(math.Asinh(arg.Float())), nil
}
