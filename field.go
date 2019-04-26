package jubjub

import (
	"io"
	"math/big"
	"math/bits"
)

type Field struct {
	p *FieldElement
	// inp = (-p^{-1} mod 2^b) where b = 64
	// p2  = p-2
	// rN1 = r^1 modp
	// r1  = r modp
	// r2  = r^2 modp
	// r3  = r^3 modp
	inp uint64
	p2  *FieldElement
	rN1 *FieldElement
	r1  *FieldElement
	r2  *FieldElement
	r3  *FieldElement
}

// Given prime number as big.Int,
// field constants are precomputed
func NewField(pBig *big.Int) *Field {
	p := new(FieldElement).Unmarshal(pBig.Bytes())
	inp := bn().ModInverse(bn().Neg(pBig), bn().Exp(big2, big64, nil))
	r1Big := bn().Exp(big2, big256, nil)
	r1 := new(FieldElement).Unmarshal(bn().Mod(r1Big, pBig).Bytes())
	r2 := new(FieldElement).Unmarshal(bn().Exp(r1Big, big2, pBig).Bytes())
	r3 := new(FieldElement).Unmarshal(bn().Exp(r1Big, big3, pBig).Bytes())
	rN1 := new(FieldElement).Unmarshal(bn().ModInverse(r1Big, pBig).Bytes())
	p2 := new(FieldElement).Unmarshal(bn().Sub(pBig, big2).Bytes())
	return &Field{
		p:   p,
		inp: inp.Uint64(),
		p2:  p2,
		r1:  r1,
		rN1: rN1,
		r2:  r2,
		r3:  r3,
	}
}

// Returns new element in Montgomery domain
func (f *Field) NewElement(in []byte) *FieldElement {
	fe := new(FieldElement).Unmarshal(in)
	f.Mul(fe, fe, f.r2)
	return fe
}

// Adapted from https://github.com/golang/go/blob/master/src/crypto/rand/util.go
func (f *Field) RandElement(fe *FieldElement, r io.Reader) error {
	// assuming p > 2^192
	bitLen := bits.Len64(f.p[3]) + 64 + 64 + 64
	// k is the maximum byte length needed to encode a value < max.
	k := (bitLen + 7) / 8
	// b is the number of bits in the most significant byte of max-1.
	b := uint(bitLen % 8)
	if b == 0 {
		b = 8
	}
	bytes := make([]byte, k)
	for {
		_, err := io.ReadFull(r, bytes)
		if err != nil {
			return err
		}
		// Clear bits in the first byte to increase the probability
		// that the candidate is < max.
		bytes[0] &= uint8(int(1<<b) - 1)
		fe.Unmarshal(bytes)

		if fe.Cmp(f.p) < 0 {
			break
		}
	}
	return nil
}

// c = (a + b) modp
func (f *Field) Add(c, a, b *FieldElement) {
	a0 := a[0]
	a1 := a[1]
	a2 := a[2]
	a3 := a[3]
	b0 := b[0]
	b1 := b[1]
	b2 := b[2]
	b3 := b[3]
	p0 := f.p[0]
	p1 := f.p[1]
	p2 := f.p[2]
	p3 := f.p[3]
	var e, e2, ne uint64

	u0 := a0 + b0
	e = (a0&b0 | (a0|b0)&^u0) >> 63
	u1 := a1 + b1 + e
	e = (a1&b1 | (a1|b1)&^u1) >> 63
	u2 := a2 + b2 + e
	e = (a2&b2 | (a2|b2)&^u2) >> 63
	u3 := a3 + b3 + e
	e = (a3&b3 | (a3|b3)&^u3) >> 63

	v0 := u0 - p0
	e2 = (^u0&p0 | (^u0|p0)&v0) >> 63
	v1 := u1 - p1 - e2
	e2 = (^u1&p1 | (^u1|p1)&v1) >> 63
	v2 := u2 - p2 - e2
	e2 = (^u2&p2 | (^u2|p2)&v2) >> 63
	v3 := u3 - p3 - e2
	e2 = (^u3&p3 | (^u3|p3)&v3) >> 63

	e = e - e2
	ne = ^e

	c[0] = (u0 & e) | (v0 & ne)
	c[1] = (u1 & e) | (v1 & ne)
	c[2] = (u2 & e) | (v2 & ne)
	c[3] = (u3 & e) | (v3 & ne)
}

