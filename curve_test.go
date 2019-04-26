package jubjub

import (
	"math/big"
	"testing"
)

func TestCurveAffineGenerator1(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	g := &AffinePoint{
		field: field,
		x:     fe(field, "0x11dafe5d23e1218086a365b99fbf3d3be72f6afd7d1f72623e6b071492d1122b"),
		y:     fe(field, "0x1d523cf1ddab1a1793132e78c866c0c33e26ba5cc220fed7cc3f870e59d292aa"),
	}
	if !g.IsOnCurve(curve.d) {
		t.Errorf("generator is not on curve")
	}
}

func TestCurveAffineGenerator2(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	g := &AffinePoint{
		field: field,
		x:     fe(field, "0x62edcbb8bf3787c88b0f03ddd60a8187caf55d1b29bf81afe4b3d35df1a7adfe"),
		y:     fe(field, "0x0b"),
	}
	if !g.IsOnCurve(curve.d) {
		t.Errorf("generator is not on curve")
	}
}

func TestCurveAffineGenerator3(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	g := &AffinePoint{
		field: field,
		x:     fe(field, "0x187d2619ff114316d237e86684fb6e3c6b15e9b924fa4e322764d3177508297a"),
		y:     fe(field, "0x6230c613f1b460e026221be21cf4eabd5a8ea552db565cb18d3cabc39761eb9b"),
	}
	if !g.IsOnCurve(curve.d) {
		t.Errorf("generator is not on curve")
	}
}

func TestCurveExtendedGenerator1(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	g := &AffinePoint{
		field: field,
		x:     fe(field, "0x11dafe5d23e1218086a365b99fbf3d3be72f6afd7d1f72623e6b071492d1122b"),
		y:     fe(field, "0x1d523cf1ddab1a1793132e78c866c0c33e26ba5cc220fed7cc3f870e59d292aa"),
	}
	e := g.ToExtended()
	if !e.IsOnCurve(curve.d) {
		t.Errorf("generator is not on curve")
	}
}

func TestCurveExtendedGenerator2(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	g := &AffinePoint{
		field: field,
		x:     fe(field, "0x62edcbb8bf3787c88b0f03ddd60a8187caf55d1b29bf81afe4b3d35df1a7adfe"),
		y:     fe(field, "0x0b"),
	}
	e := g.ToExtended()
	if !e.IsOnCurve(curve.d) {
		t.Errorf("generator is not on curve")
	}
}

func TestCurveExtendedGenerator3(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	g := &AffinePoint{
		field: field,
		x:     fe(field, "0x187d2619ff114316d237e86684fb6e3c6b15e9b924fa4e322764d3177508297a"),
		y:     fe(field, "0x6230c613f1b460e026221be21cf4eabd5a8ea552db565cb18d3cabc39761eb9b"),
	}
	e := g.ToExtended()
	if !e.IsOnCurve(curve.d) {
		t.Errorf("generator is not on curve")
	}
}

func TestCurveIdentity(t *testing.T) {
	curve := NewJubjub()
	g1 := curve.NewAffinePoint()
	g2 := curve.NewExtendedPoint()
	if !g1.IsOnCurve(curve.d) {
		t.Errorf("affine coordinate identity is not curve")
	}
	if !g2.IsOnCurve(curve.d) {
		t.Errorf("extended coordinate identity is not curve")
	}
}

func TestCurveIdentityConversion1(t *testing.T) {
	curve := NewJubjub()
	g1 := curve.NewAffinePoint()
	g2 := curve.NewExtendedPoint().ToAffine()
	if !g1.Eq(g2) {
		t.Errorf("coordinate conversion fails")
	}
}

func TestCurveIdentityConversion2(t *testing.T) {
	curve := NewJubjub()
	g1 := curve.NewAffinePoint().ToExtended()
	g2 := curve.NewExtendedPoint()
	if !g1.Eq(g2) {
		t.Errorf("coordinate conversion fails")
	}
}

func TestCurveExtendedToAffineConvertion(t *testing.T) {
	field := NewJubjub().field
	g := &AffinePoint{
		field: field,
		x:     fe(field, "0x187d2619ff114316d237e86684fb6e3c6b15e9b924fa4e322764d3177508297a"),
		y:     fe(field, "0x6230c613f1b460e026221be21cf4eabd5a8ea552db565cb18d3cabc39761eb9b"),
	}
	h := g.ToExtended().ToAffine()
	if !g.Eq(h) {
		t.Errorf("bad affine to extended conversion")
	}
}

