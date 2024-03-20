package mathxf

import (
	"fmt"
	"github.com/shopspring/decimal"
	"reflect"
	"strings"
)

type assignmentResolver struct {
	variable *variableResolver
	value    IEvaluator
}

func (a assignmentResolver) Execute(ctx EvaluatorContext) error {
	val, err := a.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Println("--------assignmentResolver---------Execute---------", val)
	return nil
}

func (a assignmentResolver) GetPositionToken() *Token {
	return a.variable.GetPositionToken()
}

func (a assignmentResolver) Evaluate(ctx EvaluatorContext) (*Value, error) {
	err := a.variable.SetPartValue(ctx, a.value.Evaluate)
	if err != nil {
		return nil, err
	}
	return AsValue(nil), nil
}

type numberResolver struct {
	locationToken *Token
	val           float64
}

func (f numberResolver) Execute(ctx EvaluatorContext) error {
	val, err := f.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Println(val)
	return nil
}

func (f numberResolver) GetPositionToken() *Token {
	return f.locationToken
}

func (f numberResolver) Evaluate(ctx EvaluatorContext) (*Value, error) {
	if ctx.IsHighPrecision {
		return AsValue(decimal.NewFromFloat(f.val)), nil
	}
	return AsValue(f.val), nil
}

type boolResolver struct {
	locationToken *Token
	val           bool
}

func (b boolResolver) Execute(ctx EvaluatorContext) error {
	val, err := b.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Println(val)
	return nil
}

func (b boolResolver) GetPositionToken() *Token {
	return b.locationToken
}

func (b boolResolver) Evaluate(ctx EvaluatorContext) (*Value, error) {
	return AsValue(b.val), nil
}

type stringResolver struct {
	locationToken *Token
	val           string
}

func (s stringResolver) Execute(ctx EvaluatorContext) error {
	val, err := s.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Println(val)
	return nil
}

func (s stringResolver) GetPositionToken() *Token {
	return s.locationToken
}

func (s stringResolver) Evaluate(ctx EvaluatorContext) (*Value, error) {
	return AsValue(s.val), nil
}

type variableResolver struct {
	locationToken *Token
	parts         []*variablePart
}

func (v variableResolver) Execute(ctx EvaluatorContext) error {
	val, err := v.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Println(val)
	return nil
}

