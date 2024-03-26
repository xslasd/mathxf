package mathxf

import (
	"context"
	"math"
	"reflect"
)

type template struct {
	tpl   string // the string being scanned
	ctx   *EvaluatorContext
	tags  map[string]TagParser
	debug bool
}

func (t *template) Debug(b bool) {
	t.debug = b
}

func (t *template) AddFuncOrConst(name string, val any) error {
	if _, ok := t.ctx.ValMap[name]; ok {
		return t.ParseErr()(ConstRegisteredErr.SetMessagef(name))
	}
	switch reflect.ValueOf(val).Kind() {
	case reflect.Func:
		t.ctx.ValMap[name] = NewConstValElement(val, true)
	default:
		t.ctx.ValMap[name] = NewConstValElement(val, false)
	}
	return nil
}

func (t *template) SetDefResultKey(key string) {
	t.ctx.defResultKey = key
}

func (t *template) SetParseErrFn(fn ParseECodeFn) {
	t.ctx.parseErrFn = fn
}
func (t *template) SetPublicVarMap(values map[string]any) {
	for k, v := range values {
		t.ctx.ValMap[k] = NewPublicValElement(v)
	}
}
func (t *template) SetContext(ctx context.Context) {
	t.ctx.Context = ctx
}

//	func (t *template) AddResultKey(key string) error {
//		_, ok := t.ctx.ResultMap[key]
//		if ok {
//			return t.ParseErr()(ResultKeyRegisteredErr.SetMessagef(key))
//		}
//		t.ctx.ResultMap[key] = make(ValMap)
//		return nil
//	}
func (t *template) HighPrecision(b bool) {
	t.ctx.IsHighPrecision = b
}

func NewTemplate(tpl string) (*template, error) {
	t := &template{
		tpl:  tpl,
		tags: defTags(),
		ctx: &EvaluatorContext{
			Context:         context.TODO(),
			IsHighPrecision: true,
			ValMap:          DefFunc,
			ResultMap:       make(map[string]ValMap),
			defResultKey:    "res",
			parseErrFn:      ParseErr,
		},
	}
	t.ctx.ValMap["pi"] = NewConstValElement(math.Pi, false)
	t.ctx.ResultMap[t.ctx.defResultKey] = make(ValMap)
	return t, nil
}
func (t *template) ParseErr() ParseECodeFn {
	return t.ctx.parseErrFn
}
func (t *template) Execute() (ValElementMap, error) {
	l := lex(t.tpl)
	l.run()
	parse := &Parser{
		lex:  l,
		tags: t.tags,
	}
	var err error
	root, err := parse.ParseDocument()
	if err != nil {
		return nil, err
	}
	err = root.Execute(t.ctx)
	if err != nil {
		return nil, err
	}
	return t.getValMap(ResultVal), nil
}
func (t *template) PublicValMap() ValElementMap {
	return t.getValMap(PublicVal)
}

func (t *template) ResultValMap() map[string]ValMap {
	return t.ctx.ResultMap
}

func (t *template) getValMap(valType ValType) ValElementMap {
	res := make(ValElementMap)
	for k, v := range t.ctx.ValMap {
		if v.ValType != valType {
			continue
		}
		res[k] = v
	}
	return res
}
