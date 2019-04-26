package jubjub

import "fmt"

type AffinePoint struct {
	field *Field
	x     *FieldElement
	y     *FieldElement
}

func (p *AffinePoint) NewPoint(field *Field) *AffinePoint {
	p.field = field
	return p.Identity()
}

func (p *AffinePoint) MontgomeryDecodeUncompressed(out []byte) {
	p.field.Mul(p.x, &FieldElement{1, 0, 0, 0}, p.x)
	p.field.Mul(p.y, &FieldElement{1, 0, 0, 0}, p.y)
	p.x.Marshal(out[:32])
	p.y.Marshal(out[32:64])
}

func (p *AffinePoint) Set(p2 *AffinePoint) *AffinePoint {
	p.NewPoint(p2.field)
	p.x.Set(p2.x)
	p.y.Set(p2.y)
	return p
}

func (p *AffinePoint) Eq(q *AffinePoint) bool {
	return p.x.Eq(q.x) && p.y.Eq(q.y)
}

func (p *AffinePoint) Identity() *AffinePoint {
	p.x = &FieldElement{0, 0, 0, 0}
	p.y = new(FieldElement).Set(p.field.r1)
	return p
}

func (p *AffinePoint) Neg() *AffinePoint {
	p.field.Neg(p.x, p.x)
	return p
}

func (p *AffinePoint) IsOnCurve(d *FieldElement) bool {
	//  y^2 == 1 + d.x^2.y^2 + x^2
	var A, B, C FieldElement
	p.field.Square(&B, p.x)
	p.field.Square(&C, p.y)
	p.field.Mul(&A, &C, &B)
	p.field.Mul(&A, &A, d)
	p.field.Add(&A, &A, p.field.r1)
	p.field.Add(&A, &A, &B)
	return A.Eq(&C)
}

func (p *AffinePoint) String() string {
	return fmt.Sprintf("(%s, %s)", p.x.String(), p.y.String())
}

func (p *AffinePoint) ToExtended() *ExtendedPoint {
	extended := new(ExtendedPoint).NewPoint(p.field)
	extended.x.Set(p.x)
	extended.y.Set(p.y)
	extended.z.Set(p.field.r1)
	p.field.Mul(extended.t, p.x, p.y)
	return extended
}

func (p *AffinePoint) toProjective() *projectivePoint {
	projective := new(projectivePoint).newPoint(p.field)
	projective.x.Set(p.x)
	projective.y.Set(p.y)
	projective.z.Set(p.field.r1)
	return projective
}

type ExtendedPoint struct {
	field *Field
	x     *FieldElement
	y     *FieldElement
	z     *FieldElement
	t     *FieldElement
}

func (p *ExtendedPoint) NewPoint(field *Field) *ExtendedPoint {
	p.field = field
	return p.Identity()
}

func (p *ExtendedPoint) Set(p2 *ExtendedPoint) *ExtendedPoint {
	p.NewPoint(p2.field)
	p.x.Set(p2.x)
	p.y.Set(p2.y)
	p.z.Set(p2.z)
	p.t.Set(p2.t)
	return p
}

func (p *ExtendedPoint) IsValid() bool {
	var r1, r2 FieldElement
	p.field.Mul(&r1, p.x, p.y)
	p.field.Mul(&r2, p.z, p.t)
	return r1.Eq(&r2)
}

func (p *ExtendedPoint) Eq(p2 *ExtendedPoint) bool {
	var r1, r2, r3, r4 FieldElement
	p.field.Mul(&r1, p.x, p2.z)
	p.field.Mul(&r2, p.z, p2.x)
	p.field.Mul(&r3, p.y, p2.z)
	p.field.Mul(&r4, p.z, p2.y)
	return r1.Eq(&r2) && r3.Eq(&r4)
}

