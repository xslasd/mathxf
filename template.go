package mathxf

import (
	"context"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
)

const DefResultKey = "res"
const DefResultEnvKey = "env"

var logger = log.New(os.Stdout, "[mathxf] ", log.LstdFlags|log.Lshortfile)
var debug bool

// Logging function (internally used)
func logf(format string, items ...any) {
	if debug {
		logger.Printf(format, items...)
	}
}

type template struct {
	tpl      string // the string being scanned
	ctx      *EvaluatorContext
	tags     map[string]TagParser
	keyOrder []string
	strMap   map[string]string

	root *nodeDocument
}

func Debug(b bool) {
	debug = b
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
func (t *template) SetParseErrFn(fn ParseECodeFn) {
	t.ctx.parseErrFn = fn
}
func (t *template) SetContext(ctx context.Context) {
	t.ctx.Context = ctx
}
func (t *template) SetDefResultKey(key string) {
	t.ctx.defResultKey = key
}
func (t *template) AddResultKeys(keys ...string) error {
	for _, key := range keys {
		if key == DefResultEnvKey {
			return t.ParseErr()(ResultKeyRegisteredErr.SetMessagef(key))
		}
		_, ok := t.ctx.ResultMap[key]
		if ok {
			return t.ParseErr()(ResultKeyRegisteredErr.SetMessagef(key))
		}
		t.ctx.ResultMap[key] = make(ValMap)
	}
	return nil
}
func (t *template) HighPrecision(b bool) {
	t.ctx.IsHighPrecision = b
}

// ReplaceStrMap strMap map[string]string ,map value must Ensure that the string contains at least one letter,
// while allowing numbers and underscores.For example: 'abc'、'abc123'、'abc_123'".
func (t *template) ReplaceStrMap(strMap map[string]string) error {
	for k, v := range strMap {
		if !containsAtLeastOneLetter(v) {
			return t.ParseErr()(InvalidReplaceStrErr.SetMessagef(v))
		}
		t.keyOrder = append(t.keyOrder, k)
	}
	sort.Slice(t.keyOrder, func(i, j int) bool {
		return len(t.keyOrder[i]) > len(t.keyOrder[j])
	})
	t.strMap = strMap
	return nil
}

func NewTemplate(tpl string) (*template, error) {
	t := &template{
		tpl:    tpl,
		strMap: make(map[string]string),
		tags:   defTags(),
		ctx: &EvaluatorContext{
			Context:         context.TODO(),
			IsHighPrecision: true,
			ValMap:          DefConst,
			ResultMap:       make(map[string]ValMap),
			defResultKey:    DefResultKey,
			parseErrFn:      ParseErr,
		},
	}
	t.ctx.ResultMap[t.ctx.defResultKey] = make(ValMap)
	return t, nil
}
func (t *template) ParseErr() ParseECodeFn {
	return t.ctx.parseErrFn
}
func (t *template) Execute(env map[string]any) (map[string]ValMap, error) {
	if t.root == nil {
		for _, k := range t.keyOrder {
			t.tpl = strings.ReplaceAll(t.tpl, k, t.strMap[k])
		}
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
		t.root = root
	} else {
		t.ctx.ResultMap = make(map[string]ValMap)
		t.ctx.ResultMap[t.ctx.defResultKey] = make(ValMap)
		for k, v := range t.ctx.ValMap {
			if v.ValType != PublicVal {
				continue
			}
			delete(t.ctx.ValMap, k)
		}
	}
	for k, v := range env {
		t.ctx.ValMap[k] = NewPublicValElement(v)
	}
	err := t.root.Execute(t.ctx)
	if err != nil {
		return nil, err
	}
	_env := make(ValMap)
	for k, ele := range t.ctx.ValMap {
		if ele.ValType != PublicVal {
			continue
		}
		if ele.IsSet {
			_env[k] = AsValue(ele.Val)
		}
	}
	if len(_env) > 0 {
		t.ctx.ResultMap[DefResultEnvKey] = _env
	}
	return t.ctx.ResultMap, nil
}
func (t *template) PublicValMap() ValElementMap {
	return t.getValMap(PublicVal)
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
