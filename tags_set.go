package mathxf

type tagSetNode struct {
	setNodes []*SetNode
	isAssign bool
}

type SetNode struct {
	name       string
	expression IEvaluator
}

func (t tagSetNode) Execute(ctx *EvaluatorContext) error {
	var val *Value
	var err error
	for index, set := range t.setNodes {
		if index == 0 || !t.isAssign {
			val, err = set.expression.Evaluate(ctx)
			if err != nil {
				return err
			}
		}
		_, has := ctx.ValMap[set.name]
		_, isRes := ctx.ResultMap[set.name]
		if has || isRes {
			pos := set.expression.GetPositionToken()
			return VariableAlreadyExistsErr.SetMessagef(set.name).SetPosition(pos.line, pos.col)
		}
		ctx.ValMap[set.name] = NewPrivateValElement(val)
	}
	return nil
}
func tagSetParser(parser *Parser) (INode, error) {
	res := tagSetNode{
		setNodes: make([]*SetNode, 0),
	}
	var isAssign bool
	var isComma bool
	setNameArr := make([]string, 0)
	for {
		next := parser.NextToken()
		if next.typ != TokenIdentifier {
			parser.Backup()
			break
		}
		_, ok := TokenKeywords[next.val]
		if ok {
			if isComma {
				parser.Backup()
				break
			}
			return nil, VariableIsKeywordErr.SetMessagef(next.val).SetPosition(next.line, next.col)
		}
		setNameArr = append(setNameArr, next.val)
		assign := parser.NextToken()
		if assign.typ == TokenAssign {
			isAssign = true
			break
		}
		if assign.typ != TokenComma {
			parser.Backup()
		}
		isComma = true
	}
	if len(setNameArr) == 0 {
		peek := parser.PeekToken()
		return nil, TokenNotIdentifierErr.SetMessagef(peek.val).SetPosition(peek.line, peek.col)
	}
	var exp IEvaluator
	var err error
	if isAssign {
		exp, err = parser.ParseExpression()
		if err != nil {
			return nil, err
		}
	} else {
		exp = &numberResolver{
			val: 0,
		}
	}
	for _, name := range setNameArr {
		res.setNodes = append(res.setNodes, &SetNode{
			name:       name,
			expression: exp,
		})
	}
	res.isAssign = isAssign
	return res, nil
}

func setNodeParser(parser *Parser) (*SetNode, bool, error) {
	next := parser.NextToken()
	assign := parser.NextToken()
	if assign.typ != TokenAssign {
		res := &SetNode{
			name: next.val,
			expression: &numberResolver{
				locationToken: &next,
				val:           0,
			}}
		if assign.typ != TokenComma {
			parser.Backup()
			return res, false, nil
		} else {
			return res, true, nil
		}
	}
	expression, err := parser.ParseExpression()
	if err != nil {
		return nil, false, err
	}
	return &SetNode{
		name:       next.val,
		expression: expression,
	}, false, nil
}
