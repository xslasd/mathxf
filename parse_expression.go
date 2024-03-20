package mathxf

import (
	"errors"
	"fmt"
	"strconv"
)

func (p *Parser) ParseAssignment() (IEvaluator, error) {
	varName, err := p.parseVariableOrLiteral()
	if err != nil {
		return nil, err
	}
	variable, ok := varName.(*variableResolver)
	if !ok {
		return nil, fmt.Errorf("expected a variable, got %T", varName)
	}
	peek := p.peekToken()
	if peek.typ != TokenAssign {
		return nil, fmt.Errorf("expected '=' after variable, got %s", peek.val)
	}
	p.nextToken()
	exp, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	return &assignmentResolver{
		variable: variable,
		value:    exp,
	}, nil
}

func (p *Parser) ParseExpression() (IEvaluator, error) {
	expr1, err := p.parseRelationalExpression()
	if err != nil {
		return nil, err
	}
	exp := &Expression{
		expr1: expr1,
	}
	peek := p.peekToken()
	if peek.typ == TokenAnd || peek.typ == TokenOr {
		op := p.nextToken()
		expr2, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		exp.expr2 = expr2
		exp.opToken = &op
		return exp, nil
	}
	return exp.expr1, nil
}
func (p *Parser) parseRelationalExpression() (IEvaluator, error) {
	expr1, err := p.parseSimpleExpression()
	if err != nil {
		return nil, err
	}
	expr := &relationalExpression{
		expr1: expr1,
	}
	peek := p.peekToken()
	switch peek.typ {
	case TokenEquals, TokenNotEquals, TokenGreat, TokenGreatEquals, TokenLess, TokenLessEquals:
		op := p.nextToken()
		expr2, err := p.parseRelationalExpression()
		if err != nil {
			return nil, err
		}
		fmt.Println(expr1.GetPositionToken().val, "-----------------", op.val, "----", expr2.GetPositionToken().val)
		expr.expr2 = expr2
		expr.opToken = &op
		return expr, nil
	case TokenIn:
		op := p.nextToken()
		expr2, err := p.parseSimpleExpression()
		if err != nil {
			return nil, err
		}
		expr.expr2 = expr2
		expr.opToken = &op
		return expr, nil
	default:
		return expr.expr1, nil
	}
}
func (p *Parser) parseSimpleExpression() (IEvaluator, error) {
	term1, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	expr := &simpleExpression{
		term1: term1,
	}
	for {
		peek := p.peekToken()
		switch peek.typ {
		case TokenAdd, TokenSub:
			if expr.opToken != nil {
				expr = &simpleExpression{
					term1: expr,
				}
			}
			op := p.nextToken()
			term2, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			expr.term2 = term2
			expr.opToken = &op
		default:
			if expr.term2 == nil {
				return expr.term1, nil
			}
			return expr, nil
		}
	}

}
func (p *Parser) parseTerm() (IEvaluator, error) {
	factor1, err := p.parsePower()
	if err != nil {
		return nil, err
	}
	termObj := &termExpression{
		factor1: factor1,
	}
	for {
		peek := p.peekToken()
		switch peek.typ {
		case TokenMul, TokenDiv, TokenMod:
			if termObj.opToken != nil {
				termObj = &termExpression{
					factor1: termObj,
				}
			}
			op := p.nextToken()
			factor2, err := p.parsePower()
			if err != nil {
				return nil, err
			}
			termObj.factor2 = factor2
			termObj.opToken = &op
		default:
			if termObj.factor2 == nil {
				return termObj.factor1, nil
			}
			return termObj, nil
		}
	}
}
func (p *Parser) parsePower() (IEvaluator, error) {
	power1, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	powerObj := &powerExpression{
		power1: power1,
	}
	if p.peekToken().typ == TokenPow {
		p.nextToken()
		power2, err := p.parsePower()
		if err != nil {
			return nil, err
		}
		powerObj.power2 = power2
		return powerObj, nil
	}
	return powerObj.power1, nil
}
func (p *Parser) parseFactor() (IEvaluator, error) {
	if p.peekToken().typ == TokenLeftParen {
		p.nextToken()
		expr, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		if p.peekToken().typ == TokenRightParen {
			p.nextToken()
			return expr, nil
		}
		return nil, errors.New("expect ')' expected after expression")
	}
	return p.parseVariableOrLiteral()
}
func (p *Parser) parseVariableOrLiteral() (IEvaluator, error) {
	t := p.nextToken()
	if t.typ == TokenEOF {
		return nil, errors.New("unexpected EOF, expected a number, string, keyword or identifier")
	}
	switch t.typ {
	case TokenNumber:
		f, err := strconv.ParseFloat(t.val, 64)
		if err != nil {
			return nil, err
		}
		fr := &numberResolver{
			locationToken: &t,
			val:           f,
		}
		return fr, nil
	case TokenBool:
		b, err := strconv.ParseBool(t.val)
		if err != nil {
			return nil, err
		}
		br := &boolResolver{
			locationToken: &t,
			val:           b,
		}
		return br, nil
	case TokenString:
		s := &stringResolver{
			locationToken: &t,
			val:           t.val,
		}
		return s, nil
	case TokenLeftBrackets:
		arr := &arrayResolver{
			locationToken: &t,
		}
		if p.peekToken().typ == TokenRightBrackets {
			p.nextToken()
			return arr, nil
		}
		for {
			if p.peekToken().typ == TokenEOF {
				return nil, errors.New("unexpected EOF, expected a number, string, keyword or identifier")
			}
			exprArg, err := p.ParseExpression()
			if err != nil {
				return nil, err
			}
			arr.parts = append(arr.parts, &variablePart{
				typ:       VariablePartTypeArray,
				subscript: exprArg,
			})
			if p.peekToken().typ == TokenRightBrackets {
				p.nextToken()
				break
			}
			if p.nextToken().typ != TokenComma {
				return nil, errors.New("missing comma or closing bracket after argument")
			}
		}
		return arr, nil
	}
	return p.ParseVariable(t)
}

