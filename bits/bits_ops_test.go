package bits

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetBit(t *testing.T) {
	tbl := []struct {
		in  int
		pos uint8
		out int
	}{
		{0b11110111, 3, 0b11111111},
		{0b11100111, 4, 0b11110111},
		{0b11110111, 4, 0b11110111},
	}
	for _, el := range tbl {
		res := SetBit(el.in, el.pos)
		assert.Equal(t, el.out, res, "incorrect bit set value:\nexpected: %b\nactual  : %b", el.out, res)
	}
}

func TestClearBit(t *testing.T) {
	tbl := []struct {
		in  int
		pos uint8
		out int
	}{
		{0b11111111, 3, 0b11110111},
		{0b11110111, 4, 0b11100111},
		{0b11110111, 3, 0b11110111},
	}
	for _, el := range tbl {
		res := ClearBit(el.in, el.pos)
		assert.Equal(t, el.out, res, "incorrect clear bit value:\nexpected: %b\nactual  : %b", el.out, res)
	}
}

func TestHasBit(t *testing.T) {
	tbl := []struct {
		in  int
		pos uint8
		out bool
	}{
		{0b11110111, 3, false},
		{0b11110111, 4, true},
	}
	for _, el := range tbl {
		res := HasBit(el.in, el.pos)
		assert.Equal(t, el.out, res, "incorrect bit check value:\nexpected: %t\nactual  : %t", el.out, res)
	}
}