func TestCurveAdditiveIdentity(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	a := &AffinePoint{
		field: field,
		x:     fe(field, "0x187d2619ff114316d237e86684fb6e3c6b15e9b924fa4e322764d3177508297a"),
		y:     fe(field, "0x6230c613f1b460e026221be21cf4eabd5a8ea552db565cb18d3cabc39761eb9b"),
	}
	g := a.ToExtended()
	h := curve.NewExtendedPoint()
	r := curve.NewExtendedPoint()
	curve.Add(r, g, h)
	if !r.ToAffine().Eq(a) {
		t.Errorf("identity addition is failed")
	}
}

func TestCurveAdditiveIdentity2(t *testing.T) {
	curve := NewJubjub()
	g := curve.newProjectivePoint()
	r := curve.NewExtendedPoint()
	identity := curve.NewExtendedPoint()
	curve.double(r, g)
	if !r.Eq(identity) {
		t.Errorf("identity addition is failed")
	}
}

func TestCurveNegation(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	ga := &AffinePoint{
		field: field,
		x:     fe(field, "0x187d2619ff114316d237e86684fb6e3c6b15e9b924fa4e322764d3177508297a"),
		y:     fe(field, "0x6230c613f1b460e026221be21cf4eabd5a8ea552db565cb18d3cabc39761eb9b"),
	}
	gan := new(AffinePoint).Set(ga).Neg()
	r := curve.NewExtendedPoint()
	identity := curve.NewExtendedPoint()
	curve.Add(r, ga.ToExtended(), gan.ToExtended())
	if !r.Eq(identity) {
		t.Errorf("negation fails")
	}
}

func TestCurveNegation2(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	ga := &AffinePoint{
		field: field,
		x:     fe(field, "0x187d2619ff114316d237e86684fb6e3c6b15e9b924fa4e322764d3177508297a"),
		y:     fe(field, "0x6230c613f1b460e026221be21cf4eabd5a8ea552db565cb18d3cabc39761eb9b"),
	}
	ga2 := new(AffinePoint).Set(ga)
	ga.Neg().Neg()
	if !ga.Eq(ga2) {
		t.Errorf("negation fails")
	}
}

func TestCurveNegation3(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	ga := &AffinePoint{
		field: field,
		x:     fe(field, "0x187d2619ff114316d237e86684fb6e3c6b15e9b924fa4e322764d3177508297a"),
		y:     fe(field, "0x6230c613f1b460e026221be21cf4eabd5a8ea552db565cb18d3cabc39761eb9b"),
	}
	ge := ga.ToExtended()
	ge2 := new(ExtendedPoint).Set(ge)
	ge.Neg()
	ge2.Neg()
	if !ge.Eq(ge2) {
		t.Errorf("negation fails")
	}
}

func TestCurveMultiplication(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	scalarField := NewJubjubScalarField()
	ga := &AffinePoint{
		field: field,
		x:     fe(field, "0x11dafe5d23e1218086a365b99fbf3d3be72f6afd7d1f72623e6b071492d1122b"),
		y:     fe(field, "0x1d523cf1ddab1a1793132e78c866c0c33e26ba5cc220fed7cc3f870e59d292aa"),
	}
	ge1 := ga.ToExtended()
	for i := 0; i < 1; i++ {
		ge2 := curve.NewExtendedPoint()
		ge3 := curve.NewExtendedPoint()
		s1 := scalarField.NewRandElement()
		s2 := scalarField.NewRandElement()
		s3 := scalarField.NewElement()
		scalarField.Mul(s3, s1, s2)
		curve.Mul(ge2, ge1, s1)
		curve.Mul(ge2, ge2, s2)
		curve.Mul(ge3, ge1, s3)
		if !ge2.Eq(ge3) {
			t.Errorf("bad multiplication")
		}
	}
}

func TestCurveAddition(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	scalarField := NewJubjubScalarField()
	ga := &AffinePoint{
		field: field,
		x:     fe(field, "0x11dafe5d23e1218086a365b99fbf3d3be72f6afd7d1f72623e6b071492d1122b"),
		y:     fe(field, "0x1d523cf1ddab1a1793132e78c866c0c33e26ba5cc220fed7cc3f870e59d292aa"),
	}
	ge := ga.ToExtended()
	for i := 0; i < 10; i++ {
		ge1 := curve.NewExtendedPoint()
		ge2 := curve.NewExtendedPoint()
		ge3 := curve.NewExtendedPoint()
		s1 := scalarField.NewRandElement()
		s2 := scalarField.NewRandElement()
		s3 := scalarField.NewElement()
		scalarField.Add(s3, s1, s2)
		curve.Mul(ge1, ge, s1)
		curve.Mul(ge2, ge, s2)
		curve.Mul(ge3, ge, s3)
		curve.Add(ge2, ge1, ge2)
		if !ge2.Eq(ge3) {
			t.Errorf("bad addition")
		}
	}
}