func (v variableResolver) GetPositionToken() *Token {
	return v.locationToken
}
func (v variableResolver) String() string {
	parts := make([]string, 0, len(v.parts))
	for _, p := range v.parts {
		parts = append(parts, p.String())
	}

	return strings.Join(parts, ".")
}
func (v variableResolver) SetPartValue(ctx EvaluatorContext, valueFunc func(ctx EvaluatorContext) (*Value, error)) error {
	var varData reflect.Value
	var keyName string
	var keyInd int
	pLen := len(v.parts)
	for index, part := range v.parts {
		if part.isFunctionCall {
			return fmt.Errorf("can't set value of function call")
		}
		keyName = part.name
		if index == 0 {
			if _, ok := ctx.ResValues[keyName]; ok {
				varData = reflect.ValueOf(&ctx.ResValues).Elem()
				varData = varData.MapIndex(reflect.ValueOf(keyName))
			} else {
				varData = reflect.ValueOf(&ctx.Private).Elem()
				val := varData.MapIndex(reflect.ValueOf(keyName))
				if !val.IsValid() {
					varData = reflect.ValueOf(&ctx.Public).Elem()
					val = varData.MapIndex(reflect.ValueOf(keyName))
					if !val.IsValid() {
						return fmt.Errorf("invalid variable %s", v.String())
					}
				}
			}
		} else {
			if varData.Kind() == reflect.Ptr {
				varData = varData.Elem()
				if !varData.IsValid() {
					return fmt.Errorf("invalid value %v", v.String())
				}
			}
			switch part.typ {
			case VariablePartTypeIdent:
				switch varData.Kind() {
				case reflect.Struct:
					if index != pLen-1 {
						varData = varData.FieldByName(part.name)
					}
				case reflect.Map:
					if index != pLen-1 {
						varData = varData.MapIndex(reflect.ValueOf(part.name))
					}
				default:
					return fmt.Errorf("can't access a field by name on type %s (variable %s)",
						varData.Kind().String(), v.String())
				}
			case VariablePartTypeSubscript:
				switch varData.Kind() {
				case reflect.String, reflect.Array, reflect.Slice:
					eVal, err := part.subscript.Evaluate(ctx)
					if err != nil {
						return err
					}
					ind := eVal.Integer()
					keyInd = ind
					if ind >= 0 && varData.Len() > ind && index != pLen-1 {
						varData = varData.Index(ind)
					} else {
						return fmt.Errorf("index out of bounds %s: 0-%d (index %d)", part.name, varData.Len(), ind)
					}
				case reflect.Struct:
					eVal, err := part.subscript.Evaluate(ctx)
					if err != nil {
						return err
					}
					keyName = eVal.String()
					if index != pLen-1 {
						varData = varData.FieldByName(keyName)
					}
				case reflect.Map:
					eVal, err := part.subscript.Evaluate(ctx)
					if err != nil {
						return err
					}
					if eVal.IsNil() {
						return fmt.Errorf("invalid value %v", v.String())
					}
					if !eVal.val.Type().AssignableTo(varData.Type().Key()) {
						return fmt.Errorf("invalid key type %v", v.String())
					}
					keyName = eVal.String()
					if index != pLen-1 {
						varData = varData.MapIndex(eVal.val)
					}
				default:
					return fmt.Errorf("can't access a field by index on type %s (variable %s)",
						varData.Kind().String(), v.String())
				}
			default:
				return fmt.Errorf("unimplemented")
			}

		}
	}
	if !varData.IsValid() {
		return fmt.Errorf("invalid variable %s", v.String())
	}

	val, err := valueFunc(ctx)
	if err != nil {
		return err
	}
	switch varData.Kind() {
	case reflect.Struct:
		varData.FieldByName(keyName).Set(reflect.ValueOf(val.Interface()))
	case reflect.Map:
		varData.SetMapIndex(reflect.ValueOf(keyName), reflect.ValueOf(val.val))
	case reflect.String:
		varData.SetString(val.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		varData.SetInt(int64(val.Integer()))
	case reflect.Array:
		varData.Index(keyInd).Set(reflect.ValueOf(val.Interface()))
	case reflect.Slice:
		varData.Index(keyInd).Set(reflect.ValueOf(val.Interface()))
	case reflect.Float32, reflect.Float64:
		if !varData.CanSet() {
			return fmt.Errorf("can't set value of %s", v.String())
		}
		varData.SetFloat(val.Float())
	case reflect.Bool:
		if !varData.CanSet() {
			return fmt.Errorf("can't set value of %s", v.String())
		}
		varData.SetBool(val.IsTrue())
	default:
		return fmt.Errorf("can't set value of %s", v.String())
	}
	return nil
}

func (v variableResolver) Evaluate(ctx EvaluatorContext) (*Value, error) {
	var varData reflect.Value
	fmt.Println("--------varData---------", v.String())
	for index, part := range v.parts {
		if index == 0 {
			var ok bool
			name := part.name
			varData, ok = ctx.Private[name]
			if !ok {
				varData = ctx.Public[name]
			}
		} else {
			isFunc := false
			if part.typ == VariablePartTypeIdent {
				funcValue := varData.MethodByName(part.name)
				if funcValue.IsValid() {
					varData = funcValue
					isFunc = true
				}
			}
			if !isFunc {
				if varData.Kind() == reflect.Ptr {
					varData = varData.Elem()
					if !varData.IsValid() {
						// Value is not valid (anymore)
						return AsValue(nil), nil
					}
				}
				switch part.typ {
				case VariablePartTypeIdent:
					switch varData.Kind() {
					case reflect.Struct:
						varData = varData.FieldByName(part.name)
					case reflect.Map:
						varData = varData.MapIndex(reflect.ValueOf(part.name))
					default:
						return nil, fmt.Errorf("can't access a field by name on type %s (variable %s)",
							varData.Kind().String(), v.String())
					}
				case VariablePartTypeSubscript:
					switch varData.Kind() {
					case reflect.String, reflect.Array, reflect.Slice:
						eVal, err := part.subscript.Evaluate(ctx)
						if err != nil {
							return nil, err
						}
						ind := eVal.Integer()
						if ind >= 0 && varData.Len() > ind {
							varData = varData.Index(ind)
						} else {
							return nil, fmt.Errorf("index out of bounds %s: 0-%d (index %d)", part.name, varData.Len(), ind)
						}
					case reflect.Struct:
						eVal, err := part.subscript.Evaluate(ctx)
						if err != nil {
							return nil, err
						}
						varData = varData.FieldByName(eVal.String())
					case reflect.Map:
						eVal, err := part.subscript.Evaluate(ctx)
						if err != nil {
							return nil, err
						}
						if eVal.IsNil() {
							return AsValue(nil), nil
						}
						if eVal.val.Type().AssignableTo(varData.Type().Key()) {
							varData = varData.MapIndex(eVal.val)
						} else {
							return AsValue(nil), nil
						}
					default:
						return nil, fmt.Errorf("can't access a field by index on type %s (variable %s)",
							varData.Kind().String(), v.String())
					}
				default:
					return nil, fmt.Errorf("unimplemented")
				}
			}
		}
		if !varData.IsValid() {
			// Value is not valid (anymore)
			return AsValue(nil), nil
		}
		if varData.Type() == TypeOfValuePtr {
			tmpValue := varData.Interface().(*Value)
			varData = tmpValue.val
		}
		if varData.Kind() == reflect.Interface {
			varData = reflect.ValueOf(varData.Interface())
		}
		if part.isFunctionCall {
			if varData.Kind() != reflect.Func {
				return nil, fmt.Errorf("variable %s is not a function", v.String())
			}
			funcT := varData.Type()
			numIn := funcT.NumIn()
			numOut := funcT.NumOut()
			isVariadic := funcT.IsVariadic()

			currArgs := part.callingArgs
			if len(currArgs) != numIn && !(len(currArgs) >= numIn-1 && isVariadic) {
				return nil, fmt.Errorf("function %s requires %d arguments, but got %d", v.String(), numIn, len(currArgs))
			}
			if numOut < 1 && numOut > 2 {
				return nil, fmt.Errorf("'%s' must have exactly 1 or 2 output arguments, the second argument must be of type error", v.String())
			}
			var args []reflect.Value

			ind := 0
			if numIn > 0 && funcT.In(0) == TypeOfEvaluatorContext {
				args = append(args, reflect.ValueOf(ctx))
				ind = 1
			}
			for i, arg := range currArgs {
				pv, err := arg.Evaluate(ctx)
				if err != nil {
					return nil, err
				}
				var fnArg reflect.Type
				inds := ind + 1
				if isVariadic && inds >= numIn-1 {
					fnArg = funcT.In(numIn - 1).Elem()
				} else {
					fnArg = funcT.In(inds)
				}
				fmt.Println(fnArg, "-------fnArg----------")
				if fnArg != TypeOfValuePtr {
					if !isVariadic {
						if fnArg != reflect.TypeOf(pv.Interface()) && fnArg.Kind() != reflect.Interface {
							return nil, fmt.Errorf("function input argument %d of '%s' must be of type %s or *mathxf.Value (not %T)",
								i, v.String(), fnArg.String(), pv.Interface())
						}
					} else {
						if fnArg != reflect.TypeOf(pv.Interface()) && fnArg.Kind() != reflect.Interface {
							return nil, fmt.Errorf("function variadic input argument of '%s' must be of type %s or *mathxf.Value (not %T)",
								v.String(), fnArg.String(), pv.Interface())
						}
					}
					if pv.IsNil() {
						// Workaround to present an interface nil as reflect.Value
						var empty any = nil
						args = append(args, reflect.ValueOf(&empty).Elem())
					} else {
						args = append(args, reflect.ValueOf(pv.Interface()))
					}
				} else {
					args = append(args, reflect.ValueOf(pv))
				}
			}
			for _, arg1 := range args {
				if arg1.Kind() == reflect.Invalid {
					return nil, fmt.Errorf("invalid argument")
				}
			}
			results := varData.Call(args)
			rVal := results[0]
			if numOut == 2 {
				errVal := results[1].Interface()
				if errVal != nil {
					return nil, errVal.(error)
				}
			}
			if rVal.Type() != TypeOfValuePtr {
				varData = reflect.ValueOf(rVal.Interface())
			} else {
				varData = rVal.Interface().(*Value).val
			}
		}
		if !varData.IsValid() {
			// Value is not valid (e. g. NIL value)
			return AsValue(nil), nil
		}
	}
	return &Value{val: varData}, nil
}

type VariablePartType int

const (
	VariablePartTypeArray VariablePartType = iota
	VariablePartTypeIdent
	VariablePartTypeSubscript
)

type variablePart struct {
	typ       VariablePartType
	name      string
	subscript IEvaluator

	isFunctionCall bool
	callingArgs    []functionCallArgument // needed for a function call, represents all argument nodes (INode supports nested function calls)
}

func (v variablePart) String() string {
	switch v.typ {
	case VariablePartTypeIdent:
		return v.name
	case VariablePartTypeSubscript:
		return "[subscript]"
	case VariablePartTypeArray:
		return "[array]"
	}
	panic("unimplemented")
}

type arrayResolver struct {
	locationToken *Token
	parts         []*variablePart
}

func (a arrayResolver) Execute(ctx EvaluatorContext) error {
	val, err := a.Evaluate(ctx)
	if err != nil {
		return err
	}
	fmt.Println(val)
	return nil
}

func (a arrayResolver) GetPositionToken() *Token {
	return a.locationToken
}

func (a arrayResolver) Evaluate(ctx EvaluatorContext) (*Value, error) {
	if len(a.parts) == 0 {
		return &Value{}, nil
	}
	items := make([]*Value, 0)
	for _, part := range a.parts {
		item, err := part.subscript.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &Value{
		val: reflect.ValueOf(items),
	}, nil
}

func (a arrayResolver) String() string {
	parts := make([]string, 0, len(a.parts))
	for _, p := range a.parts {
		parts = append(parts, p.String())
	}

	return strings.Join(parts, ".")
}
