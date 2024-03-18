package mathxf

import (
	"context"
	"fmt"
	"testing"
)

func TestLex(t *testing.T) {
	input := `if 1>0 {}`
	p, err := Parse(input)
	if err != nil {
		t.Error(err)
	}
	doc, err := p.ParseDocument()
	if err != nil {
		t.Error(err)
	}
	ctx := NewEvaluatorContext(context.Background())
	for _, v := range doc.Nodes {
		fmt.Println(v.Execute(ctx))
	}
}