func TestCurveOrder(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	scalarField := NewJubjubScalarField()
	rbig, _ := new(big.Int).SetString("0e7db4ea6533afa906673b0101343b00a6682093ccc81082d0970e5ed6f72cb7", 16)
	r := scalarField.NewElementFromBig(new(big.Int).Set(rbig))
	identity := curve.NewExtendedPoint()
	ga := &AffinePoint{
		field: field,
		x:     fe(field, "0x11dafe5d23e1218086a365b99fbf3d3be72f6afd7d1f72623e6b071492d1122b"),
		y:     fe(field, "0x1d523cf1ddab1a1793132e78c866c0c33e26ba5cc220fed7cc3f870e59d292aa"),
	}
	g1 := ga.ToExtended()
	g2 := curve.NewExtendedPoint()
	for i := 0; i < 8; i++ {
		curve.Mul(g2, g1, r)
		r.n.Add(r.n, rbig)
		if !g2.Eq(identity) {
			t.Errorf("identity element expected")
		}
	}
}

func TestCurveTorsions1(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	scalarField := NewJubjubScalarField()
	cofactor := scalarField.NewElementFromUint64(8)
	identity := curve.NewExtendedPoint()
	te := make([]ExtendedPoint, 8)
	ta := make([]AffinePoint, 8)
	ta[0] = AffinePoint{
		field: field,
		x:     fe(field, "0x71d4df38ba9e7973eaaae086a16618d17aa41ac43dae8582d92e6a7927200d43"),
		y:     fe(field, "0x4958bdb21966982e16a13035ad4d72669106ee90f384a4a1ff0d2068eff496dd"),
	}

	ta[1] = AffinePoint{
		field: field,
		x:     fe(field, "0x73eda753299d7d47a5e80b39939ed33467baa40089fb5bfefffeffff00000001"),
		y:     fe(field, "0x0000000000000000000000000000000000000000000000000000000000000000"),
	}

	ta[2] = AffinePoint{
		field: field,
		x:     fe(field, "0x71d4df38ba9e7973eaaae086a16618d17aa41ac43dae8582d92e6a7927200d43"),
		y:     fe(field, "0x2a94e9a11036e51a1c98a7d25c54659ec2b6b5720c79b75d00f2df96100b6924"),
	}

	ta[3] = AffinePoint{
		field: field,
		x:     fe(field, "0x0000000000000000000000000000000000000000000000000000000000000000"),
		y:     fe(field, "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000000"),
	}

	ta[4] = AffinePoint{
		field: field,
		x:     fe(field, "0x0218c81a6eff03d4488ef781683bbf33d919893ec24fd67c26d19585d8dff2be"),
		y:     fe(field, "0x2a94e9a11036e51a1c98a7d25c54659ec2b6b5720c79b75d00f2df96100b6924"),
	}

	ta[5] = AffinePoint{
		field: field,
		x:     fe(field, "0x00000000000000008d51ccce760304d0ec030002760300000001000000000000"),
		y:     fe(field, "0x0000000000000000000000000000000000000000000000000000000000000000"),
	}

	ta[6] = AffinePoint{
		field: field,
		x:     fe(field, "0x0218c81a6eff03d4488ef781683bbf33d919893ec24fd67c26d19585d8dff2be"),
		y:     fe(field, "0x4958bdb21966982e16a13035ad4d72669106ee90f384a4a1ff0d2068eff496dd"),
	}

	ta[7] = AffinePoint{
		field: field,
		x:     fe(field, "0x0000000000000000000000000000000000000000000000000000000000000000"),
		y:     fe(field, "0x0000000000000000000000000000000000000000000000000000000000000001"),
	}
	for i := 0; i < 8; i++ {
		curve.Mul(&te[i], ta[i].ToExtended(), cofactor)
		if !te[i].Eq(identity) {
			t.Errorf("identity element expected")
		}
	}
}

