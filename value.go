package mathxf

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

var (
	TypeOfValuePtr         = reflect.TypeOf(new(Value))
	TypeOfValMapPtr        = reflect.TypeOf(make(ValMap))
	TypeOfValElementMapPrt = reflect.TypeOf(new(ValElementMap))
	TypeOfValElementPrt    = reflect.TypeOf(new(ValElement))
	TypeOfEvaluatorContext = reflect.TypeOf(new(EvaluatorContext))
	TypeOfDecimalPtr       = reflect.TypeOf(new(decimal.Decimal))
)

type Value struct {
	Val reflect.Value
}

// AsValue converts any given value to a pongo2.Value
// Usually being used within own functions passed to a template
// through a Context or within filter functions.
//
// Example:
//
//	AsValue("my string")
func AsValue(i any) *Value {
	return &Value{
		Val: reflect.ValueOf(i),
	}
}
func (v *Value) getResolvedValue() reflect.Value {
	if v.Val.IsValid() && v.Val.Kind() == reflect.Ptr {
		return v.Val.Elem()
	}
	return v.Val
}

// IsString checks whether the underlying value is a string
func (v *Value) IsString() bool {
	return v.getResolvedValue().Kind() == reflect.String
}

// IsBool checks whether the underlying value is a bool
func (v *Value) IsBool() bool {
	return v.getResolvedValue().Kind() == reflect.Bool
}

// IsFloat checks whether the underlying value is a float
func (v *Value) IsFloat() bool {
	val := v.getResolvedValue()
	return val.Kind() == reflect.Float32 ||
		val.Kind() == reflect.Float64
}
func (v *Value) IsDecimal() bool {
	val := v.getResolvedValue()
	if !val.IsValid() {
		return false
	}
	return val.Type() == TypeOfDecimalPtr.Elem() || v.IsNumber()
}
func (v *Value) Decimal() decimal.Decimal {
	val := v.getResolvedValue()
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return decimal.NewFromInt(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return decimal.NewFromInt(int64(val.Uint()))
	case reflect.Float32, reflect.Float64:
		return decimal.NewFromFloat(val.Float())
	case reflect.String:
		// Try to convert from string to float64 (base 10)
		f, err := strconv.ParseFloat(v.getResolvedValue().String(), 64)
		if err != nil {
			return decimal.Decimal{}
		}
		return decimal.NewFromFloat(f)
	default:
		if val.IsValid() && val.Type() == TypeOfDecimalPtr.Elem() {
			b, ok := val.Interface().(decimal.Decimal)
			if ok {
				return b
			}
		}
		logf("Value.Float() not available for type: %s\n", v.getResolvedValue().Kind().String())
		return decimal.Decimal{}
	}
}

// IsInteger checks whether the underlying value is an integer
func (v *Value) IsInteger() bool {
	val := v.getResolvedValue()
	return val.Kind() == reflect.Int ||
		val.Kind() == reflect.Int8 ||
		val.Kind() == reflect.Int16 ||
		val.Kind() == reflect.Int32 ||
		val.Kind() == reflect.Int64 ||
		val.Kind() == reflect.Uint ||
		val.Kind() == reflect.Uint8 ||
		val.Kind() == reflect.Uint16 ||
		val.Kind() == reflect.Uint32 ||
		val.Kind() == reflect.Uint64 ||
		val.Type() == TypeOfDecimalPtr.Elem()
}

// IsNumber checks whether the underlying value is either an integer
// or a float.
func (v *Value) IsNumber() bool {
	return v.IsInteger() || v.IsFloat()
}

// IsTime checks whether the underlying value is a time.Time.
func (v *Value) IsTime() bool {
	_, ok := v.Interface().(time.Time)
	return ok
}

// IsNil checks whether the underlying value is NIL
func (v *Value) IsNil() bool {
	// fmt.Printf("%+v\n", v.getResolvedValue().Type().String())
	return !v.getResolvedValue().IsValid()
}

