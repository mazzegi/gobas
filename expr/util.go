package expr

import (
	"reflect"

	"github.com/pkg/errors"
)

var typeOfFloat64 = reflect.TypeOf(float64(0))

func convertToFloat(v interface{}) (float64, error) {
	rv := reflect.ValueOf(v)
	if !rv.CanConvert(typeOfFloat64) {
		return 0, errors.Errorf("cannot convert value of type %T to type %s", v, typeOfFloat64.String())
	}
	return rv.Convert(typeOfFloat64).Interface().(float64), nil
}

func canConvertToFloat(v interface{}) bool {
	rv := reflect.ValueOf(v)
	return rv.CanConvert(typeOfFloat64)
}
