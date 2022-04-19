package math

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinMax(t *testing.T) {
	lst := []struct {
		x   int
		y   int
		min int
		max int
	}{
		{1, 2, 1, 2},
		{2, 1, 1, 2},
		{1, 1, 1, 1},
		{2, 2, 2, 2},
		{-1, -2, -2, -1},
		{-2, -1, -2, -1},
		{-1, -1, -1, -1},
		{-2, -2, -2, -2},
	}
	for _, el := range lst {
		assert.Equal(t, el.min, Min(el.x, el.y), "Min(%d, %d) should be %d, but got %d", el.x, el.y, el.min, Min(el.x, el.y))
		assert.Equal(t, el.max, Max(el.x, el.y), "Max(%d, %d) should be %d, but got %d", el.x, el.y, el.max, Max(el.x, el.y))
	}
}