// String returns a string for the underlying value. If this value is not
// of type string, pongo2 tries to convert it. Currently the following
// types for underlying values are supported:
//
//  1. string
//  2. int/uint (any size)
//  3. float (any precision)
//  4. bool
//  5. time.Time
//  6. String() will be called on the underlying value if provided
//
// NIL values will lead to an empty string. Unsupported types are leading
// to their respective type name.
func (v *Value) String() string {
	if v.IsNil() {
		return ""
	}

	if t, ok := v.Interface().(fmt.Stringer); ok {
		return t.String()
	}

	switch v.getResolvedValue().Kind() {
	case reflect.String:
		return v.getResolvedValue().String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.getResolvedValue().Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.getResolvedValue().Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", v.getResolvedValue().Float())
	case reflect.Bool:
		if v.Bool() {
			return "True"
		}
		return "False"
	}

	logf("Value.String() not implemented for type: %s\n", v.getResolvedValue().Kind().String())
	return v.getResolvedValue().String()
}

// Integer returns the underlying value as an integer (converts the underlying
// value, if necessary). If it'name not possible to convert the underlying value,
// it will return 0.
func (v *Value) Integer() int {
	val := v.getResolvedValue()
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(val.Uint())
	case reflect.Float32, reflect.Float64:
		return int(val.Float())
	case reflect.String:
		// Try to convert from string to int (base 10)
		f, err := strconv.ParseFloat(val.String(), 64)
		if err != nil {
			return 0
		}
		return int(f)
	default:
		if val.Type() == TypeOfDecimalPtr.Elem() {
			b, ok := val.Interface().(decimal.Decimal)
			if ok {
				f, _ := b.Float64()
				return int(f)
			}
			return 0
		}
		logf("Value.Integer() not available for type: %s\n", v.getResolvedValue().Kind().String())
		return 0
	}
}

// Float returns the underlying value as a float (converts the underlying
// value, if necessary). If it'name not possible to convert the underlying value,
// it will return 0.0.
func (v *Value) Float() float64 {
	val := v.getResolvedValue()
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(val.Uint())
	case reflect.Float32, reflect.Float64:
		return val.Float()
	case reflect.String:
		// Try to convert from string to float64 (base 10)
		f, err := strconv.ParseFloat(val.String(), 64)
		if err != nil {
			return 0.0
		}
		return f
	default:
		if val.Type() == TypeOfDecimalPtr.Elem() {
			b, ok := val.Interface().(decimal.Decimal)
			if ok {
				f, _ := b.Float64()
				return f
			}
			return 0
		}
		logf("Value.Float() not available for type: %s\n", v.getResolvedValue().Kind().String())
		return 0.0
	}
}

// Bool returns the underlying value as bool. If the value is not bool, false
// will always be returned. If you're looking for true/false-evaluation of the
// underlying value, have a look on the IsTrue()-function.
func (v *Value) Bool() bool {
	switch v.getResolvedValue().Kind() {
	case reflect.Bool:
		return v.getResolvedValue().Bool()
	default:
		logf("Value.Bool() not available for type: %s\n", v.getResolvedValue().Kind().String())
		return false
	}
}

// Time returns the underlying value as time.Time.
// If the underlying value is not a time.Time, it returns the zero value of time.Time.
func (v *Value) Time() time.Time {
	tm, ok := v.Interface().(time.Time)
	if ok {
		return tm
	}
	return time.Time{}
}

// IsTrue tries to evaluate the underlying value the Pythonic-way:
//
// Returns TRUE in one the following cases:
//
//   - int != 0
//   - uint != 0
//   - float != 0.0
//   - len(array/chan/map/slice/string) > 0
//   - bool == true
//   - underlying value is a struct
//
// Otherwise returns always FALSE.
func (v *Value) IsTrue() bool {
	val := v.getResolvedValue()
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return val.Float() != 0
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return val.Len() > 0
	case reflect.Bool:
		return val.Bool()
	case reflect.Struct:
		return true // struct instance is always true
	default:
		if val.Type() == TypeOfDecimalPtr.Elem() {
			b, ok := val.Interface().(decimal.Decimal)
			if ok {
				return b.Cmp(decimal.Zero) > 0
			}
			return false
		}
		logf("Value.IsTrue() not available for type: %s\n", v.getResolvedValue().Kind().String())
		return false
	}
}

