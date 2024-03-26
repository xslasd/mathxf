package mathxf

type INode interface {
	Execute(ctx *EvaluatorContext) error
}

type nodeDocument struct {
	Nodes []INode
}

func (n *nodeDocument) Execute(ctx *EvaluatorContext) error {
	for _, node := range n.Nodes {
		if err := node.Execute(ctx); err != nil {
			return ParseErr(err)
		}
	}
	return nil
}

type NodeWrapper struct {
	nodes []INode
}

func (wrapper *NodeWrapper) Execute(ctx *EvaluatorContext) error {
	for _, n := range wrapper.nodes {
		err := n.Execute(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

type NodeResData struct {
	name string
	evl  IEvaluator
}

func (n NodeResData) Execute(ctx *EvaluatorContext) error {
	val, err := n.evl.Evaluate(ctx)
	if err != nil {
		return err
	}
	ctx.ResultMap[ctx.defResultKey][n.name] = val
	return nil
}

type NodeAssignment struct {
	variable *variableResolver
	value    IEvaluator
}

func (n NodeAssignment) Execute(ctx *EvaluatorContext) error {
	return n.variable.SetPartValue(ctx, n.value)
}
