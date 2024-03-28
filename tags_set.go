package mathxf

type tagSetNode struct {
	name       string
	expression IEvaluator
}

func (t tagSetNode) Execute(ctx *EvaluatorContext) error {
	val, err := t.expression.Evaluate(ctx)
	if err != nil {
		return err
	}
	_, has := ctx.ValMap[t.name]
	_, isRes := ctx.ResultMap[t.name]
	if has || isRes {
		pos := t.expression.GetPositionToken()
		return VariableAlreadyExistsErr.SetMessagef(t.name).SetPosition(pos.line, pos.col)
	}
	ctx.ValMap[t.name] = NewPrivateValElement(val)
	return nil
}
func tagSetParser(parser *Parser) (INode, error) {
	next := parser.NextToken()
	if next.typ != TokenIdentifier {
		return nil, TokenNotIdentifierErr.SetMessagef(next.val).SetPosition(next.line, next.col)
	}
	assign := parser.NextToken()
	if assign.typ != TokenAssign {
		return nil, UnexpectedTokenErr.SetMessagef("set parser", assign.val).SetPosition(assign.line, assign.col)
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
