package bits

// source https://stackoverflow.com/questions/23192262/how-would-you-set-and-clear-a-single-bit-in-go

// SetBit sets the bit at position pos to 1
func SetBit(n int, pos uint8) int {
	n |= 1 << pos
	return n
}

// ClearBit sets the bit at position pos to 0
func ClearBit(n int, pos uint8) int {
	return n &^ (1 << pos)
}

// HasBit returns true if the bit at position pos is 1
func HasBit(n int, pos uint8) bool {
	val := n & (1 << pos)
	return val > 0
}