// c = (a + a) modp
func (f *Field) Double(c, a *FieldElement) {

	a0 := a[0]
	a1 := a[1]
	a2 := a[2]
	a3 := a[3]
	p0 := f.p[0]
	p1 := f.p[1]
	p2 := f.p[2]
	p3 := f.p[3]

	e := a3 >> 63
	u3 := a3<<1 | a2>>63
	u2 := a2<<1 | a1>>63
	u1 := a1<<1 | a0>>63
	u0 := a0 << 1

	v0 := u0 - p0
	e2 := (^u0&p0 | (^u0|p0)&v0) >> 63
	v1 := u1 - p1 - e2
	e2 = (^u1&p1 | (^u1|p1)&v1) >> 63
	v2 := u2 - p2 - e2
	e2 = (^u2&p2 | (^u2|p2)&v2) >> 63
	v3 := u3 - p3 - e2
	e2 = (^u3&p3 | (^u3|p3)&v3) >> 63

	e = e - e2
	ne := ^e

	c[0] = (u0 & e) | (v0 & ne)
	c[1] = (u1 & e) | (v1 & ne)
	c[2] = (u2 & e) | (v2 & ne)
	c[3] = (u3 & e) | (v3 & ne)
}

// c = (a - b) modp
func (f *Field) Sub(c, a, b *FieldElement) {
	a0 := a[0]
	a1 := a[1]
	a2 := a[2]
	a3 := a[3]
	b0 := b[0]
	b1 := b[1]
	b2 := b[2]
	b3 := b[3]
	p0 := f.p[0]
	p1 := f.p[1]
	p2 := f.p[2]
	p3 := f.p[3]

	var e, e2, ne uint64

	u0 := a0 - b0
	e = (^a0&b0 | (^a0|b0)&u0) >> 63
	u1 := a1 - b1 - e
	e = (^a1&b1 | (^a1|b1)&u1) >> 63
	u2 := a2 - b2 - e
	e = (^a2&b2 | (^a2|b2)&u2) >> 63
	u3 := a3 - b3 - e
	e = (^a3&b3 | (^a3|b3)&u3) >> 63

	v0 := u0 + p0
	e2 = (u0&p0 | (u0|p0)&^v0) >> 63
	v1 := u1 + p1 + e2
	e2 = (u1&p1 | (u1|p1)&^v1) >> 63
	v2 := u2 + p2 + e2
	e2 = (u2&p2 | (u2|p2)&^v2) >> 63
	v3 := u3 + p3 + e2

	e--
	ne = ^e
	c[0] = (u0 & e) | (v0 & ne)
	c[1] = (u1 & e) | (v1 & ne)
	c[2] = (u2 & e) | (v2 & ne)
	c[3] = (u3 & e) | (v3 & ne)
}

func (f *Field) Neg(c, a *FieldElement) {
	f.Sub(c, f.p, a)
}

// Sets c as a^2(R^-1) modp
func (f *Field) Square(c, a *FieldElement) {
	var T [8]uint64
	square256(&T, *a)
	f.montReduce(c, T)
}

// Sets c as ab(R^-1) modp
func (f *Field) Mul(c, a, b *FieldElement) {
	var T [8]uint64
	mul256(&T, *a, *b)
	f.montReduce(c, T)
}

