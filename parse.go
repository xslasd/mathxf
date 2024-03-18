package mathxf

import "fmt"

func Parse(tpl string) (*Parser, error) {
	l := lex(tpl)
	l.run()
	return &Parser{lex: l}, nil
}
func (p *Parser) ParseDocument() (*nodeDocument, error) {
	doc := &nodeDocument{
		Nodes: make([]INode, 0),
	}
	for p.peekToken().typ != TokenEOF {
		node, err := p.parseDocElement()
		if err != nil {
			return nil, err
		}
		doc.Nodes = append(doc.Nodes, node)
	}
	return doc, nil
}
func (p *Parser) parseDocElement() (INode, error) {
	t := p.peekToken()
	if t.typ == TokenIdentifier {
		tag1, ok := tags[t.val]
		if ok {
			p.nextToken()
			return tag1.parser(p)
		}
	}
	return p.ParseExpression()
}

func (p *Parser) WrapUntil() (*NodeWrapper, error) {
	peek := p.peekToken()
	if peek.typ == TokenLeftBigBrackets {
		p.nextToken()
		wrapper := new(NodeWrapper)
		for {
			t := p.peekToken()
			switch t.typ {
			case TokenEOF:
				return nil, fmt.Errorf("unclosed wrapper")
			case TokenRightBigBrackets:
				p.nextToken()
				return wrapper, nil
			}
			tag1, ok := tags[t.val]
			if ok {
				p.nextToken()
				tagNode, err := tag1.parser(p)
				if err != nil {
					return nil, err
				}
				wrapper.nodes = append(wrapper.nodes, tagNode)
				continue
			}
			n := p.nextToken()
			vRes, err := p.ParseVariable(n)
			if err != nil {
				return nil, err
			}
			next := p.nextToken()
			if next.typ != TokenAssign {
				return nil, fmt.Errorf("assign unexpected token %s", next)
			}
			exp2, err := p.ParseExpression()
			if err != nil {
				return nil, err
			}
			assign := &assignmentResolver{
				vRes,
				exp2,
			}
			fmt.Println("-----add----assignmentResolver-------------")
			wrapper.nodes = append(wrapper.nodes, assign)
		}
	}
	return nil, fmt.Errorf("WrapUntil unexpected token %s", peek.val)
}

type Parser struct {
	lex *lexer

	peekTokens [3]Token // three-token lookahead for Parser.
	peekCount  int
}

// errorf formats the error and terminates processing.
func (p *Parser) errorf(format string, args ...interface{}) {
	format = fmt.Sprintf("template %d: %s", p.lex.line, format)
	panic(fmt.Errorf(format, args...))
}

// error terminates processing.
func (p *Parser) error(err error) {
	p.errorf("%name", err)
}

func (p *Parser) peekToken() Token {
	if p.peekCount > 0 {
		return p.peekTokens[p.peekCount-1]
	}
	p.peekCount = 1
	p.peekTokens[0] = p.lex.nextToken()
	return p.peekTokens[0]
}

func (p *Parser) nextToken() Token {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.peekTokens[0] = p.lex.nextToken()
	}
	return p.peekTokens[p.peekCount]
}

func (p *Parser) backup() {
	p.peekCount++
}
