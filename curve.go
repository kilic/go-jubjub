package jubjub

import "fmt"

type Curve struct {
	d         *FieldElement
	twoD      *FieldElement
	field     *Field
	generator *ExtendedPoint
}

func (e *Curve) NewAffinePoint() *AffinePoint {
	return new(AffinePoint).NewPoint(e.field)
}

func (e *Curve) NewExtendedPoint() *ExtendedPoint {
	return new(ExtendedPoint).NewPoint(e.field)
}

func (e *Curve) newProjectivePoint() *projectivePoint {
	return new(projectivePoint).newPoint(e.field)
}

func (e *Curve) NewAffinePointFromUncompressed(in []byte) (*AffinePoint, error) {
	if len(in) != 64 {
		return nil, fmt.Errorf("bad uncompressed point input size")
	}
	point := new(AffinePoint).NewPoint(e.field)
	point.x = e.field.NewElement(in[0:32])
	point.y = e.field.NewElement(in[32:64])
	if !e.AffinePointIsOnCurve(point) {
		return nil, fmt.Errorf("point is not on curve")
	}
	return point, nil
}

// TODO
// require square root implementation
// func (e *Curve) NewAffinePointFromCompressed(in []byte) *AffinePoint, error {
// }

func (e *Curve) NewExtendedPointFromUncompressed(in []byte) (*ExtendedPoint, error) {
	if len(in) != 64 {
		return nil, fmt.Errorf("bad uncompressed point input size")
	}
	point := new(AffinePoint).NewPoint(e.field)
	point.x = e.field.NewElement(in[0:32])
	point.y = e.field.NewElement(in[32:64])
	if !e.AffinePointIsOnCurve(point) {
		return nil, fmt.Errorf("point is not on curve")
	}
	return point.ToExtended(), nil
}

// TODO
// require square root implementation
// func (e *Curve) NewExtendedPointFromCompressed(in []byte) *ExtendedPoint, error {
// }

// add-2008-hwcd-3
// http://www.hyperelliptic.org/EFD/g1p/auto-twisted-extended-1
// 2008 Hisil–Wong–Carter–Dawson, http://eprint.iacr.org/2008/522, Section 3.1.
func (e *Curve) Add(r *ExtendedPoint, p *ExtendedPoint, q *ExtendedPoint) {
	o := e.field
	var t1, t2, t3, t4, t5, t6 FieldElement
	o.Sub(&t1, p.y, p.x)
	o.Sub(&t2, q.y, q.x)
	o.Mul(&t3, &t1, &t2)
	o.Add(&t1, p.y, p.x)
	o.Add(&t2, q.y, q.x)
	o.Mul(&t4, &t1, &t2)
	o.Mul(&t5, p.t, e.twoD)
	o.Mul(&t5, &t5, q.t)
	o.Mul(&t6, p.z, q.z)
	o.Double(&t6, &t6)
	o.Sub(&t1, &t4, &t3)
	o.Sub(&t2, &t6, &t5)
	o.Add(&t3, &t4, &t3)
	o.Add(&t4, &t6, &t5)
	o.Mul(r.x, &t1, &t2)
	o.Mul(r.y, &t4, &t3)
	o.Mul(r.t, &t1, &t3)
	o.Mul(r.z, &t2, &t4)
}

func (e *Curve) Sub(r *ExtendedPoint, p *ExtendedPoint, q *ExtendedPoint) {
	e.Add(r, p, q.Neg())
}

// dbl-2008-bbjlp
// https://hyperelliptic.org/EFD/g1p/auto-twisted-projective.html
// 2008 Bernstein–Birkner–Joye–Lange–Peters http://eprint.iacr.org/2008/013 Section 6
// reduced version
// https://github.com/zkcrypto/jubjub/blob/master/src/lib.rs
func (e *Curve) double(r *ExtendedPoint, p *projectivePoint) {
	o := e.field
	var t1, t2, t3, t4 FieldElement
	o.Add(&t1, p.x, p.y)
	o.Square(&t1, &t1)
	o.Square(&t2, p.x)
	o.Square(&t3, p.y)
	o.Add(&t4, &t3, &t2)
	o.Sub(&t2, &t3, &t2)
	o.Square(&t3, p.z)
	o.Double(&t3, &t3)
	o.Sub(&t3, &t3, &t2)
	o.Sub(&t1, &t1, &t4)
	o.Mul(r.x, &t1, &t3)
	o.Mul(r.y, &t2, &t4)
	o.Mul(r.z, &t2, &t3)
	o.Mul(r.t, &t1, &t4)
}

func (c *Curve) AddBase(r *ExtendedPoint, p *ExtendedPoint) {
	c.Add(r, c.generator, p)
}

func (e *Curve) Mul(r *ExtendedPoint, p *ExtendedPoint, a *ScalarFieldElement) {

	N := new(ExtendedPoint).Set(p)
	Q := new(ExtendedPoint).NewPoint(e.field)
	l := a.bitLength()
	for i := 0; i < l; i++ {
		if a.bit(i) == 1 {
			e.Add(Q, Q, N)
		}
		e.double(N, N.toProjective())
	}
	r.Set(Q)
}

func (e *Curve) MulBase(r *ExtendedPoint, s *ScalarFieldElement) {
	e.Mul(r, e.generator, s)
}

func (e *Curve) AffinePointIsOnCurve(p *AffinePoint) bool {
	return p.IsOnCurve(e.d)
}

func (e *Curve) ExtendedPointIsOnCurve(p *ExtendedPoint) bool {
	return p.IsOnCurve(e.d)
}

// TODO
// requires field square root
// func (c *Curve) RandomPoint() *AffinePoint {
// }