// Len returns the length for an array, chan, map, slice or string.
// Otherwise it will return 0.
func (v *Value) Len() int {
	switch v.getResolvedValue().Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return v.getResolvedValue().Len()
	case reflect.String:
		runes := []rune(v.getResolvedValue().String())
		return len(runes)
	default:
		logf("Value.Len() not available for type: %s\n", v.getResolvedValue().Kind().String())
		return 0
	}
}

// Slice slices an array, slice or string. Otherwise it will
// return an empty []int.
func (v *Value) Slice(i, j int) *Value {
	switch v.getResolvedValue().Kind() {
	case reflect.Array, reflect.Slice:
		return AsValue(v.getResolvedValue().Slice(i, j).Interface())
	case reflect.String:
		runes := []rune(v.getResolvedValue().String())
		return AsValue(string(runes[i:j]))
	default:
		logf("Value.Slice() not available for type: %s\n", v.getResolvedValue().Kind().String())
		return AsValue([]int{})
	}
}

// Index gets the i-th item of an array, slice or string. Otherwise
// it will return NIL.
func (v *Value) Index(i int) *Value {
	switch v.getResolvedValue().Kind() {
	case reflect.Array, reflect.Slice:
		if i >= v.Len() {
			return AsValue(nil)
		}
		return AsValue(v.getResolvedValue().Index(i).Interface())
	case reflect.String:
		s := v.getResolvedValue().String()
		runes := []rune(s)
		if i < len(runes) {
			return AsValue(string(runes[i]))
		}
		return AsValue("")
	default:
		logf("Value.Slice() not available for type: %s\n", v.getResolvedValue().Kind().String())
		return AsValue([]int{})
	}
}

// Contains checks whether the underlying value (which must be of type struct, map,
// string, array or slice) contains of another Value (e. g. used to check
// whether a struct contains of a specific field or a map contains a specific key).
//
// Example:
//
//	AsValue("Hello, World!").Contains(AsValue("World")) == true
func (v *Value) Contains(other *Value) bool {
	baseValue := v.getResolvedValue()
	switch baseValue.Kind() {
	case reflect.Struct:
		fieldValue := baseValue.FieldByName(other.String())
		return fieldValue.IsValid()
	case reflect.Map:
		// We can't check against invalid types
		if !other.Val.IsValid() {
			return false
		}
		// Ensure that map key type is equal to other'name type.
		if baseValue.Type().Key() != other.Val.Type() {
			return false
		}

		var mapValue reflect.Value
		switch other.Interface().(type) {
		case int:
			mapValue = baseValue.MapIndex(other.getResolvedValue())
		case string:
			mapValue = baseValue.MapIndex(other.getResolvedValue())
		default:
			logf("Value.Contains() does not support lookup type '%s'\n", other.getResolvedValue().Kind().String())
			return false
		}
		return mapValue.IsValid()
	case reflect.String:
		return strings.Contains(baseValue.String(), other.String())

	case reflect.Slice, reflect.Array:
		for i := 0; i < baseValue.Len(); i++ {
			item := baseValue.Index(i)
			if item.Type() == TypeOfValuePtr {
				tmpValue := item.Interface().(*Value)
				item = tmpValue.Val
			}
			if other.EqualValueTo(AsValue(item.Interface())) {
				return true
			}
		}
		return false
	default:
		logf("Value.Contains() not available for type: %s\n", baseValue.Kind().String())
		return false
	}
}

// CanSlice checks whether the underlying value is of type array, slice or string.
// You normally would use CanSlice() before using the Slice() operation.
func (v *Value) CanSlice() bool {
	switch v.getResolvedValue().Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		return true
	}
	return false
}

// Iterate iterates over a map, array, slice or a string. It calls the
// function'name first argument for every value with the following arguments:
//
//	idx      current 0-index
//	count    total amount of items
//	key      *Value for the key or item
//	value    *Value (only for maps, the respective value for a specific key)
//
// If the underlying value has no items or is not one of the types above,
// the empty function (function'name second argument) will be called.
func (v *Value) Iterate(fn func(idx, count int, key, value *Value) bool, empty func()) {
	v.IterateOrder(fn, empty, false, false)
}

