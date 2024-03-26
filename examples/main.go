package main

import (
	"fmt"
	"github.com/xslasd/mathxf"
	"math"
)

func main() {
	input := `
set f=3
if cc.aa.dd>1 and f<100{
  res.bb=3
}else{
res.cc=20
}
`
	tpl, err := mathxf.NewTemplate(input)
	if err != nil {
		panic(err)
	}
	m := map[string]any{
		"cc": map[string]any{
			"aa": map[string]any{
				"dd": 100,
			},
		},
	}
	tpl.SetPublicVarMap(m)
	tpl.AddFuncOrConst("ff", 100)
	_, err = tpl.Execute()
	if err != nil {
		fmt.Println("Execute:", err)
		return
	}
	for k, v := range tpl.PublicValMap() {
		//val, _ := v.Interface().(decimal.Decimal)
		////f := val.String()
		fmt.Printf("--Public.%s---%v\n", k, v.Val)
	}
	for k, v := range tpl.ResultValMap() {
		//val, _ := v.Interface().(decimal.Decimal)
		////f := val.String()
		fmt.Printf("--Result.%s---%v\n", k, v.Val)
	}
}

func Round(f float64, n int) float64 {
	p := math.Pow10(n)
	return math.Round(f*p) / p
}
