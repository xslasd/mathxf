package mathxf

import "reflect"

type ResValues map[string]reflect.Value

func (r ResValues) Float(key string) float64 {
	if _, ok := r[key]; !ok {
		return 0.0
	}
	return r[key].Float()
}
