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
	input := `
set aa =100
set bb = 0
if 200<max(aa,bb+400) {
	res.cc=pi-3
    res.dd=(-150/100)*0.17
    res.ff=bb/3
} else {
	res["cc"]=2+3
}`
	p, err := mathxf.Parse(input)
	if err != nil {
		panic(err)
	}
	doc, err := p.ParseDocument()
	if err != nil {
		panic(err)
	}

	ctx := mathxf.NewEvaluatorContext(context.Background())
	ctx.Private["ff"] = reflect.ValueOf(0.03)
	ctx.ResValues["res"] = make(mathxf.VarMap)
	err = doc.Execute(ctx)
	fmt.Println(err, "------------err")

	for k, v := range ctx.ResValues["res"] {
		val, _ := v.Interface().(decimal.Decimal)
		f, _ := val.Float64()
		fmt.Printf("%s--ResValues---%.16f\n", k, f)
	}
}

func Round(f float64, n int) float64 {
	p := math.Pow10(n)
	return math.Round(f*p) / p
}
