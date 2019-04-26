package jubjub

import (
	"crypto/rand"
	"testing"
)

var nBox = 100000

func TestBoxFieldELementByteInOut(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	var a, b = new(FieldElement), new(FieldElement)
	bytes := make([]byte, 32)
	for i := 0; i < nBox; i++ {
		field.RandElement(a, rand.Reader)
		a.Marshal(bytes)
		field.Mul(a, a, field.r2)
		b = field.NewElement(bytes)
		if !b.Eq(a) {
			t.Errorf("bad byte conversion in:%s, out:%s",
				a.String(), b.String())
		}
	}
}

func TestBoxAdditiveAssoc(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	var a, b, c, u, v FieldElement
	for i := 0; i < nBox; i++ {
		field.RandElement(&a, rand.Reader)
		field.RandElement(&b, rand.Reader)
		field.RandElement(&c, rand.Reader)
		field.Add(&u, &a, &b)
		field.Add(&u, &u, &c)
		field.Add(&v, &b, &c)
		field.Add(&v, &v, &a)
		if !u.Eq(&v) {
			t.Errorf("additive associativity does not hold a:%s, b:%s, c:%s, u:%s, v:%s",
				a.String(), b.String(), c.String(), u.String(), v.String())
		}
	}
}

func TestBoxSubractiveAssoc(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	var a, b, c, u, v FieldElement
	for i := 0; i < nBox; i++ {
		field.RandElement(&a, rand.Reader)
		field.RandElement(&b, rand.Reader)
		field.RandElement(&c, rand.Reader)
		field.Sub(&u, &a, &c)
		field.Sub(&u, &u, &b)
		field.Sub(&v, &a, &b)
		field.Sub(&v, &v, &c)
		if !u.Eq(&v) {
			t.Errorf("subtractive associativity does not hold a:%s, b:%s, c:%s, u:%s, v:%s",
				a.String(), b.String(), c.String(), u.String(), v.String())
		}
	}
}

func TestBoxMultiplicativeAssoc(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	var a, b, c, u, v FieldElement
	for i := 0; i < nBox; i++ {
		field.RandElement(&a, rand.Reader)
		field.RandElement(&b, rand.Reader)
		field.RandElement(&c, rand.Reader)
		field.Mul(&u, &a, &b)
		field.Mul(&u, &u, &c)
		field.Mul(&v, &b, &c)
		field.Mul(&v, &v, &a)
		if !u.Eq(&v) {
			t.Errorf("multiplicative associativity does not hold a:%s, b:%s, c:%s, u:%s, v:%s",
				a.String(), b.String(), c.String(), u.String(), v.String())
		}
	}
}

func TestBoxAdditiveCommutativity(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	var a, b, u, v FieldElement
	for i := 0; i < nBox; i++ {
		field.RandElement(&a, rand.Reader)
		field.RandElement(&b, rand.Reader)
		field.Add(&u, &a, &b)
		field.Add(&v, &b, &a)
		if !u.Eq(&v) {
			t.Errorf("additive commutativity  does not hold a:%s, b:%s, u:%s",
				a.String(), b.String(), u.String())
		}
	}
}

func TestBoxMultiplicativeCommutativity(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	var a, b, u, v FieldElement
	for i := 0; i < nBox; i++ {
		field.RandElement(&a, rand.Reader)
		field.RandElement(&b, rand.Reader)
		field.Mul(&u, &a, &b)
		field.Mul(&v, &b, &a)
		if !u.Eq(&v) {
			t.Errorf("multiplicative commutativity does not hold a:%s, b:%s, u:%s",
				a.String(), b.String(), u.String())
		}
	}
}

func TestBoxNegation(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	var a, b, u, v FieldElement
	for i := 0; i < nBox; i++ {
		field.RandElement(&a, rand.Reader)
		field.RandElement(&b, rand.Reader)
		field.Sub(&u, &a, &b)
		field.Neg(&a, &a)
		field.Neg(&b, &b)
		field.Sub(&v, &b, &a)
		if !u.Eq(&v) {
			t.Errorf("subtraction check does not hold a:%s, b:%s, u:%s",
				a.String(), b.String(), u.String())
		}
	}
}

func TestBoxNegation2(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	var a, b, c FieldElement
	var zero = &FieldElement{0, 0, 0, 0}
	for i := 0; i < nBox; i++ {
		field.RandElement(&a, rand.Reader)
		field.Neg(&b, &a)
		field.Add(&c, &a, &b)
		if !zero.Eq(&c) {
			t.Errorf("bad negation a:%s, b:%s",
				a.String(), b.String())
		}
	}
}

func TestBoxDoubling(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	var a, b, c, monttwo FieldElement
	field.Mul(&monttwo, &FieldElement{2, 0, 0, 0}, field.r2)
	for i := 0; i < nBox; i++ {
		field.RandElement(&a, rand.Reader)
		field.Double(&b, &a)
		field.Mul(&c, &a, &monttwo)
		if !b.Eq(&c) {
			t.Errorf("bad doubling c:%s, b:%s",
				c.String(), b.String())
		}
	}
}

func TestBoxAdditiveIdentity(t *testing.T) {

	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	var a, c FieldElement
	identity := &FieldElement{0, 0, 0, 0}
	for i := 0; i < nBox; i++ {
		field.RandElement(&a, rand.Reader)
		field.Add(&c, &a, identity)
		if !c.Eq(&a) {
			t.Errorf("additive identity does not hold, have: %s, want: %s", c.String(), a.String())
		}
		field.Add(&a, &c, identity)
		if !c.Eq(&a) {
			t.Errorf("additive identity does not hold, have: %s, want: %s", c.String(), a.String())
		}
	}
}

func TestBoxMultiplicativeIdentity(t *testing.T) {

	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	var a, c FieldElement
	identity := field.r1
	for i := 0; i < nBox; i++ {
		field.RandElement(&a, rand.Reader)
		field.Mul(&c, &a, identity)
		if !c.Eq(&a) {
			t.Errorf("multiplicative identity does not hold, have: %s, want: %s", c.String(), a.String())
		}
		field.Mul(&a, &c, identity)
		if !c.Eq(&a) {
			t.Errorf("multiplicative identity does not hold, have: %s, want: %s", c.String(), a.String())
		}
	}
}

func TestBoxInverseDown(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	var a, b, c FieldElement
	e := &FieldElement{1, 0, 0, 0}
	for i := 0; i < nBox; i++ {
		field.RandElement(&a, rand.Reader)
		field.InvMontDown(&b, &a)
		field.Mul(&c, &b, &a)
		if !c.Eq(e) {
			t.Errorf("bad montgomery downgrade inversion have: %s, want: %s", c.String(), e.String())
		}
	}
}

func TestBoxInverse(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	var a, b, c FieldElement
	e := field.r1
	for i := 0; i < nBox; i++ {
		field.RandElement(&a, rand.Reader)
		field.InvMontUp(&b, &a)
		field.Mul(&c, &b, &a)
		if !c.Eq(e) {
			t.Errorf("bad montgomery upgrade inversion, have: %s, want: %s", c.String(), e.String())
		}
	}
}

func TestBoxSquare(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	var a, b, c FieldElement
	for i := 0; i < nBox; i++ {
		field.RandElement(&a, rand.Reader)
		field.Square(&b, &a)
		field.Mul(&c, &a, &a)
		if !c.Eq(&b) {
			t.Errorf("bad squaring, have: %s, want: %s", c.String(), b.String())
		}
	}
}
