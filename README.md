# MathXF
[![license](https://badgen.net/github/license/xslasd/mathxf/)](https://github.com/xslasd/mathxf/blob/master/LICENSE)
[![release](https://badgen.net/github/release/xslasd/mathxf/stable)](https://github.com/xslasd/mathxf/releases)  
这是一个用 Go 语言编写的强大数学计算规则引擎，旨在提供灵活高效的数学表达式计算能力。它支持多种运算符，包括基本的算术运算（如加、减、乘、除）、指数运算、括号分组运算，甚至还包括逻辑条件运算符如 if 语句。该解析器允许用户输入并解析复杂的数学公式，精确计算出表达式的最终结果值。

## 功能特点：
1. **全面的运算符支持：** 囊括基础算术运算的同时，还融入了指数等高级运算功能。  
2. **条件逻辑集成：** 创新性地支持了 if 条件语句，使解析器能够处理含有条件逻辑的数学表达式。 
3. **高度可扩展性：** 采用模块化设计，易于扩展，可根据需要轻松添加新的运算符或函数。  
4. **准确性保障：** 运用了健壮的解析算法，确保即使是最复杂的公式也能正确解析并计算，同时保持浮点运算的高精度。 
5. **易用性佳：** 提供了一套简洁明了、用户友好的API接口，便于无缝集成到各类项目中。 
#### 支持操作符：
1. 算术运算符：+（加法）-（减法）*（乘法）/（除法）%（求余数或模运算）**（幂运算）
2. 比较运算符：<（小于）>（大于）<=（小于等于）>=（大于等于）==（等于）!= 或 <>（不等于）
3. 逻辑运算符：&& and（逻辑与）|| or（逻辑或） in (包含)  

**todo** 还没实现 ! not （逻辑非) 
#### 支持语法：
1. if条件判断： if<条件>{ }else if<条件>else{ } 
2. val定义变量：val a;val a,b,c; val a=1;var a,b,c=1 
3. 代码注释： //单行注释; /* */多行注释
4. 赋值操作： a=1; (**常量不能赋值**)  

#### 支持常量(可动态扩展)：
1. pi=math.Pi 
2. e=math.E 

如添加常量ff
```
	tpl.AddFuncOrConst("ff", 100)
```
#### 支持函数(可动态扩展):
 sum ,avg ,max ,min ,cbrt ,sqrt ,round ,floor ,ceil ,abs ,sin ,cos ,tan ,asin ,acos ,atan ,atan2 ,sinh ,cosh ,tanh ,asinh  

函数格式为 func(ctx *mathxf.EvaluatorContext,arg *mathxf.Value)(res1,error)  
ctx *EvaluatorContext 可以省略 
arg *mathxf.Value 可以多个或使用args ...*mathxf.Value  
返回参数必须 1-2个，第2个必须为error  

如添加函数 DblEquals
```
tpl.AddFuncOrConst("DblEquals", func(a, b *mathxf.Value) bool {
  return a.Decimal().Cmp(b.Decimal()) == 0
})
```

#### 支持模板内容替换：
将显示代码转换成可执行代码
```go
package main
import (
	"fmt"
	"github.com/xslasd/mathxf"
)
func main() {
    input := `
    if 用户显示变量 A > 10 {
      res.aa=5
    }
    `
    tpl, err := mathxf.NewTemplate(input)
    tpl.ReplaceStrMap(map[string]string{
    "用户显示变量 A":"UserA",
    })
    env := map[string]any{
    "UserA":100,
    }
    res, err := tpl.Execute(env)
    if err != nil {
        panic(err)
    }
    for k, v := range res {
        fmt.Printf("--Result.%s---%v\n", k, v)
    }
}
```
## 用法
go mod github.com/xslasd/mathxf

1. mathxf 默认开启精度计算功能，HighPrecision(false) 关闭。精度库使用：github.com/shopspring/decimal     
2. mathxf 计算结果默认放到map[string]map[string]*mathxf.Value中,默认前缀key为”res“,可以使用AddResultKeys(keys ...string)添加返回值前缀Key  
3. mathxf 计算结果map中, 前缀key为"env"中，存放着被修改env值。 
4. mathxf 支持直接计算，但不能和其它语法混用。 

#### 直接计算
```go
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
```
#### 组合用法
```go
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
```


