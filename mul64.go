// +build pure_go

package jubjub

// https://github.com/golang/go/blob/master/src/math/bits/bits.go
func mul64(a, b uint64) (hi, lo uint64) {
	const mask32 = 1<<32 - 1
	a0 := a & mask32
	a1 := a >> 32
	b0 := b & mask32
	b1 := b >> 32
	w0 := a0 * b0
	t := a1*b0 + w0>>32
	w1 := t & mask32
	w2 := t >> 32
	w1 += a0 * b1
	hi = a1*b1 + w2 + w1>>32
	lo = a * b
	return
}
