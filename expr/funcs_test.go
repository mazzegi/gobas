package expr

import (
	"fmt"
	"testing"
)

func TestScanArgs(t *testing.T) {
	vs := []interface{}{int(2), float64(63.78), "foo"}
	var a1 float64
	var a2 float64
	var a3 float64
	err := ScanArgs(vs, &a1, &a2, &a3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	fmt.Printf("[%f], [%f], [%f]\n", a1, a2, a3)
}
