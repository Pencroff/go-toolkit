package extended_type

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/Pencroff/go-toolkit/bits"
	"math/big"
	"testing"
)

func BenchmarkArithmetic(b *testing.B) {
	randBuf := make([]byte, 17)
	randUint128 := func() Uint128 {
		_, err := rand.Read(randBuf)
		if err != nil {
			return Uint128{}
		}
		var Lo, Hi uint64
		if randBuf[16]&1 != 0 {
			Lo = binary.LittleEndian.Uint64(randBuf[:8])
		}
		if randBuf[16]&2 != 0 {
			Hi = binary.LittleEndian.Uint64(randBuf[8:])
		}
		return New(Lo, Hi)
	}
	x, y := randUint128(), randUint128()
	x64, y64 := x.Lo, y.Lo

	b.Run("Add native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = x64 + y64
		}
	})

	b.Run("Sub native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = x64 - y64
		}
	})

	b.Run("Mul native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = x64 * y64
		}
	})

	b.Run("Lsh native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = bits.RotateL64(x.Lo, 17)
		}
	})

	b.Run("Rsh native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = bits.RotateR64(x.Lo, 17)
		}
	})

	b.Run("Add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x.Add(y)
		}
	})

	b.Run("Sub", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x.Sub(y)
		}
	})

	b.Run("Mul", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x.Mul(y)
		}
	})

	b.Run("Lsh", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x.Lsh(17)
		}
	})

	b.Run("Rsh", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x.Rsh(17)
		}
	})

	b.Run("Cmp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x.Cmp(y)
		}
	})

	b.Run("Cmp64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x.Cmp64(y.Lo)
		}
	})
}

func BenchmarkDivision(b *testing.B) {
	randBuf := make([]byte, 8)
	randU64 := func() uint64 {
		_, err := rand.Read(randBuf)
		if err != nil {
			return uint64(0)
		}
		return binary.LittleEndian.Uint64(randBuf) | 3 // avoid divide-by-zero
	}
	x64 := From64(randU64())
	y64 := From64(randU64())
	x128 := New(randU64(), randU64())
	y128 := New(randU64(), randU64())

	b.Run("native 64/64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = x64.Lo / y64.Lo
		}
	})
	b.Run("Div64 64/64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x64.Div64(y64.Lo)
		}
	})
	b.Run("Div64 128/64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x128.Div64(y64.Lo)
		}
	})
	b.Run("Div 64/64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x64.Div(y64)
		}
	})
	b.Run("Div 128/64-Lo", func(b *testing.B) {
		x := x128
		x.Hi = y64.Lo - 1
		for i := 0; i < b.N; i++ {
			x.Div(y64)
		}
	})
	b.Run("Div 128/64-Hi", func(b *testing.B) {
		x := x128
		x.Hi = y64.Lo + 1
		for i := 0; i < b.N; i++ {
			x.Div(y64)
		}
	})
	b.Run("Div 128/128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x128.Div(y128)
		}
	})
	b.Run("big.Int 128/64", func(b *testing.B) {
		xb, yb := x128.Big(), y64.Big()
		q := new(big.Int)
		for i := 0; i < b.N; i++ {
			q = q.Div(xb, yb)
		}
	})
	b.Run("big.Int 128/128", func(b *testing.B) {
		xb, yb := x128.Big(), y128.Big()
		q := new(big.Int)
		for i := 0; i < b.N; i++ {
			q = q.Div(xb, yb)
		}
	})
}

func BenchmarkString(b *testing.B) {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return
	}
	x := New(
		binary.LittleEndian.Uint64(buf[:8]),
		binary.LittleEndian.Uint64(buf[8:]),
	)
	xb := x.Big()
	b.Run("Uint128", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = x.String()
		}
	})
	b.Run("big.Int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = xb.String()
		}
	})
}
