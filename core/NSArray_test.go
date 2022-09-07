package core_test

import (
	"testing"

	"github.com/progrium/macdriver/core"
	"github.com/stretchr/testify/assert"
)

func TestNSArraySize(t *testing.T) {
	arr := core.NSArray_array()
	arr = arr.ArrayByAddingObject_(core.String("a"))
	arr = arr.ArrayByAddingObject_(core.String("b"))
	arr = arr.ArrayByAddingObject_(core.String("c"))
	assert.EqualValues(t, 3, arr.Count())
	assert.Equal(t, []string{"a", "b", "c"}, arr.Strings())
}
