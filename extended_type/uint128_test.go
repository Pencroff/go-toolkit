package extended_type

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"log"
	"math/big"
	"testing"
)

func getRandBuffer(size int) []byte {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	return buf
}

func randUint128() Uint128 {
	randBuf := getRandBuffer(16)
	return FromBytes(randBuf)
}

func TestUint128(t *testing.T) {
	// test non-arithmetic methods
	for i := 0; i < 1000; i++ {
		x, y := randUint128(), randUint128()
		if i%3 == 0 {
			x = x.Rsh(64)
		} else if i%7 == 0 {
			x = x.Lsh(64)
		}

		if FromBig(x.Big()) != x {
			t.Fatal("FromBig is not the inverse of Big for", x)
		}

		b := make([]byte, 16)
		x.PutBytes(b)
		if FromBytes(b) != x {
			t.Fatal("FromBytes is not the inverse of PutBytes for", x)
		}

		if !x.Equals(x) {
			t.Fatalf("%v does not equal itself", x.Lo)
		}
		if !From64(x.Lo).Equals64(x.Lo) {
			t.Fatalf("%v does not equal itself", x.Lo)
		}

		if x.Cmp(y) != x.Big().Cmp(y.Big()) {
			t.Fatalf("mismatch: cmp(%v,%v) should equal %v, got %v", x, y, x.Big().Cmp(y.Big()), x.Cmp(y))
		} else if x.Cmp(x) != 0 {
			t.Fatalf("%v does not equal itself", x)
		}

		if x.Cmp64(y.Lo) != x.Big().Cmp(From64(y.Lo).Big()) {
			t.Fatalf("mismatch: cmp64(%v,%v) should equal %v, got %v", x, y.Lo, x.Big().Cmp(From64(y.Lo).Big()), x.Cmp64(y.Lo))
		} else if From64(x.Lo).Cmp64(x.Lo) != 0 {
			t.Fatalf("%v does not equal itself", x.Lo)
		}
	}

	// Check FromBig panics
	checkPanic := func(fn func(), msg string) {
		defer func() {
			r := recover()
			if s, ok := r.(string); !ok || s != msg {
				t.Errorf("expected %q, got %q", msg, r)
			}
		}()
		fn()
	}
	checkPanic(func() { _ = FromBig(big.NewInt(-1)) }, "value cannot be negative")
	checkPanic(func() { _ = FromBig(new(big.Int).Lsh(big.NewInt(1), 129)) }, "value overflows Uint128")
}

