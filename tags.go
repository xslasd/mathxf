package mathxf

import "fmt"

type TagParser func(parser *Parser) (INode, error)

func (t *template) RegisterTag(name string, parserFn TagParser) error {
	_, ok := t.tags[name]
	if ok {
		return t.ParseErr()(TagRegisteredErr.SetMessagef(name))
	}
	fmt.Printf("registering tag '%s' \n", name)
	t.tags[name] = parserFn
	return nil

}
func defTags() map[string]TagParser {
	return map[string]TagParser{
		"if":  tagIfParser,
		"set": tagSetParser,
	}
}
