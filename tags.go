package mathxf

import "fmt"

type TagParser func(parser *Parser) (INode, error)

type tag struct {
	name   string
	parser TagParser
}

var tags map[string]*tag

func init() {
	tags = make(map[string]*tag)
}
func RegisterTag(name string, parserFn TagParser) error {
	_, ok := tags[name]
	if ok {
		return fmt.Errorf("tag with name '%s' is already registered", name)
	}
	fmt.Printf("registering tag '%s' %T \n", name, parserFn)
	tags[name] = &tag{
		name:   name,
		parser: parserFn,
	}
	return nil
}
