package gobas

import "testing"

type set struct {
	cs  []int
	val int
}

func TestArray(t *testing.T) {
	dims := []int{10, 20, 30}
	a := NewArray[int](dims)

	cs := []int{0, 0, 0}
	inc := func() bool {
		for i := 0; i < len(cs); i++ {
			if cs[i] < dims[i]-1 {
				cs[i]++
				return true
			}
			cs[i] = 0
		}
		return false
	}
	hash := func(cs []int) int {
		var h int
		for i, c := range cs {
			h += (i + 1) * c
		}
		return 23 + h%22
	}

	for {
		ok := inc()
		if !ok {
			break
		}
		err := a.Set(cs, hash(cs))
		if err != nil {
			t.Fatalf("set-error: %v", err)
		}
	}

	cs = []int{0, 0, 0}
	for {
		ok := inc()
		if !ok {
			break
		}
		v, err := a.Get(cs)
		if err != nil {
			t.Fatalf("get-error: %v", err)
		}
		h := hash(cs)
		if h != v {
			t.Fatalf("expect %d, got %d", h, v)
		}
	}
}