// IterateOrder behaves like Value.Iterate, but can iterate through an array/slice/string in reverse. Does
// not affect the iteration through a map because maps don't have any particular order.
// However, you can force an order using the `sorted` keyword (and even use `reversed sorted`).
func (v *Value) IterateOrder(fn func(idx, count int, key, value *Value) bool, empty func(), reverse bool, sorted bool) {
	switch v.getResolvedValue().Kind() {
	case reflect.Map:
		keys := sortedKeys(v.getResolvedValue().MapKeys())
		if sorted {
			if reverse {
				sort.Sort(sort.Reverse(keys))
			} else {
				sort.Sort(keys)
			}
		}
		keyLen := len(keys)
		for idx, key := range keys {
			value := v.getResolvedValue().MapIndex(key)
			if !fn(idx, keyLen, &Value{Val: key}, &Value{Val: value}) {
				return
			}
		}
		if keyLen == 0 {
			empty()
		}
		return // done
	case reflect.Array, reflect.Slice:
		var items valuesList

		itemCount := v.getResolvedValue().Len()
		for i := 0; i < itemCount; i++ {
			items = append(items, &Value{Val: v.getResolvedValue().Index(i)})
		}

		if sorted {
			if reverse {
				sort.Sort(sort.Reverse(items))
			} else {
				sort.Sort(items)
			}
		} else {
			if reverse {
				for i := 0; i < itemCount/2; i++ {
					items[i], items[itemCount-1-i] = items[itemCount-1-i], items[i]
				}
			}
		}

		if len(items) > 0 {
			for idx, item := range items {
				if !fn(idx, itemCount, item, nil) {
					return
				}
			}
		} else {
			empty()
		}
		return // done
	case reflect.String:
		s := v.getResolvedValue().String()
		rs := []rune(s)
		charCount := len(rs)

		if charCount > 0 {
			if sorted {
				sort.SliceStable(rs, func(i, j int) bool {
					return rs[i] < rs[j]
				})
			}

			if reverse {
				for i, j := 0, charCount-1; i < j; i, j = i+1, j-1 {
					rs[i], rs[j] = rs[j], rs[i]
				}
			}

			for i := 0; i < charCount; i++ {
				if !fn(i, charCount, &Value{Val: reflect.ValueOf(string(rs[i]))}, nil) {
					return
				}
			}
		} else {
			empty()
		}
		return // done
	default:
		logf("Value.Iterate() not available for type: %s\n", v.getResolvedValue().Kind().String())
	}
	empty()
}

// Interface gives you access to the underlying value.
func (v *Value) Interface() any {
	if v.Val.IsValid() {
		return v.Val.Interface()
	}
	return nil
}

// EqualValueTo checks whether two values are containing the same value or object (if comparable).
func (v *Value) EqualValueTo(other *Value) bool {
	if v.IsInteger() && other.IsInteger() {
		return v.Integer() == other.Integer()
	}
	if v.IsTime() && other.IsTime() {
		return v.Time().Equal(other.Time())
	}
	if !v.Val.IsValid() || !other.Val.IsValid() {
		return false
	}
	if v.IsFloat() && other.IsFloat() {
		return v.Float() == other.Float()
	}
	return v.Val.Equal(other.Val)
}

type sortedKeys []reflect.Value

func (sk sortedKeys) Len() int {
	return len(sk)
}

func (sk sortedKeys) Less(i, j int) bool {
	vi := &Value{Val: sk[i]}
	vj := &Value{Val: sk[j]}
	switch {
	case vi.IsInteger() && vj.IsInteger():
		return vi.Integer() < vj.Integer()
	case vi.IsFloat() && vj.IsFloat():
		return vi.Float() < vj.Float()
	default:
		return vi.String() < vj.String()
	}
}

func (sk sortedKeys) Swap(i, j int) {
	sk[i], sk[j] = sk[j], sk[i]
}

type valuesList []*Value

func (vl valuesList) Len() int {
	return len(vl)
}

func (vl valuesList) Less(i, j int) bool {
	vi := vl[i]
	vj := vl[j]
	switch {
	case vi.IsInteger() && vj.IsInteger():
		return vi.Integer() < vj.Integer()
	case vi.IsFloat() && vj.IsFloat():
		return vi.Float() < vj.Float()
	default:
		return vi.String() < vj.String()
	}
}

func (vl valuesList) Swap(i, j int) {
	vl[i], vl[j] = vl[j], vl[i]
}
