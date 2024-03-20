package mathxf

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math"
)

type IEvaluator interface {
	INode
	GetPositionToken() *Token
	Evaluate(ctx EvaluatorContext) (*Value, error)
}
type functionCallArgument interface {
	Evaluate(ctx EvaluatorContext) (*Value, error)
}

// Expression 处理TokenAnd 和 TokenOr
type Expression struct {
	expr1   IEvaluator
	expr2   IEvaluator
	opToken *Token
}

func (e Expression) Execute(ctx EvaluatorContext) error {
	_, err := e.expr1.Evaluate(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (e Expression) GetPositionToken() *Token {
	return e.expr1.GetPositionToken()
}

func (e Expression) Evaluate(ctx EvaluatorContext) (*Value, error) {
	v1, err := e.expr1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if e.expr2 != nil {
		switch e.opToken.typ {
		case TokenAnd:
			if !v1.IsTrue() {
				fmt.Println(" --------TokenAnd---v1---------", v1.IsTrue())
				return AsValue(false), nil
			} else {
				v2, err := e.expr2.Evaluate(ctx)
				fmt.Println(" --------TokenAnd---v2---------", v2.IsTrue())
				if err != nil {
					return nil, err
				}
				return AsValue(v2.IsTrue()), nil
			}
		case TokenOr:
			if v1.IsTrue() {
				return AsValue(true), nil
			} else {
				v2, err := e.expr2.Evaluate(ctx)
				if err != nil {
					return nil, err
				}
				return AsValue(v2.IsTrue()), nil
			}
		default:
			return nil, nil //ctx.Error(fmt.Sprintf("unimplemented: %name", r.opToken.Val), r.opToken)
		}
	} else {
		return v1, nil
	}
}

// relationalExpression 处理  TokenEqual  TokenNotEqual  TokenLess  TokenLessEqual  TokenGreater  TokenGreaterEqual
type relationalExpression struct {
	expr1   IEvaluator
	expr2   IEvaluator
	opToken *Token
}

func (r relationalExpression) Execute(ctx EvaluatorContext) error {
	v, err := r.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Println(v, "------relationalExpression---------", r.opToken.val)
	return nil
}

func (r relationalExpression) GetPositionToken() *Token {
	return r.expr1.GetPositionToken()
}

func (r relationalExpression) Evaluate(ctx EvaluatorContext) (*Value, error) {
	v1, err := r.expr1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}

	if r.expr2 == nil {
		return v1, nil

	}
	v2, err := r.expr2.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	switch r.opToken.typ {
	case TokenLessEquals:
		if ctx.IsHighPrecision {
			return AsValue(v1.Decimal().Cmp(v2.Decimal()) <= 0), nil
		}
		if v1.IsFloat() || v2.IsFloat() {
			return AsValue(v1.Float() <= v2.Float()), nil
		}
		if v1.IsTime() && v2.IsTime() {
			tm1, tm2 := v1.Time(), v2.Time()
			return AsValue(tm1.Before(tm2) || tm1.Equal(tm2)), nil
		}
		return AsValue(v1.Integer() <= v2.Integer()), nil
	case TokenGreatEquals:
		if ctx.IsHighPrecision {
			return AsValue(v1.Decimal().Cmp(v2.Decimal()) >= 0), nil
		}
		if v1.IsFloat() || v2.IsFloat() {
			return AsValue(v1.Float() >= v2.Float()), nil
		}
		if v1.IsTime() && v2.IsTime() {
			tm1, tm2 := v1.Time(), v2.Time()
			return AsValue(tm1.After(tm2) || tm1.Equal(tm2)), nil
		}
		return AsValue(v1.Integer() >= v2.Integer()), nil
	case TokenEquals:
		if ctx.IsHighPrecision {
			return AsValue(v1.Decimal().Cmp(v2.Decimal()) == 0), nil
		}
		return AsValue(v1.EqualValueTo(v2)), nil
	case TokenGreat:
		if ctx.IsHighPrecision {
			return AsValue(v1.Decimal().Cmp(v2.Decimal()) > 0), nil
		}
		if v1.IsFloat() || v2.IsFloat() {
			return AsValue(v1.Float() > v2.Float()), nil
		}
		if v1.IsTime() && v2.IsTime() {
			return AsValue(v1.Time().After(v2.Time())), nil
		}
		return AsValue(v1.Integer() > v2.Integer()), nil
	case TokenLess:
		if ctx.IsHighPrecision {
			return AsValue(v1.Decimal().Cmp(v2.Decimal()) < 0), nil
		}
		if v1.IsFloat() || v2.IsFloat() {
			return AsValue(v1.Float() < v2.Float()), nil
		}
		if v1.IsTime() && v2.IsTime() {
			return AsValue(v1.Time().Before(v2.Time())), nil
		}
		return AsValue(v1.Integer() < v2.Integer()), nil
	case TokenNotEquals:
		if ctx.IsHighPrecision {
			return AsValue(v1.Decimal().Cmp(v2.Decimal()) != 0), nil
		}
		return AsValue(!v1.EqualValueTo(v2)), nil
	case TokenIn:
		return AsValue(v2.Contains(v1)), nil
	default:
		return nil, fmt.Errorf("unimplemented: %v", r.opToken.val)
	}
}

// simpleExpression 处理 TokenAdd TokenSub
type simpleExpression struct {
	term1   IEvaluator
	term2   IEvaluator
	opToken *Token
}

func (s simpleExpression) Execute(ctx EvaluatorContext) error {
	v, err := s.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("simpleExpression---------------------------:%.16f\n", v.Float())
	return nil
}

func (s simpleExpression) GetPositionToken() *Token {
	return s.term1.GetPositionToken()
}

func (s simpleExpression) Evaluate(ctx EvaluatorContext) (*Value, error) {
	t1, err := s.term1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if s.term2 == nil {
		return t1, nil
	}
	t2, err := s.term2.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	switch s.opToken.typ {
	case TokenAdd:
		if t1.IsString() || t2.IsString() {
			// Result will be a string
			return AsValue(t1.String() + t2.String()), nil
		}
		if ctx.IsHighPrecision {
			fmt.Println(t1.Decimal(), t2.Decimal(), "--------IsHighPrecision>>>>>-------------", t1.Decimal().Add(t2.Decimal()))
			return AsValue(t1.Decimal().Add(t2.Decimal())), nil
		}
		if t1.IsFloat() || t2.IsFloat() {
			return AsValue(t1.Float() + t2.Float()), nil
		}
		return AsValue(t1.Integer() + t2.Integer()), nil
	case TokenSub:
		if ctx.IsHighPrecision {
			return AsValue(t1.Decimal().Sub(t2.Decimal())), nil
		}
		if t1.IsFloat() || t2.IsFloat() {
			return AsValue(t1.Float() - t2.Float()), nil
		}
		return AsValue(t1.Integer() - t2.Integer()), nil
	default:
		return nil, fmt.Errorf("simpleExpression unimplemented %s", s.GetPositionToken().String())
	}
}

// termExpression 处理 TokenMul TokenDiv TokenMod
type termExpression struct {
	factor1 IEvaluator
	factor2 IEvaluator
	opToken *Token
}

func (t termExpression) Execute(ctx EvaluatorContext) error {
	v, err := t.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Println("termExpression---------", v)
	return nil
}

func (t termExpression) GetPositionToken() *Token {
	return t.factor1.GetPositionToken()
}

func (t termExpression) Evaluate(ctx EvaluatorContext) (*Value, error) {
	f1, err := t.factor1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if t.factor2 == nil {
		return f1, nil
	}
	f2, err := t.factor2.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	switch t.opToken.typ {
	case TokenMul:
		if ctx.IsHighPrecision {
			return AsValue(f1.Decimal().Mul(f2.Decimal())), nil
		}
		if f1.IsFloat() || f2.IsFloat() {
			return AsValue(f1.Float() * f2.Float()), nil
		}
		return AsValue(f1.Integer() * f2.Integer()), nil
	case TokenDiv:
		if ctx.IsHighPrecision {
			divisor := f2.Decimal()
			//todo 精度控制
			if divisor.Cmp(decimal.Zero) == 0 {
				return nil, fmt.Errorf("float divide by zero ")
			}
			return AsValue(f1.Decimal().Div(divisor)), nil
		}

		if f1.IsFloat() || f2.IsFloat() {
			divisor := f2.Float()
			if divisor == 0 {
				return nil, fmt.Errorf("float divide by zero ")
			}
			return AsValue(f1.Float() / divisor), nil
		}
		divisor := f2.Integer()
		if divisor == 0 {
			return nil, fmt.Errorf("integer divide by zero")
		}
		return AsValue(f1.Integer() / divisor), nil
	case TokenMod:
		if ctx.IsHighPrecision {
			divisor := f2.Decimal()
			//todo 精度控制
			if divisor.Cmp(decimal.Zero) == 0 {
				return nil, fmt.Errorf("float divide by zero ")
			}
			return AsValue(f1.Decimal().Mod(divisor)), nil
		}
		divisor := f2.Integer()
		if divisor == 0 {
			return nil, fmt.Errorf("integer divide by zero")
		}
		return AsValue(f1.Integer() % divisor), nil
	default:
		return nil, fmt.Errorf("unimplemented %v", t.opToken)
	}

}

// powerExpression 处理 returns x**y, the base-x exponential of y.
type powerExpression struct {
	power1 IEvaluator
	power2 IEvaluator
}

func (p powerExpression) Execute(ctx EvaluatorContext) error {
	v, err := p.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Println("powerExpression--------", v)
	return nil
}

func (p powerExpression) GetPositionToken() *Token {
	return p.power1.GetPositionToken()
}

func (p powerExpression) Evaluate(ctx EvaluatorContext) (*Value, error) {
	p1, err := p.power1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if p.power2 == nil {
		return p1, nil
	}
	p2, err := p.power2.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if ctx.IsHighPrecision {
		return AsValue(p1.Decimal().Pow(p2.Decimal())), nil
	}
	return AsValue(math.Pow(p1.Float(), p2.Float())), nil
}
