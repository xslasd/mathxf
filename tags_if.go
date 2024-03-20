package mathxf

import "fmt"

type tagIfNode struct {
	conditions []IEvaluator
	wrappers   []*NodeWrapper
}

func (t *tagIfNode) Execute(ctx EvaluatorContext) error {
	fmt.Println("========tagIfNode=======Execute====")
	cLength := len(t.conditions)
	wLength := len(t.wrappers)

	for index, condition := range t.conditions {
		res, err := condition.Evaluate(ctx)
		if err != nil {
			return err
		}
		if res.IsTrue() {
			fmt.Println("========tagIfNode======IsTrue=====", res.IsFloat())
			return t.wrappers[index].Execute(ctx)
		}
		last := index + 1
		if cLength == last && wLength > last {
			return t.wrappers[last].Execute(ctx)
		}
	}
	return nil
}

func tagIfParser(parser *Parser) (INode, error) {
	condition, err := parser.ParseExpression()
	if err != nil {
		return nil, err
	}
	ifNode := new(tagIfNode)
	ifNode.conditions = append(ifNode.conditions, condition)
	for {
		wrapper, err := parser.WrapUntil()
		if err != nil {
			return nil, err
		}
		ifNode.wrappers = append(ifNode.wrappers, wrapper)
		if parser.peekToken().val == "else" {
			parser.nextToken()
			if parser.peekToken().val == "if" {
				parser.nextToken()
				elseIfCondition, err := parser.ParseExpression()
				if err != nil {
					return nil, err
				}
				ifNode.conditions = append(ifNode.conditions, elseIfCondition)
				continue
			}
			elseWrapper, err := parser.WrapUntil()
			if err != nil {
				return nil, err
			}
			ifNode.wrappers = append(ifNode.wrappers, elseWrapper)
		}
		return ifNode, nil
	}
}
func init() {
	RegisterTag("if", tagIfParser)
}