func (p *ExtendedPoint) Identity() *ExtendedPoint {
	p.x = &FieldElement{0, 0, 0, 0}
	p.y = new(FieldElement).Set(p.field.r1)
	p.z = new(FieldElement).Set(p.field.r1)
	p.t = &FieldElement{0, 0, 0, 0}
	return p
}

func (p *ExtendedPoint) Neg() *ExtendedPoint {
	p.field.Neg(p.x, p.x)
	p.field.Neg(p.t, p.t)
	return p
}

func (p *ExtendedPoint) IsOnCurve(d *FieldElement) bool {
	// (Y^2 - X^2).Z^2 = Z^4 + d.X^2.Y^2
	var A, B, C, D FieldElement
	p.field.Square(&B, p.x)
	p.field.Square(&C, p.y)
	p.field.Square(&D, p.z)
	p.field.Sub(&A, &C, &B)
	p.field.Mul(&A, &A, &D)
	p.field.Square(&D, &D)
	p.field.Mul(&B, &C, &B)
	p.field.Mul(&B, &B, d)
	p.field.Add(&B, &B, &D)
	return A.Eq(&B)
}

func (p *ExtendedPoint) String() string {
	return fmt.Sprintf("(%s, %s, %s, %s)", p.x.String(), p.y.String(), p.z.String(), p.t.String())
}

// Given (X:Y:T:Z) in Ee passing to E is cost-free by simply ignoring T
func (p *ExtendedPoint) toProjective() *projectivePoint {
	projective := new(projectivePoint).newPoint(p.field)
	projective.x.Set(p.x)
	projective.y.Set(p.y)
	projective.z.Set(p.z)
	return projective
}

// the triplet (X:Y:Z) corresponds to the affine point (X/Z,Y/Z)
func (p *ExtendedPoint) ToAffine() *AffinePoint {
	affine := new(AffinePoint).NewPoint(p.field)
	zinv := new(FieldElement)
	p.field.InvMontUp(zinv, p.z)
	p.field.Mul(affine.x, zinv, p.x)
	p.field.Mul(affine.y, zinv, p.y)
	return affine
}

type projectivePoint struct {
	field *Field
	x     *FieldElement
	y     *FieldElement
	z     *FieldElement
}

func (p *projectivePoint) newPoint(field *Field) *projectivePoint {
	p.field = field
	return p.identity()
}

func (p *projectivePoint) set(p2 *projectivePoint) *projectivePoint {
	p.newPoint(p2.field)
	p.x.Set(p2.x)
	p.y.Set(p2.y)
	p.z.Set(p2.z)
	return p
}

func (p *projectivePoint) eq() {
	// TODO impl
}

func (p *projectivePoint) identity() *projectivePoint {
	p.x = &FieldElement{0, 0, 0, 0}
	p.y = new(FieldElement).Set(p.field.r1)
	p.z = new(FieldElement).Set(p.field.r1)
	return p
}

func (p *projectivePoint) neg() {
	// TODO impl
}

func (p *projectivePoint) isOnCurve(d *FieldElement) bool {
	// TODO impl
	return false
}

func (p *projectivePoint) string() string {
	return fmt.Sprintf("(%s, %s, %s)", p.x.String(), p.y.String(), p.z.String())
}

func (p *projectivePoint) toExtended() *ExtendedPoint {
	extended := new(ExtendedPoint).NewPoint(p.field)
	p.field.Mul(extended.y, p.y, p.z)
	p.field.Mul(extended.x, p.x, p.z)
	p.field.Square(extended.z, p.z)
	p.field.Mul(extended.t, p.x, p.y)
	extended.field = p.field
	return extended
}

func (p *projectivePoint) toAffine() *AffinePoint {
	affine := new(AffinePoint).NewPoint(p.field)
	zinv := new(FieldElement)
	p.field.InvMontUp(zinv, p.z)
	p.field.Mul(affine.x, zinv, p.x)
	p.field.Mul(affine.y, zinv, p.y)
	affine.field = p.field
	return affine
}
