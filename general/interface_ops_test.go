package general

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterfaceSlice(t *testing.T) {
	type tstStruct struct{}
	var strSlice []interface{}
	res := InterfaceSlice(nil)
	assert.Equal(t, []interface{}(nil), res)
	assert.Equal(t, 0, len(res))
	strSlice = append(strSlice, "a", "b")
	assert.Equal(t, strSlice, InterfaceSlice([]string{"a", "b"}))
	var boolSlice []interface{}
	boolSlice = append(boolSlice, true, false, true)
	assert.Equal(t, boolSlice, InterfaceSlice([]bool{true, false, true}))
	var intSlice []interface{}
	intSlice = append(intSlice, -1, 1, 2)
	assert.Equal(t, intSlice, InterfaceSlice([]int{-1, 1, 2}))
	var uintSlice []interface{}
	uintSlice = append(uintSlice, uint(1), uint(2))
	assert.Equal(t, uintSlice, InterfaceSlice([]uint{1, 2}))
	var float64Slice []interface{}
	float64Slice = append(float64Slice, 1.1, 2.2)
	assert.Equal(t, float64Slice, InterfaceSlice([]float64{1.1, 2.2}))
	var structSlice []interface{}
	structSlice = append(structSlice, tstStruct{}, tstStruct{})
	assert.Equal(t, structSlice, InterfaceSlice([]tstStruct{tstStruct{}, tstStruct{}}))
	assert.Panicsf(t, func() { InterfaceSlice("abc") },
		"InterfaceSlice() given a non-slice type")
}
