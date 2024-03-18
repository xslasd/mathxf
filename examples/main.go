package main

import (
	"context"
	"fmt"
	"github.com/xslasd/mathxf"
	"reflect"
)

func main() {
	input := `
set aa = 300
set bb = 0
if 200<aa and bb>=0 {
	res["cc"]=ff+3
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
	ctx.Private["ff"] = reflect.ValueOf(200)
	for _, v := range doc.Nodes {
		fmt.Println(v.Execute(ctx))
	}
	for k, v := range ctx.ResValues {
		fmt.Println("--ResValues---", k, v)
	}
}
