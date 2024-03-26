package mathxf

import (
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
)

type ValMap map[string]reflect.Value
type ValElementMap map[string]*ValElement

type ValElement struct {
	ValType ValType
	IsSet   bool
	IsFunc  bool
	Val     reflect.Value
}

func NewConstValElement(val any, isFunc bool) *ValElement {
	return &ValElement{
		ValType: ConstVal,
		IsFunc:  isFunc,
		Val:     reflect.ValueOf(val),
	}
}

func NewPublicValElement(val any) *ValElement {
	return &ValElement{
		ValType: PublicVal,
		Val:     reflect.ValueOf(val),
	}
}

func NewPrivateValElement(val any) *ValElement {
	return &ValElement{
		ValType: PrivateVal,
		Val:     reflect.ValueOf(val),
	}
}
func NewResultValElement(val any) *ValElement {
	return &ValElement{
		ValType: ResultVal,
		Val:     reflect.ValueOf(val),
	}
}

type ParseECodeFn func(error) error

type ValType int

const (
	PublicVal ValType = iota
	PrivateVal
	ConstVal
	ResultVal
)

type EvaluatorContext struct {
	context.Context
	IsHighPrecision bool
	ValMap          ValElementMap
	//ResultMap       map[string]ValMap

	defResultKey string
	parseErrFn   ParseECodeFn
}

func NewEvaluatorContext(ctx context.Context) *EvaluatorContext {
	valMap := DefFunc
	valMap["pi"] = NewConstValElement(math.Pi, false)
	fmt.Printf("registering const '%s' \n", "pi")
	res := EvaluatorContext{
		Context:         ctx,
		IsHighPrecision: true,
		ValMap:          valMap,
		//ResultMap:       make(map[string]ValMap),
		defResultKey: "res",
		parseErrFn:   ParseErr,
	}
	return &res
}

func ParseErr(err error) error {
	if err == nil {
		return nil
	}
	e := Cause(err)
	line, col := e.Position()
	return errors.New(fmt.Sprintf("line: %d, col: %d, %s", line, col, e.Message()))
}