// Reduces T as T (R^-1) modp
// Handbook of Applied Cryptography
// Hankerson, Menezes, Vanstone
// Algorithm 14.32 Montgomery reduction
func (f *Field) montReduce(c *FieldElement, w [8]uint64) {
	w0 := w[0]
	w1 := w[1]
	w2 := w[2]
	w3 := w[3]
	w4 := w[4]
	w5 := w[5]
	w6 := w[6]
	w7 := w[7]
	p0 := f.p[0]
	p1 := f.p[1]
	p2 := f.p[2]
	p3 := f.p[3]
	var e1, e2, el, res uint64
	var t1, t2, u uint64

	// i = 0
	u = w0 * f.inp
	//
	e1, res = mul64(u, p0)
	t1 = res + w0
	e1 += (res&w0 | (res|w0)&^t1) >> 63
	w0 = t1
	//
	e2, res = mul64(u, p1)
	t1 = res + e1
	e2 += (res&e1 | (res|e1)&^t1) >> 63
	t2 = t1 + w1
	e2 += (t1&w1 | (t1|w1)&^t2) >> 63
	w1 = t2
	//
	e1, res = mul64(u, p2)
	t1 = res + e2
	e1 += (res&e2 | (res|e2)&^t1) >> 63
	t2 = t1 + w2
	e1 += (t1&w2 | (t1|w2)&^t2) >> 63
	w2 = t2
	//
	e2, res = mul64(u, p3)
	t1 = res + e1
	e2 += (res&e1 | (res|e1)&^t1) >> 63
	t2 = t1 + w3
	e2 += (t1&w3 | (t1|w3)&^t2) >> 63
	w3 = t2
	//
	t1 = w4 + el
	e1 = (w4&el | (w4|el)&^t1) >> 63
	t2 = t1 + e2
	e1 += (t1&e2 | (t1|e2)&^t2) >> 63
	w4 = t2
	el = e1

	// i = 1
	u = w1 * f.inp
	//
	e1, res = mul64(u, p0)
	t1 = res + w1
	e1 += (res&w1 | (res|w1)&^t1) >> 63
	w1 = t1
	//
	e2, res = mul64(u, p1)
	t1 = res + e1
	e2 += (res&e1 | (res|e1)&^t1) >> 63
	t2 = t1 + w2
	e2 += (t1&w2 | (t1|w2)&^t2) >> 63
	w2 = t2
	//
	e1, res = mul64(u, p2)
	t1 = res + e2
	e1 += (res&e2 | (res|e2)&^t1) >> 63
	t2 = t1 + w3
	e1 += (t1&w3 | (t1|w3)&^t2) >> 63
	w3 = t2
	//
	e2, res = mul64(u, p3)
	t1 = res + e1
	e2 += (res&e1 | (res|e1)&^t1) >> 63
	t2 = t1 + w4
	e2 += (t1&w4 | (t1|w4)&^t2) >> 63
	w4 = t2
	//
	t1 = w5 + el
	e1 = (w5&el | (w5|el)&^t1) >> 63
	t2 = t1 + e2
	e1 += (t1&e2 | (t1|e2)&^t2) >> 63
	w5 = t2
	el = e1

	// i = 2
	u = w2 * f.inp
	//
	e1, res = mul64(u, p0)
	t1 = res + w2
	e1 += (res&w2 | (res|w2)&^t1) >> 63
	w2 = t1
	//
	e2, res = mul64(u, p1)
	t1 = res + e1
	e2 += (res&e1 | (res|e1)&^t1) >> 63
	t2 = t1 + w3
	e2 += (t1&w3 | (t1|w3)&^t2) >> 63
	w3 = t2
	//
	e1, res = mul64(u, p2)
	t1 = res + e2
	e1 += (res&e2 | (res|e2)&^t1) >> 63
	t2 = t1 + w4
	e1 += (t1&w4 | (t1|w4)&^t2) >> 63
	w4 = t2
	//
	e2, res = mul64(u, p3)
	t1 = res + e1
	e2 += (res&e1 | (res|e1)&^t1) >> 63
	t2 = t1 + w5
	e2 += (t1&w5 | (t1|w5)&^t2) >> 63
	w5 = t2
	//
	t1 = w6 + el
	e1 = (w6&el | (w6|el)&^t1) >> 63
	t2 = t1 + e2
	e1 += (t1&e2 | (t1|e2)&^t2) >> 63
	w6 = t2
	el = e1

	// i = 3
	u = w3 * f.inp
	//
	e1, res = mul64(u, p0)
	t1 = res + w3
	e1 += (res&w3 | (res|w3)&^t1) >> 63
	w3 = t1
	//
	e2, res = mul64(u, p1)
	t1 = res + e1
	e2 += (res&e1 | (res|e1)&^t1) >> 63
	t2 = t1 + w4
	e2 += (t1&w4 | (t1|w4)&^t2) >> 63
	w4 = t2
	//
	e1, res = mul64(u, p2)
	t1 = res + e2
	e1 += (res&e2 | (res|e2)&^t1) >> 63
	t2 = t1 + w5
	e1 += (t1&w5 | (t1|w5)&^t2) >> 63
	w5 = t2
	//
	e2, res = mul64(u, p3)
	t1 = res + e1
	e2 += (res&e1 | (res|e1)&^t1) >> 63
	t2 = t1 + w6
	e2 += (t1&w6 | (t1|w6)&^t2) >> 63
	w6 = t2
	//
	t1 = w7 + el
	e1 = (w7&el | (w7|el)&^t1) >> 63
	t2 = t1 + e2
	e1 += (t1&e2 | (t1|e2)&^t2) >> 63
	w7 = t2

	e1--
	c[0] = w4 - ((p0) & ^e1)
	e2 = (^w4&p0 | (^w4|p0)&c[0]) >> 63
	c[1] = w5 - ((p1 + e2) & ^e1)
	e2 = (^w5&p1 | (^w5|p1)&c[1]) >> 63
	c[2] = w6 - ((p2 + e2) & ^e1)
	e2 = (^w6&p2 | (^w6|p2)&c[2]) >> 63
	c[3] = w7 - ((p3 + e2) & ^e1)

	f.Sub(c, c, f.p)
}

