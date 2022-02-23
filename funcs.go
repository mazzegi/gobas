package gobas

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/mazzegi/gobas/expr"
	"github.com/pkg/errors"
)

func BuiltinFuncs() *expr.Funcs {
	fs := expr.NewFuncs()

	fs.AddFloatFunc("ABS", func(vs []interface{}) (float64, error) {
		var a float64
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		return math.Abs(a), nil
	})
	fs.AddFunc("ASC", func(vs []interface{}) (interface{}, error) {
		var a string
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		if a == "" {
			return nil, errors.Errorf("ASC on empty string")
		}
		return int(a[0]), nil
	})
	fs.AddFloatFunc("ATN", func(vs []interface{}) (float64, error) {
		var a float64
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		return math.Atan(a), nil
	})
	fs.AddFunc("CHR$", func(vs []interface{}) (interface{}, error) {
		var a int
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		return string([]byte{byte(a)}), nil
	})
	fs.AddFloatFunc("COS", func(vs []interface{}) (float64, error) {
		var a float64
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		return math.Cos(a), nil
	})
	fs.AddFloatFunc("EXP", func(vs []interface{}) (float64, error) {
		var a float64
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		return math.Exp(a), nil
	})
	fs.AddFloatFunc("INT", func(vs []interface{}) (float64, error) {
		var a float64
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		return math.Round(a), nil
	})
	fs.AddFunc("LEFT$", func(vs []interface{}) (interface{}, error) {
		var s string
		var num int
		if err := expr.ScanArgs(vs, &s, &num); err != nil {
			return 0, err
		}
		if num >= len(s) {
			num = len(s) - 1
		}
		return s[:num], nil
	})
	fs.AddFunc("LEN", func(vs []interface{}) (interface{}, error) {
		var s string
		if err := expr.ScanArgs(vs, &s); err != nil {
			return 0, err
		}
		return len(s), nil
	})
	fs.AddFloatFunc("LOG", func(vs []interface{}) (float64, error) {
		var a float64
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		return math.Log(a), nil
	})
	fs.AddFunc("MID$", func(vs []interface{}) (interface{}, error) {
		var s string
		var idx, num int
		if err := expr.ScanArgs(vs, &s, &idx, &num); err != nil {
			return 0, err
		}
		toIdx := idx + num
		if toIdx >= len(s) {
			toIdx = len(s) - 1
		}
		return s[idx:toIdx], nil
	})
	fs.AddFloatFunc("RND", func(vs []interface{}) (float64, error) {
		return rand.Float64(), nil
	})
	fs.AddFunc("RIGHT$", func(vs []interface{}) (interface{}, error) {
		var s string
		var num int
		if err := expr.ScanArgs(vs, &s, &num); err != nil {
			return 0, err
		}
		fromIdx := len(s) - num
		if fromIdx < 0 {
			fromIdx = 0
		}

		return s[fromIdx:], nil
	})
	fs.AddFloatFunc("SGN", func(vs []interface{}) (float64, error) {
		var a float64
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		switch {
		case a == 0:
			return 0, nil
		case a < 0:
			return -1, nil
		default:
			return 1, nil
		}
	})
	fs.AddFloatFunc("SIN", func(vs []interface{}) (float64, error) {
		var a float64
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		return math.Sin(a), nil
	})
	fs.AddFloatFunc("SQR", func(vs []interface{}) (float64, error) {
		var a float64
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		return math.Sqrt(a), nil
	})
	fs.AddFunc("STR$", func(vs []interface{}) (interface{}, error) {
		var a float64
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		return fmt.Sprintf("%f", a), nil
	})
	fs.AddFunc("TAB", func(vs []interface{}) (interface{}, error) {
		var a int
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		return strings.Repeat(" ", a), nil
	})
	fs.AddFloatFunc("TAN", func(vs []interface{}) (float64, error) {
		var a float64
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		return math.Tan(a), nil
	})
	fs.AddFloatFunc("VAL", func(vs []interface{}) (float64, error) {
		var a string
		if err := expr.ScanArgs(vs, &a); err != nil {
			return 0, err
		}
		f, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return 0, err
		}
		return f, nil
	})

	return fs
}
