package mathxf

import (
	"context"
	"reflect"
)

type VarMap map[string]reflect.Value

type EvaluatorContext struct {
	context.Context
	Public    VarMap
	Private   VarMap
	ResValues VarMap
}

func NewEvaluatorContext(ctx context.Context) EvaluatorContext {
	return EvaluatorContext{
		Context:   ctx,
		Public:    make(VarMap),
		Private:   make(VarMap),
		ResValues: make(VarMap),
	}
}
