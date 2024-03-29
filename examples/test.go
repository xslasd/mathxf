package main

import (
	"fmt"
	"github.com/xslasd/mathxf"
)

func main() {
	input := `1 + 2 * 6 / 4 + (456 - 8 * 9.2) - (2 + 4 ^ 5)`
	tpl, err := mathxf.NewTemplate(input)
	res, err := tpl.Execute(nil)
	if err != nil {
		panic(err)
	}
	for k, v := range res {
		fmt.Printf("--Result.%s---%v\n", k, v)
	}
}
