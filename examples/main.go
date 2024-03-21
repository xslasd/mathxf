package main

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/xslasd/mathxf"
	"math"
	"reflect"
)

func main() {
	input := `ssss`
	p, err := mathxf.Parse(input)
	if err != nil {
		fmt.Println("---------err != nil---------", err)
		panic(err)
	}
	doc, err := p.ParseDocument()
	if err != nil {
		code := mathxf.Cause(err)
		line, col := code.Position()
		fmt.Println(err, "--------", code.Message(), "------------err", line, col)
		return
	}

	ctx := mathxf.NewEvaluatorContext(context.Background())
	ctx.Private["ff"] = reflect.ValueOf(0.03)
	ctx.Private["cc"] = reflect.ValueOf(33)
	ctx.ResValues["res"] = make(mathxf.VarMap)
	err = doc.Execute(ctx)
	if err != nil {
		code := mathxf.Cause(err)
		line, col := code.Position()
		fmt.Println(err, "--------", code.Message(), "------------err", line, col)
		return
	}

	for k, v := range ctx.ResValues["res"] {
		//val, _ := v.Interface().(decimal.Decimal)
		////f := val.String()

		fmt.Printf("%s--ResValues---%v %T\n", k, v, v.Interface())
	}
	value := decimal.NewFromFloat(123.456789)

	// 保留3位小数
	roundedValue := value.Round(3)

	// 输出结果
	fmt.Println(roundedValue.Float64())
}

func Round(f float64, n int) float64 {
	p := math.Pow10(n)
	return math.Round(f*p) / p
}
