package jubjub

func add64(a, b, cin uint64) (sum, cout uint64) {
	// assumes cin = {1,0}
	sum = a + b + cin
	cout = (b&a | (b|a)&^sum)
	return
}

func sub64(a, b, bin uint64) (diff, bout uint64) {
	// asumes bin = {1,0}
	diff = a - (b + bin)
	bout = (^a&b | (^a|b)&diff)
	return
}

// https://github.com/golang/go/blob/master/src/math/bits/bits.go
func mul64_g(a, b uint64) (hi, lo uint64) {
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

func square64(a uint64) (hi, lo uint64) {
	return mul64(a, a)
	// const mask32 = 1<<32 - 1
	// a0 := a & mask32
	// a1 := a >> 32
	// w0 := a0 * a0
	// t := a1*a0 + w0>>32
	// w1 := t & mask32
	// w2 := t >> 32
	// w1 += a0 * a1
	// hi = a1*a1 + w2 + w1>>32
	// lo = a * a
	// return
}

func add256(a [4]uint64, b [4]uint64) ([4]uint64, uint64) {
	var e uint64
	a0, a1, a2, a3 := a[0], a[1], a[2], a[3]
	b0, b1, b2, b3 := b[0], b[1], b[2], b[3]
	r0 := a0 + b0
	e = (a0&b0 | (a0|b0)&^r0) >> 63
	r1 := a1 + b1 + e
	e = (a1&b1 | (a1|b1)&^r1) >> 63
	r2 := a2 + b2 + e
	e = (a2&b2 | (a2|b2)&^r2) >> 63
	r3 := a3 + b3 + e
	e = (a3&b3 | (a3|b3)&^r3) >> 63
	return [4]uint64{r0, r1, r2, r3}, e
}

func sub256(a, b [4]uint64) (diff [4]uint64, e uint64) {
	a0, a1, a2, a3 := a[0], a[1], a[2], a[3]
	b0, b1, b2, b3 := b[0], b[1], b[2], b[3]

	diff0 := a0 - b0
	e = (^a0&b0 | (^a0|b0)&diff0) >> 63
	diff1 := a1 - b1 - e
	e = (^a1&b1 | (^a1|b1)&diff1) >> 63
	diff2 := a2 - b2 - e
	e = (^a2&b2 | (^a2|b2)&diff2) >> 63
	diff3 := a3 - b3 - e
	e = (^a3&b3 | (^a3|b3)&diff3) >> 63
	return [4]uint64{diff0, diff1, diff2, diff3}, e
}

// Handbook of Applied Cryptography
// Hankerson, Menezes, Vanstone
// 14.12 Algorithm Multiple-precision multiplication
func mul256(w *[8]uint64, a [4]uint64, b [4]uint64) {

	var w0, w1, w2, w3, w4, w5, w6, w7 uint64
	var a0 = a[0]
	var a1 = a[1]
	var a2 = a[2]
	var a3 = a[3]
	var b0 = b[0]
	var b1 = b[1]
	var b2 = b[2]
	var b3 = b[3]
	var u, v, c, t uint64

	// i = 0, j = 0
	c, w0 = mul64(a0, b0)

	// i = 0, j = 1
	u, v = mul64(a1, b0)
	w1 = v + c
	c = u + (v&c|(v|c)&^w1)>>63

	// i = 0, j = 2
	u, v = mul64(a2, b0)
	w2 = v + c
	c = u + (v&c|(v|c)&^w2)>>63

	// i = 0, j = 3
	u, v = mul64(a3, b0)
	w3 = v + c
	w4 = u + (v&c|(v|c)&^w3)>>63

	//
	// i = 1, j = 0
	c, v = mul64(a0, b1)
	t = v + w1
	c += (v&w1 | (v|w1)&^t) >> 63
	w1 = t

	// i = 1, j = 1
	u, v = mul64(a1, b1)
	t = v + w2
	u += (v&w2 | (v|w2)&^t) >> 63
	w2 = t + c
	c = u + (t&c|(t|c)&^w2)>>63

	// i = 1, j = 2
	u, v = mul64(a2, b1)
	t = v + w3
	u += (v&w3 | (v|w3)&^t) >> 63
	w3 = t + c
	c = u + (t&c|(t|c)&^w3)>>63

	// i = 1, j = 3
	u, v = mul64(a3, b1)
	t = v + w4
	u += (v&w4 | (v|w4)&^t) >> 63
	w4 = t + c
	w5 = u + (t&c|(t|c)&^w4)>>63

	//
	// i = 2, j = 0
	c, v = mul64(a0, b2)
	t = v + w2
	c += (v&w2 | (v|w2)&^t) >> 63
	w2 = t

	// i = 2, j = 1
	u, v = mul64(a1, b2)
	t = v + w3
	u += (v&w3 | (v|w3)&^t) >> 63
	w3 = t + c
	c = u + (t&c|(t|c)&^w3)>>63

	// i = 2, j = 2
	u, v = mul64(a2, b2)
	t = v + w4
	u += (v&w4 | (v|w4)&^t) >> 63
	w4 = t + c
	c = u + (t&c|(t|c)&^w4)>>63

	// i = 2, j = 3
	u, v = mul64(a3, b2)
	t = v + w5
	u += (v&w5 | (v|w5)&^t) >> 63
	w5 = t + c
	w6 = u + (t&c|(t|c)&^w5)>>63

	//
	// i = 3, j = 0
	c, v = mul64(a0, b3)
	t = v + w3
	c += (v&w3 | (v|w3)&^t) >> 63
	w3 = t

	// i = 3, j = 1
	u, v = mul64(a1, b3)
	t = v + w4
	u += (v&w4 | (v|w4)&^t) >> 63
	w4 = t + c
	c = u + (t&c|(t|c)&^w4)>>63

	// i = 3, j = 2
	u, v = mul64(a2, b3)
	t = v + w5
	u += (v&w5 | (v|w5)&^t) >> 63
	w5 = t + c
	c = u + (t&c|(t|c)&^w5)>>63

	// i = 3, j = 3
	u, v = mul64(a3, b3)
	t = v + w6
	u += (v&w6 | (v|w6)&^t) >> 63
	w6 = t + c
	w7 = u + (t&c|(t|c)&^w6)>>63

	w[0] = w0
	w[1] = w1
	w[2] = w2
	w[3] = w3
	w[4] = w4
	w[5] = w5
	w[6] = w6
	w[7] = w7
}

// Handbook of Applied Cryptography
// Hankerson, Menezes, Vanstone
// 14.16 Algorithm Multiple-precision squaring
func square256(w *[8]uint64, a [4]uint64) {

	var w0, w1, w2, w3, w4, w5, w6, w7 uint64
	var u, v, c, vv, uu, z1, z2, z3, e uint64
	var a0 = a[0]
	var a1 = a[1]
	var a2 = a[2]
	var a3 = a[3]

	// i = 0
	c, w0 = square64(a0)

	// i = 0, j = 1
	u, v = mul64(a0, a1)
	z1 = u >> 63 // z1 for w2
	u = u<<1 + v>>63
	v = v << 1
	w1 = v + c
	e = (v&c | (v|c)&^w1) >> 63
	uu = u + e
	z1 += (u&e | (u|e)&^uu) >> 63
	c = uu

	// i = 0, j = 2
	u, v = mul64(a0, a2)
	z2 = u >> 63 // z2 for w3
	u = u<<1 + v>>63
	v = v << 1
	w2 = v + c
	e = z1 + (v&c|(v|c)&^w2)>>63
	uu = u + e
	z2 += (u&e | (u|e)&^uu) >> 63
	c = uu

	// i = 0, j = 3
	u, v = mul64(a0, a3)
	z1 = u >> 63 // z1 for w4
	u = u<<1 + v>>63
	v = v << 1
	w3 = v + c
	e = z2 + (v&c|(v|c)&^w3)>>63
	w4 = u + e
	z1 += (u&e | (u|e)&^w4) >> 63

	// i = 1
	c, v = square64(a1)
	vv = v + w2
	c += (v&w2 | (v|w2)&^vv) >> 63
	w2 = vv

	// i = 1, j = 2
	u, v = mul64(a1, a2)
	z2 = u >> 63 // z2 for w4
	u = u<<1 + v>>63
	v = v << 1
	vv = v + w3
	e = (v&w3 | (v|w3)&^vv) >> 63
	uu = u + e
	z2 += (u&e | (u|e)&^uu) >> 63
	w3 = vv + c
	e = (vv&c | (vv|c)&^w3) >> 63
	c = uu + e
	z2 += (uu&e | (uu|e)&^c) >> 63

	// i = 1, j = 3
	u, v = mul64(a1, a3)
	z3 = u >> 63 // z3 for w5
	u = u<<1 + v>>63
	v = v << 1
	vv = v + w4
	e = z1 + z2 + (v&w4|(v|w4)&^vv)>>63
	uu = u + e
	z3 += (u&e | (u|e)&^uu) >> 63
	w4 = vv + c
	e = (vv&c | (vv|c)&^w4) >> 63
	w5 = uu + e
	z3 += (uu&e | (uu|e)&^w5) >> 63

	// i = 2
	c, v = square64(a2)
	vv = v + w4
	c += (v&w4 | (v|w4)&^vv) >> 63
	w4 = vv

	// i = 2, j = 3
	u, v = mul64(a2, a3)
	z1 = u >> 63 // z1 for w6
	u = u<<1 + v>>63
	v = v << 1
	vv = v + w5
	e = z3 + (v&w5|(v|w5)&^vv)>>63
	uu = u + e
	z1 += (u&e | (u|e)&^uu) >> 63
	v = vv
	u = uu
	w5 = vv + c
	e = (vv&c | (vv|c)&^w5) >> 63
	w6 = uu + e
	z1 += (uu&e | (uu|e)&^w6) >> 63

	// i = 3
	c, v = square64(a3)
	vv = v + w6
	w7 = c + z1 + (v&w6|(v|w6)&^vv)>>63
	w6 = vv

	w[0] = w0
	w[1] = w1
	w[2] = w2
	w[3] = w3
	w[4] = w4
	w[5] = w5
	w[6] = w6
	w[7] = w7
}
