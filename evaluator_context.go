package mathxf

import (
	"context"
	"math"
	"reflect"
)

type VarMap map[string]reflect.Value

type EvaluatorContext struct {
	context.Context
	IsHighPrecision bool
	Public          VarMap
	Private         VarMap
	ResValues       map[string]VarMap
}

func NewEvaluatorContext(ctx context.Context) EvaluatorContext {
	public := DefFunc
	public["pi"] = reflect.ValueOf(math.Pi)
	return EvaluatorContext{
		Context:         ctx,
		IsHighPrecision: true,
		Public:          public,
		Private:         make(VarMap),
		ResValues:       make(map[string]VarMap),
	}
}