func (p *Parser) ParseVariable(t Token) (*variableResolver, error) {
	if t.typ != TokenIdentifier {
		return nil, errors.New("unexpected token, expected a number, string, keyword or identifier")
	}
	resolver := &variableResolver{
		locationToken: &t,
	}
	resolver.parts = append(resolver.parts, &variablePart{
		typ:  VariablePartTypeIdent,
		name: t.val,
	})
	for {
		next := p.nextToken()
		switch next.typ {
		case TokenField:
			resolver.parts = append(resolver.parts, &variablePart{
				typ:  VariablePartTypeIdent,
				name: next.val,
			})
		case TokenLeftBrackets:
			exprSubscript, err := p.ParseExpression()
			if err != nil {
				return nil, err
			}
			resolver.parts = append(resolver.parts, &variablePart{
				typ:       VariablePartTypeSubscript,
				subscript: exprSubscript,
			})
			if p.nextToken().typ != TokenRightBrackets {
				return nil, errors.New("missing closing bracket after subscript")
			}
		case TokenLeftParen:
			funcPart := resolver.parts[len(resolver.parts)-1]
			funcPart.isFunctionCall = true
		argumentLoop:
			for {
				peek := p.peekToken()
				if peek.typ == TokenEOF {
					return nil, errors.New("unexpected EOF, expected a number, string, keyword or identifier")
				}
				if peek.typ == TokenRightParen {
					p.nextToken()
					break argumentLoop
				}
				exprArg, err := p.ParseExpression()
				if err != nil {
					return nil, err
				}
				funcPart.callingArgs = append(funcPart.callingArgs, exprArg)
				next2 := p.nextToken()
				fmt.Println(next2, "--------------------------")
				if next2.typ == TokenRightParen {
					break argumentLoop
				}
				if next2.typ != TokenComma {
					p.nextToken()
					return nil, errors.New("missing comma or closing bracket after argument")
				}
			}
		default:
			p.backup()
		}
		break
	}
	return resolver, nil
}