func TestCurveTorsions2(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	scalarField := NewJubjubScalarField()
	ta := make([]AffinePoint, 8)
	rbig, _ := new(big.Int).SetString("0e7db4ea6533afa906673b0101343b00a6682093ccc81082d0970e5ed6f72cb7", 16)
	r := scalarField.NewElementFromBig(new(big.Int).Set(rbig))
	ga := &AffinePoint{
		field: field,
		x:     fe(field, "0x62edcbb8bf3787c88b0f03ddd60a8187caf55d1b29bf81afe4b3d35df1a7adfe"),
		y:     fe(field, "0x0b"),
	}
	ge1 := ga.ToExtended()
	ge2 := curve.NewExtendedPoint()

	ta[0] = AffinePoint{
		field: field,
		x:     fe(field, "0x71d4df38ba9e7973eaaae086a16618d17aa41ac43dae8582d92e6a7927200d43"),
		y:     fe(field, "0x4958bdb21966982e16a13035ad4d72669106ee90f384a4a1ff0d2068eff496dd"),
	}

	ta[1] = AffinePoint{
		field: field,
		x:     fe(field, "0x73eda753299d7d47a5e80b39939ed33467baa40089fb5bfefffeffff00000001"),
		y:     fe(field, "0x0000000000000000000000000000000000000000000000000000000000000000"),
	}

	ta[2] = AffinePoint{
		field: field,
		x:     fe(field, "0x71d4df38ba9e7973eaaae086a16618d17aa41ac43dae8582d92e6a7927200d43"),
		y:     fe(field, "0x2a94e9a11036e51a1c98a7d25c54659ec2b6b5720c79b75d00f2df96100b6924"),
	}

	ta[3] = AffinePoint{
		field: field,
		x:     fe(field, "0x0000000000000000000000000000000000000000000000000000000000000000"),
		y:     fe(field, "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000000"),
	}

	ta[4] = AffinePoint{
		field: field,
		x:     fe(field, "0x0218c81a6eff03d4488ef781683bbf33d919893ec24fd67c26d19585d8dff2be"),
		y:     fe(field, "0x2a94e9a11036e51a1c98a7d25c54659ec2b6b5720c79b75d00f2df96100b6924"),
	}

	ta[5] = AffinePoint{
		field: field,
		x:     fe(field, "0x00000000000000008d51ccce760304d0ec030002760300000001000000000000"),
		y:     fe(field, "0x0000000000000000000000000000000000000000000000000000000000000000"),
	}

	ta[6] = AffinePoint{
		field: field,
		x:     fe(field, "0x0218c81a6eff03d4488ef781683bbf33d919893ec24fd67c26d19585d8dff2be"),
		y:     fe(field, "0x4958bdb21966982e16a13035ad4d72669106ee90f384a4a1ff0d2068eff496dd"),
	}

	ta[7] = AffinePoint{
		field: field,
		x:     fe(field, "0x0000000000000000000000000000000000000000000000000000000000000000"),
		y:     fe(field, "0x0000000000000000000000000000000000000000000000000000000000000001"),
	}

	for i := 0; i < 8; i++ {
		curve.Mul(ge2, ge1, r)
		r.n.Add(r.n, rbig)
		if !ta[i].Eq(ge2.ToAffine()) {
			t.Errorf("torsion element expected")
		}
	}
}

func BenchmarkCurveAddition(t *testing.B) {
	curve := NewJubjub()
	field := curve.field
	ga := &AffinePoint{
		field: field,
		x:     fe(field, "0x11dafe5d23e1218086a365b99fbf3d3be72f6afd7d1f72623e6b071492d1122b"),
		y:     fe(field, "0x1d523cf1ddab1a1793132e78c866c0c33e26ba5cc220fed7cc3f870e59d292aa"),
	}
	ge1 := ga.ToExtended()
	ge2 := ga.ToExtended()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		curve.Add(ge1, ge1, ge2)
	}
}

func BenchmarkCurveDoubling(t *testing.B) {
	curve := NewJubjub()
	field := curve.field
	ga := &AffinePoint{
		field: field,
		x:     fe(field, "0x11dafe5d23e1218086a365b99fbf3d3be72f6afd7d1f72623e6b071492d1122b"),
		y:     fe(field, "0x1d523cf1ddab1a1793132e78c866c0c33e26ba5cc220fed7cc3f870e59d292aa"),
	}
	ge := ga.ToExtended()
	gp := ga.toProjective()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		curve.double(ge, gp)
	}
}

func BenchmarkCurveMultiplication(t *testing.B) {
	curve := NewJubjub()
	field := curve.field
	scalarField := NewJubjubScalarField()
	ga := &AffinePoint{
		field: field,
		x:     fe(field, "0x11dafe5d23e1218086a365b99fbf3d3be72f6afd7d1f72623e6b071492d1122b"),
		y:     fe(field, "0x1d523cf1ddab1a1793132e78c866c0c33e26ba5cc220fed7cc3f870e59d292aa"),
	}
	s := scalarField.NewRandElement()
	ge := ga.ToExtended()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		curve.Mul(ge, ge, s)
	}
}
