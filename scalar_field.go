package jubjub

import (
	"crypto/rand"
	"math/big"
)

type ScalarField struct {
	q *big.Int
}

func (field *ScalarField) NewElement() *ScalarFieldElement {
	return &ScalarFieldElement{
		n: new(big.Int),
	}
}

func (field *ScalarField) NewRandElement() *ScalarFieldElement {
	e := field.NewElement()
	e.n, _ = rand.Int(rand.Reader, field.q)
	return e
}

func (field *ScalarField) NewElementFromUint64(a uint64) *ScalarFieldElement {
	return new(ScalarFieldElement).set(new(big.Int).SetUint64(a))
}

func (field *ScalarField) NewElementFrom16(s string) *ScalarFieldElement {
	if s[:2] == "0x" {
		s = s[2:]
	}
	n, _ := new(big.Int).SetString(s, 16)
	return new(ScalarFieldElement).set(n)
}

func (fe *ScalarField) NewElementFromBig(b *big.Int) *ScalarFieldElement {
	return new(ScalarFieldElement).set(b)
}

type ScalarFieldElement struct {
	n *big.Int
}

func (e *ScalarFieldElement) set(a *big.Int) *ScalarFieldElement {
	e.n = new(big.Int).Set(a)
	return e
}

func (e *ScalarFieldElement) bit(i int) uint {
	return e.n.Bit(i)
}

func (e *ScalarFieldElement) bitLength() int {
	return e.n.BitLen()
}

func (field *ScalarField) Add(c *ScalarFieldElement, a *ScalarFieldElement, b *ScalarFieldElement) {
	c.n.Add(a.n, b.n)
	c.n.Mod(c.n, field.q)
}

func (field *ScalarField) Mul(c *ScalarFieldElement, a *ScalarFieldElement, b *ScalarFieldElement) {
	c.n.Mul(a.n, b.n)
	c.n.Mod(c.n, field.q)
}

func (field *ScalarField) Double(c *ScalarFieldElement, a *ScalarFieldElement) {
	c.n.Add(a.n, a.n)
	c.n.Mod(c.n, field.q)
}

func (field *ScalarField) Square(c *ScalarFieldElement, a *ScalarFieldElement) {
	c.n.Mul(a.n, a.n)
	c.n.Mod(c.n, field.q)
}

func (field *ScalarField) Exp(c *ScalarFieldElement, a *ScalarFieldElement, e *ScalarFieldElement) {
	c.n.Exp(a.n, e.n, field.q)
}

func (field *ScalarField) Inv(c *ScalarFieldElement, a *ScalarFieldElement) {
	c.n.ModInverse(a.n, field.q)
}