// Guide to Elliptic Curve Cryptography Algorithm
// Hankerson, Menezes, Vanstone
// Algoritm 2.22 Binary algorithm for inversion in Fp
// Input: a
// Output: a^-1
func (f *Field) InvEEA(inv, fe *FieldElement) {
	u := new(FieldElement).Set(fe)
	v := new(FieldElement).Set(f.p)
	p := new(FieldElement).Set(f.p)
	x1 := &FieldElement{1, 0, 0, 0}
	x2 := &FieldElement{0, 0, 0, 0}
	var e uint64

	for !u.IsOne() && !v.IsOne() {
		//
		for u.IsEven() {
			u.rightShift(0)
			if x1.IsEven() {
				x1.rightShift(0)
			} else {
				e = x1.add(p)
				x1.rightShift(e)
			}
		}
		//
		for v.IsEven() {
			v.rightShift(0)
			if x2.IsEven() {
				x2.rightShift(0)
			} else {
				e = x2.add(p)
				x2.rightShift(e)
			}
		}
		//
		if u.Cmp(v) == -1 {
			v.sub(u)
			f.Sub(x2, x2, x1)
		} else { //
			u.sub(v)
			f.Sub(x1, x1, x2)
		}
	}
	if u.IsOne() {
		inv.Set(x1)
		return
	}
	inv.Set(x2)
}

// Two phase Montgomery Modular Inverse
// The Montgomery Modular Inverse - Revisited
// Savas, Koc
// &
// Guide to Elliptic Curve Cryptography Algorithm
// Hankerson, Menezes, Vanstone
// Algoritm 2.23 Partial Montgomery inversion in Fp
//
// Input : a
// Output : (a^-1)R
// or
// Input : aR
// Output : (a^-1)
func (f *Field) InvMontDown(inv, fe *FieldElement) {

	u := new(FieldElement).Set(fe)
	v := new(FieldElement).Set(f.p)
	x1 := &FieldElement{1, 0, 0, 0}
	x2 := &FieldElement{0, 0, 0, 0}
	var k int
	// Phase 1
	for !v.IsZero() {
		if v.IsEven() {
			v.rightShift(0)
			x1.leftShift()
		} else if u.IsEven() {
			u.rightShift(0)
			x2.leftShift()
		} else if v.Cmp(u) == -1 {
			u.sub(v)
			u.rightShift(0)
			x1.add(x2)
			x2.leftShift()
		} else {
			v.sub(u)
			v.rightShift(0)
			x2.add(x1)
			x1.leftShift()
		}
		k = k + 1
	}
	// Phase2
	p := new(FieldElement).Set(f.p)
	k = k - 256
	var e uint64
	for i := 0; i < k; i++ {
		if x1.IsEven() {
			x1.rightShift(0)
		} else {
			e = x1.add(p)
			x1.rightShift(e)
		}
	}
	inv.Set(x1)
}

// Inverse value stays in Montgomery space
// Two phase Montgomery Modular Inverse
// The Montgomery Modular Inverse - Revisited
// Savas, Koc
// &
// Guide to Elliptic Curve Cryptography Algorithm
// Hankerson, Menezes, Vanstone
// Algoritm 2.23 Partial Montgomery inversion in Fp
// Input : aR
// Output : (a^-1)R
func (f *Field) InvMontUp(inv, fe *FieldElement) {

	u := new(FieldElement).Set(fe)
	v := new(FieldElement).Set(f.p)
	x1 := &FieldElement{1, 0, 0, 0}
	x2 := &FieldElement{0, 0, 0, 0}
	var k int

	// Phase 1
	for !v.IsZero() {
		if v.IsEven() {
			v.rightShift(0)
			x1.leftShift()
		} else if u.IsEven() {
			u.rightShift(0)
			x2.leftShift()
		} else if v.Cmp(u) == -1 {
			u.sub(v)
			u.rightShift(0)
			x1.add(x2)
			x2.leftShift()
		} else {
			v.sub(u)
			v.rightShift(0)
			x2.add(x1)
			x1.leftShift()
		}
		k = k + 1
	}
	// Phase2
	f.Sub(x1, x1, f.p)
	for i := k; i < 512; i++ {
		f.Double(x1, x1)
	}
	inv.Set(x1)
}