func TestArithmetic(t *testing.T) {
	// compare Uint128 arithmetic methods to their math/big equivalents, using
	// random values
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
	mod128 := func(i *big.Int) *big.Int {
		// wraparound semantics
		if i.Sign() == -1 {
			i = i.Add(new(big.Int).Lsh(big.NewInt(1), 128), i)
		}
		_, rem := i.QuoRem(i, new(big.Int).Lsh(big.NewInt(1), 128), new(big.Int))
		return rem
	}
	//checkBinOpX := func(x Uint128, op string, y Uint128, fn func(x, y Uint128) Uint128, fnb func(z, x, y *big.Int) *big.Int) {
	//	t.Helper()
	//	rb := fnb(new(big.Int), x.Big(), y.Big())
	//	defer func() {
	//		if r := recover(); r != nil {
	//			if rb.BitLen() <= 128 && rb.Sign() >= 0 {
	//				t.Fatalf("mismatch: %v%v%v should not panic, %v", x, op, y, rb)
	//			}
	//		} else if rb.BitLen() > 128 || rb.Sign() < 0 {
	//			t.Fatalf("mismatch: %v%v%v should panic, %v", x, op, y, rb)
	//		}
	//	}()
	//	r := fn(x, y)
	//	if r.Big().Cmp(rb) != 0 {
	//		t.Fatalf("mismatch: %v%v%v should equal %v, got %v", x, op, y, rb, r)
	//	}
	//}
	checkBinOp := func(x Uint128, op string, y Uint128, fn func(x, y Uint128) Uint128, fnb func(z, x, y *big.Int) *big.Int) {
		t.Helper()
		r := fn(x, y)
		rb := mod128(fnb(new(big.Int), x.Big(), y.Big()))
		if r.Big().Cmp(rb) != 0 {
			t.Fatalf("mismatch: %v%v%v should equal %v, got %v", x, op, y, rb, r)
		}
	}
	checkShiftOp := func(x Uint128, op string, n uint, fn func(x Uint128, n uint) Uint128, fnb func(z, x *big.Int, n uint) *big.Int) {
		t.Helper()
		r := fn(x, n)
		rb := mod128(fnb(new(big.Int), x.Big(), n))
		if r.Big().Cmp(rb) != 0 {
			t.Fatalf("mismatch: %v%v%v should equal %v, got %v", x, op, n, rb, r)
		}
	}
	//checkBinOp64X := func(x Uint128, op string, y uint64, fn func(x Uint128, y uint64) Uint128, fnb func(z, x, y *big.Int) *big.Int) {
	//	t.Helper()
	//	xb, yb := x.Big(), From64(y).Big()
	//	rb := fnb(new(big.Int), xb, yb)
	//	defer func() {
	//		if r := recover(); r != nil {
	//			if rb.BitLen() <= 128 && rb.Sign() >= 0 {
	//				t.Fatalf("mismatch: %v%v%v should not panic, %v", x, op, y, rb)
	//			}
	//		} else if rb.BitLen() > 128 || rb.Sign() < 0 {
	//			t.Fatalf("mismatch: %v%v%v should panic, %v", x, op, y, rb)
	//		}
	//	}()
	//	r := fn(x, y)
	//	if r.Big().Cmp(rb) != 0 {
	//		t.Fatalf("mismatch: %v%v%v should equal %v, got %v", x, op, y, rb, r)
	//	}
	//}
	checkBinOp64 := func(x Uint128, op string, y uint64, fn func(x Uint128, y uint64) Uint128, fnb func(z, x, y *big.Int) *big.Int) {
		t.Helper()
		xb, yb := x.Big(), From64(y).Big()
		r := fn(x, y)
		rb := mod128(fnb(new(big.Int), xb, yb))
		if r.Big().Cmp(rb) != 0 {
			t.Fatalf("mismatch: %v%v%v should equal %v, got %v", x, op, y, rb, r)
		}
	}
	for i := 0; i < 1000; i++ {
		x, y, z := randUint128(), randUint128(), uint(randUint128().Lo&0xFF)
		checkBinOp(x, "[+]", y, Uint128.Add, (*big.Int).Add)
		checkBinOp(x, "[-]", y, Uint128.Sub, (*big.Int).Sub)
		checkBinOp(x, "[*]", y, Uint128.Mul, (*big.Int).Mul)
		//checkBinOp(x, "+", y, Uint128.AddWrap, (*big.Int).Add)
		//checkBinOp(x, "-", y, Uint128.SubWrap, (*big.Int).Sub)
		//checkBinOp(x, "*", y, Uint128.MulWrap, (*big.Int).Mul)
		if !y.IsZero() {
			checkBinOp(x, "/", y, Uint128.Div, (*big.Int).Div)
			checkBinOp(x, "%", y, Uint128.Mod, (*big.Int).Mod)
		}
		checkBinOp(x, "&", y, Uint128.And, (*big.Int).And)
		checkBinOp(x, "|", y, Uint128.Or, (*big.Int).Or)
		checkBinOp(x, "^", y, Uint128.Xor, (*big.Int).Xor)
		checkShiftOp(x, "<<", z, Uint128.Lsh, (*big.Int).Lsh)
		checkShiftOp(x, ">>", z, Uint128.Rsh, (*big.Int).Rsh)

		// check 64-bit variants
		y64 := y.Lo
		checkBinOp64(x, "[+]", y64, Uint128.Add64, (*big.Int).Add)
		checkBinOp64(x, "[-]", y64, Uint128.Sub64, (*big.Int).Sub)
		checkBinOp64(x, "[*]", y64, Uint128.Mul64, (*big.Int).Mul)
		//checkBinOp64(x, "+", y64, Uint128.AddWrap64, (*big.Int).Add)
		//checkBinOp64(x, "-", y64, Uint128.SubWrap64, (*big.Int).Sub)
		//checkBinOp64(x, "*", y64, Uint128.MulWrap64, (*big.Int).Mul)
		if y64 != 0 {
			checkBinOp64(x, "/", y64, Uint128.Div64, (*big.Int).Div)
			modfn := func(x Uint128, y uint64) Uint128 {
				return From64(x.Mod64(y))
			}
			checkBinOp64(x, "%", y64, modfn, (*big.Int).Mod)
		}
		checkBinOp64(x, "&", y64, Uint128.And64, (*big.Int).And)
		checkBinOp64(x, "|", y64, Uint128.Or64, (*big.Int).Or)
		checkBinOp64(x, "^", y64, Uint128.Xor64, (*big.Int).Xor)
	}
}

