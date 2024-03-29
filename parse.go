package mathxf

import (
	"fmt"
)

func Parse(tpl string) (*Parser, error) {
	l := lex(tpl)
	l.run()
	return &Parser{lex: l, tags: defTags()}, nil
}
func (p *Parser) ParseDocument() (*nodeDocument, error) {
	doc := &nodeDocument{
		Nodes: make([]INode, 0),
	}
	ind := 1
	for p.PeekToken().typ != TokenEOF {
		node, err := p.parseDocElement(ind)
		if err != nil {
			return nil, ParseErr(err)
		}
		doc.Nodes = append(doc.Nodes, node)
		ind++
	}
	return doc, nil
}
func (p *Parser) parseDocElement(ind int) (INode, error) {
	t := p.PeekToken()
	if t.typ == TokenIdentifier {
		tagParser, ok := p.tags[t.val]
		if ok {
			p.NextToken()
			return tagParser(p)
		}
	} else {
		evl, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		return NodeResData{name: fmt.Sprintf("res%d", ind), evl: evl}, err
	}
	vRes, err := p.ParseVariable(p.NextToken())
	if err != nil {
		return nil, err
	}
	next := p.NextToken()
	if next.typ != TokenAssign {
		p.Backup()
		return NodeResData{name: fmt.Sprintf("res%d", ind), evl: vRes}, nil
	}

	exp2, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	return &NodeAssignment{
		vRes,
		exp2,
	}, nil
}

func (p *Parser) WrapUntil() (*NodeWrapper, error) {
	peek := p.PeekToken()
	if peek.typ == TokenLeftBigBrackets {
		p.NextToken()
		wrapper := new(NodeWrapper)
		for {
			t := p.PeekToken()
			switch t.typ {
			case TokenEOF:
				return nil, WrapperUnclosedErr.SetPosition(t.line, t.col)
			case TokenRightBigBrackets:
				p.NextToken()
				return wrapper, nil
			}
			tagParser, ok := p.tags[t.val]
			if ok {
				p.NextToken()
				tagNode, err := tagParser(p)
				if err != nil {
					return nil, err
				}
				wrapper.nodes = append(wrapper.nodes, tagNode)
				continue
			}
			n := p.NextToken()
			assign, err := p.ParseAssignment(n)
			if err != nil {
				return nil, err
			}
			wrapper.nodes = append(wrapper.nodes, assign)
		}
	}
	return nil, UnexpectedTokenErr.SetMessagef("wrapUntil", peek.val).SetPosition(peek.line, peek.col)
}

type Parser struct {
	lex *lexer

	peekTokens [3]Token // three-token lookahead for Parser.
	peekCount  int

	tags map[string]TagParser
}

func (p *Parser) PeekToken() Token {
	if p.peekCount > 0 {
		return p.peekTokens[p.peekCount-1]
	}
	p.peekCount = 1
	p.peekTokens[0] = p.lex.nextToken()
	return p.peekTokens[0]
}

func (p *Parser) NextToken() Token {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.peekTokens[0] = p.lex.nextToken()
	}
	return p.peekTokens[p.peekCount]
}

func (p *Parser) Backup() {
	p.peekCount++
}
func (p *Parser) Ignore() {
	p.peekCount--
}
