package gobas

import "github.com/pkg/errors"

var idxOf1 = false

func NewArray[T any](dims []int) *Array[T] {
	if !idxOf1 {
		for i := 0; i < len(dims); i++ {
			dims[i]++
		}
	}

	size := 1
	for _, dim := range dims {
		size *= dim
	}
	a := &Array[T]{
		dims: dims,
		data: make([]T, size),
	}
	return a
}

type Array[T any] struct {
	dims []int
	data []T
}

func (a *Array[T]) segmentSize(dim int) int {
	if dim < 0 || dim >= len(a.dims) {
		return 0
	}
	size := 1
	for i := dim + 1; i < len(a.dims); i++ {
		size *= a.dims[i]
	}
	return size
}

func (a *Array[T]) idx(cs []int) (int, error) {
	if len(cs) != len(a.dims) {
		return 0, errors.Errorf("invalid argument count %d, want %d", len(cs), len(a.dims))
	}

	ix := 0
	for i, c := range cs {
		//BAS indexes start with 1
		if idxOf1 {
			c = c - 1
		}
		dim := a.dims[i]
		if c < 0 || c >= dim {
			return 0, errors.Errorf("index of dim %d (%d) out of bounds. must be >= 1 and <= %d", i, cs[i], dim)
		}
		ix += c * a.segmentSize(i)
	}
	return ix, nil
}

func (a *Array[T]) Get(cs []int) (t T, err error) {
	ix, err := a.idx(cs)
	if err != nil {
		return t, err
	}
	return a.data[ix], nil
}

func (a *Array[T]) Set(cs []int, t T) error {
	ix, err := a.idx(cs)
	if err != nil {
		return err
	}
	a.data[ix] = t
	return nil
}
