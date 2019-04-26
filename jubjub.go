package jubjub

import (
	"encoding/hex"
)

func fe(field *Field, s string) *FieldElement {
	if s[:2] == "0x" {
		s = s[2:]
	}
	h, _ := hex.DecodeString(s)
	if field == nil {
		return new(FieldElement).Unmarshal(h)
	}
	return field.NewElement(h)
}

func NewJubjub() *Curve {
	field := NewField(bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001"))
	d := fe(field, "0x2a9318e74bfa2b48f5fd9207e6bd7fd4292d7f6d37579d2601065fd6d6343eb1")
	d2 := fe(field, "0x552631ce97f45691ebfb240fcd7affa8525afeda6eaf3a4c020cbfadac687d62")
	// TODO fix generator
	generator := &ExtendedPoint{
		x: &FieldElement{0, 0, 0, 0},
		y: &FieldElement{1, 0, 0, 0},
		t: &FieldElement{0, 0, 0, 0},
		z: &FieldElement{1, 0, 0, 0}}
	return &Curve{
		d:         d,
		twoD:      d2,
		field:     field,
		generator: generator,
	}
}

func NewJubjubScalarField() *ScalarField {
	q := bigFromStr16("0x0e7db4ea6533afa906673b0101343b00a6682093ccc81082d0970e5ed6f72cb7")
	scalarField := &ScalarField{
		q: q,
	}
	return scalarField
}

// var qStr = "0xe7db4ea6533afa906673b0101343b00a6682093ccc81082d0970e5ed6f72cb7"
// var pStr = "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001"
// var qMinus2Str = "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfefffffffeffffffff"
// var inpStr = "0x3d443ab0d7bf2839181b2c170004ec0653ba5bfffffe5bfdfffffffeffffffff"
// var rN1Str = "0x1bbe869330009d577204078a4f77266aab6fca8f09dc705f13f75b69fe75c040"
// var r2Str = "0x0748d9d99f59ff1105d314967254398f2b6cedcb87925c23c999e990f3f29c6d"
// var r3Str = "0x6e2a5bb9c8db33e973d13c71c7b5f4181b3e0d188cf06990c62c1807439b73af"
// var oneDStr = "0x2a9318e74bfa2b48f5fd9207e6bd7fd4292d7f6d37579d2601065fd6d6343eb1"
// var twoDStr = "0x552631ce97f45691ebfb240fcd7affa8525afeda6eaf3a4c020cbfadac687d62"
// var aStr = "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000000"
