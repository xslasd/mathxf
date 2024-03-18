package mathxf

type INode interface {
	Execute(ctx EvaluatorContext) error
}

type nodeDocument struct {
	Nodes []INode
}

func (n *nodeDocument) Execute(ctx EvaluatorContext) error {
	for _, node := range n.Nodes {
		if err := node.Execute(ctx); err != nil {
			return err
		}
	}
	return nil
}

type NodeWrapper struct {
	nodes []INode
}

func (wrapper *NodeWrapper) Execute(ctx EvaluatorContext) error {
	for _, n := range wrapper.nodes {
		err := n.Execute(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