func TestLeadingZeros(t *testing.T) {
	tcs := []struct {
		l     Uint128
		r     Uint128
		zeros int
	}{
		{
			l:     New(0x00, 0xf000000000000000),
			r:     New(0x00, 0x8000000000000000),
			zeros: 1,
		},
		{
			l:     New(0x00, 0xf000000000000000),
			r:     New(0x00, 0xc000000000000000),
			zeros: 2,
		},
		{
			l:     New(0x00, 0xf000000000000000),
			r:     New(0x00, 0xe000000000000000),
			zeros: 3,
		},
		{
			l:     New(0x00, 0xffff000000000000),
			r:     New(0x00, 0xff00000000000000),
			zeros: 8,
		},
		{
			l:     New(0x00, 0x000000000000ffff),
			r:     New(0x00, 0x000000000000ff00),
			zeros: 56,
		},
		{
			l:     New(0xf000000000000000, 0x01),
			r:     New(0x4000000000000000, 0x00),
			zeros: 63,
		},
		{
			l:     New(0xf000000000000000, 0x00),
			r:     New(0x4000000000000000, 0x00),
			zeros: 64,
		},
		{
			l:     New(0xf000000000000000, 0x00),
			r:     New(0x8000000000000000, 0x00),
			zeros: 65,
		},
		{
			l:     New(0x00, 0x00),
			r:     New(0x00, 0x00),
			zeros: 128,
		},
		{
			l:     New(0x01, 0x00),
			r:     New(0x00, 0x00),
			zeros: 127,
		},
	}

	for _, tc := range tcs {
		zeros := tc.l.Xor(tc.r).LeadingZeros()
		if zeros != tc.zeros {
			t.Errorf("mismatch (expected: %d, got: %d)", tc.zeros, zeros)
		}
	}
}

func TestString(t *testing.T) {
	for i := 0; i < 1000; i++ {
		x := randUint128()
		if x.String() != x.Big().String() {
			t.Fatalf("mismatch:\n%v !=\n%v", x.String(), x.Big().String())
		}
		y, err := FromString(x.String())
		if err != nil {
			t.Fatal(err)
		} else if !y.Equals(x) {
			t.Fatalf("mismatch:\n%v !=\n%v", x.String(), y.String())
		}
	}
	// Test 0 string
	if ZeroUint128.String() != "0" {
		t.Fatalf(`Zero.String() should be "0", got %q`, ZeroUint128.String())
	}
	// Test Max string
	if MaxUint128.String() != "340282366920938463463374607431768211455" {
		t.Fatalf(`Max.String() should be "0", got %q`, MaxUint128.String())
	}
	// Test parsing invalid strings
	if _, err := FromString("-1"); err == nil {
		t.Fatal("expected error when parsing -1")
	}
	if _, err := FromString("340282366920938463463374607431768211456"); err == nil {
		t.Fatal("expected error when parsing max+1")
	}
}

func TestToBytes(t *testing.T) {
	rndOut := getRandBuffer(16)

	tbl := []struct {
		in  Uint128
		out []byte
	}{
		{
			in:  ZeroUint128,
			out: make([]byte, 16),
		},
		{
			in:  New(0xff, 0xff),
			out: []byte{0xff, 0, 0, 0, 0, 0, 0, 0, 0xff, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			in:  FromBytes(rndOut),
			out: rndOut,
		},
	}

	for _, el := range tbl {
		expected := el.out
		res := el.in.ToBytes()
		assert.Equal(t, expected, res,
			"incorrect byte slice\nexpected: %v\n!=\nresult  : %v\n",
			expected, res)
	}

}
