package expr

import (
	"constraints"
	"reflect"

	"github.com/pkg/errors"
)

var typeOfFloat64 = reflect.TypeOf(float64(0))
var typeOfString = reflect.TypeOf(string(""))

func ConvertToFloat(v interface{}) (float64, error) {
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

func convertToString(v interface{}) (string, error) {
	rv := reflect.ValueOf(v)
	if !rv.CanConvert(typeOfString) {
		return "", errors.Errorf("cannot convert value of type %T to type %s", v, typeOfString.String())
	}
	return rv.Convert(typeOfString).Interface().(string), nil
}

func canConvertToString(v interface{}) bool {
	rv := reflect.ValueOf(v)
	return rv.CanConvert(typeOfString)
}

//
func lessFloat[T constraints.Ordered](t1, t2 T) float64 {
	if t1 < t2 {
		return 1
	}
	return 0
}

func greaterFloat[T constraints.Ordered](t1, t2 T) float64 {
	if t1 > t2 {
		return 1
	}
	return 0
}

func equalFloat[T constraints.Ordered](t1, t2 T) float64 {
	if t1 == t2 {
		return 1
	}
	return 0
}

func notEqualFloat[T constraints.Ordered](t1, t2 T) float64 {
	if t1 != t2 {
		return 1
	}
	return 0
}

func lessEqualFloat[T constraints.Ordered](t1, t2 T) float64 {
	if t1 <= t2 {
		return 1
	}
	return 0
}

func greaterEqualFloat[T constraints.Ordered](t1, t2 T) float64 {
	if t1 >= t2 {
		return 1
	}
	return 0
}

//
func floatOp(v1, v2 interface{}, op func(f1, f2 float64) float64) (interface{}, error) {
	f1, err := ConvertToFloat(v1)
	if err != nil {
		return nil, err
	}
	f2, err := ConvertToFloat(v2)
	if err != nil {
		return nil, err
	}
	return op(f1, f2), nil
}

func stringOp(v1, v2 interface{}, op func(f1, f2 string) float64) (interface{}, error) {
	f1, err := convertToString(v1)
	if err != nil {
		return nil, err
	}
	f2, err := convertToString(v2)
	if err != nil {
		return nil, err
	}
	return op(f1, f2), nil
}

func less(v1, v2 interface{}) (interface{}, error) {
	if canConvertToFloat(v1) {
		return floatOp(v1, v2, lessFloat[float64])
	} else if canConvertToString(v1) {
		return stringOp(v1, v2, lessFloat[string])
	}
	return nil, errors.Errorf("incompatible types %T, %T", v1, v2)
}

func greater(v1, v2 interface{}) (interface{}, error) {
	if canConvertToFloat(v1) {
		return floatOp(v1, v2, greaterFloat[float64])
	} else if canConvertToString(v1) {
		return stringOp(v1, v2, greaterFloat[string])
	}
	return nil, errors.Errorf("incompatible types %T, %T", v1, v2)
}

func equal(v1, v2 interface{}) (interface{}, error) {
	if canConvertToFloat(v1) {
		return floatOp(v1, v2, equalFloat[float64])
	} else if canConvertToString(v1) {
		return stringOp(v1, v2, equalFloat[string])
	}
	return nil, errors.Errorf("incompatible types %T, %T", v1, v2)
}

func notEqual(v1, v2 interface{}) (interface{}, error) {
	if canConvertToFloat(v1) {
		return floatOp(v1, v2, notEqualFloat[float64])
	} else if canConvertToString(v1) {
		return stringOp(v1, v2, notEqualFloat[string])
	}
	return nil, errors.Errorf("incompatible types %T, %T", v1, v2)
}

func lessEqual(v1, v2 interface{}) (interface{}, error) {
	if canConvertToFloat(v1) {
		return floatOp(v1, v2, lessEqualFloat[float64])
	} else if canConvertToString(v1) {
		return stringOp(v1, v2, lessEqualFloat[string])
	}
	return nil, errors.Errorf("incompatible types %T, %T", v1, v2)
}

func greaterEqual(v1, v2 interface{}) (interface{}, error) {
	if canConvertToFloat(v1) {
		return floatOp(v1, v2, greaterEqualFloat[float64])
	} else if canConvertToString(v1) {
		return stringOp(v1, v2, greaterEqualFloat[string])
	}
	return nil, errors.Errorf("incompatible types %T, %T", v1, v2)
}
