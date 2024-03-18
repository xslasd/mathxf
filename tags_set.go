package mathxf

import (
	"fmt"
	"reflect"
)

type tagSetNode struct {
	name       string
	expression IEvaluator
}

func (t tagSetNode) Execute(ctx EvaluatorContext) error {
	val, err := t.expression.Evaluate(ctx)
	if err != nil {
		return err
	}
	ctx.Private[t.name] = reflect.ValueOf(val)
	return nil
}
func tagSetParser(parser *Parser) (INode, error) {
	next := parser.nextToken()
	if next.typ != TokenIdentifier {
		return nil, fmt.Errorf("expected identifier, got %s", next.typ)
	}
	if parser.nextToken().typ != TokenAssign {
		return nil, fmt.Errorf("expected =, got %s", next.typ)
	}
	expression, err := parser.ParseExpression()
	if err != nil {
		return nil, err
	}
	return tagSetNode{
		name:       next.val,
		expression: expression,
	}, nil

}

func init() {
	RegisterTag("set", tagSetParser)
}
