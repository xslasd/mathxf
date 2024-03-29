package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/xslasd/mathxf"
)

func main() {
	input := `
// 用户优惠劵面值计算规则
if 用户A.总订单数 > 5 && 特价商品区.单价小于10元.订单数<=3 {
   优惠劵面值= 优惠算法((用户A.订单总金额-特价商品区.单价小于10元.订单数*10)*优惠比率最大值,用户A.用户等级)
}
`

	tpl, err := mathxf.NewTemplate(input)
	if err != nil {
		panic(err)
	}
	replaceStrMap := map[string]string{
		"用户A.总订单数":          "TotalOrders",
		"用户A.用户等级":          "UserLevel",
		"用户A.订单总金额":         "OrderTotalAmount",
		"特价商品区.单价小于10元.订单数": "SpecialOffer_10_Orders",
		"优惠比率最大值":           "DiscountRatioMax",
		"优惠劵面值":             "CouponFace",

		"优惠算法": "DiscountAlgorithm",
	}
	tpl.ReplaceStrMap(replaceStrMap)
	tpl.AddFuncOrConst("DiscountAlgorithm", func(a, b *mathxf.Value) decimal.Decimal {
		level := b.Integer()
		if level > 5 {
			return a.Decimal().Mul(decimal.NewFromFloat(0.5))
		} else if level > 3 {
			return a.Decimal().Mul(decimal.NewFromFloat(0.3))
		} else if level >= 1 {
			return a.Decimal().Mul(decimal.NewFromFloat(0.1))
		}
		return decimal.NewFromFloat(0)
	})
	tpl.AddFuncOrConst("DiscountRatioMax", 0.3)
	env := map[string]any{
		"TotalOrders":            10,
		"UserLevel":              1,
		"OrderTotalAmount":       100,
		"SpecialOffer_10_Orders": 3,

		"DiscountRatioMax": 0.3,

		"CouponFace": 0,
	}
	res, err := tpl.Execute(env)
	if err != nil {
		fmt.Println("Execute:", err)
		return
	}
	for k, v := range tpl.PublicValMap() {
		fmt.Printf("--Public.%s---%v\n", k, v.Val)
	}
	for k, v := range res {
		fmt.Printf("--Result.%s---%v\n", k, v)
	}
}
