package jubjub

import (
	"encoding/binary"
	"math/big"
	"reflect"
	"testing"
)

func TestFieldElementFromOneByte(t *testing.T) {
	bytes := []byte{
		1,
	}
	a := new(FieldElement).Unmarshal(bytes)
	e := &FieldElement{1, 0, 0, 0}
	if *a != *e {
		t.Errorf("cannot unmarshal bytes have: %x, want %x", a, e)
	}
}

func TestFieldElementFromOneBytes1(t *testing.T) {
	bytes := []byte{
		255, 255,
	}
	a := new(FieldElement).Unmarshal(bytes)
	e := &FieldElement{0xffff, 0, 0, 0}
	if *a != *e {
		t.Errorf("cannot unmarshal bytes have: %s, want %s", a.String(), e.String())
	}
}

func TestFieldElementFromOneBytes2(t *testing.T) {
	bytes := []byte{
		255, 255, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	a := new(FieldElement).Unmarshal(bytes)
	e := &FieldElement{0, 0xffff, 0, 0}
	if *a != *e {
		t.Errorf("cannot unmarshal bytes have: %s, want %s", a.String(), e.String())
	}
}

func TestFieldElementFromOneBytes3(t *testing.T) {
	x0, x64, x128, x192 := uint64(0x11bbbbbbbbbbbb22), uint64(0x33bbbbbbbbbbbb44),
		uint64(0x55bbbbbbbbbbbb66), uint64(0x77bbbbbbbbbbbb88)
	bytes := make([]byte, 32)
	binary.BigEndian.PutUint64(bytes[:], x0)
	binary.BigEndian.PutUint64(bytes[8:], x64)
	binary.BigEndian.PutUint64(bytes[16:], x128)
	binary.BigEndian.PutUint64(bytes[24:], x192)
	a := new(FieldElement).Unmarshal(bytes)
	e := &FieldElement{x192, x128, x64, x0}
	if *a != *e {
		t.Errorf("cannot unmarshal bytes have: %s, want %s", a.String(), e.String())
	}
}

func TestFieldELementToBytes1(t *testing.T) {
	bytes := []byte{
		100, 200,
	}
	a := new(FieldElement).Unmarshal(bytes)
	out := make([]byte, 32)
	a.Marshal(out)
	zeros := make([]byte, 30)
	if !reflect.DeepEqual(zeros[:], out[:30]) || !reflect.DeepEqual(bytes[:], out[30:]) {
		t.Errorf("cannot marshal field element")
	}
}

func TestFieldELementToBytes2(t *testing.T) {
	in := []byte{
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x10, 0x20,
		0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80, 0x90,
		0xfa, 0xfa, 0xfa, 0xfa, 0x24, 0x23, 0x22, 0x21,
	}
	a := new(FieldElement).Unmarshal(in)
	out := make([]byte, 32)
	a.Marshal(out)
	if !reflect.DeepEqual(in[:], out[:]) {
		t.Errorf("cannot marshal field element, have: %x, want: %x", out, in)
	}
}

func TestFieldELementToBytes3(t *testing.T) {
	in := []byte{
		0x11, 0x22, 0x33, 0x55, 0x66, 0x77, 0x88,
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x10, 0x20,
		0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80, 0x90,
		0xfa, 0xfa, 0xfa, 0xfa, 0x24, 0x23, 0x22, 0x21,
	}
	a := new(FieldElement).Unmarshal(in)
	out := make([]byte, 32)
	a.Marshal(out)
	if !reflect.DeepEqual(in[:], out[1:]) || out[0] != 0 {
		t.Errorf("cannot marshal field element, have: %x, want: %x", out, in)
	}
}

func TestFieldAddition(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(field, "0x6aaaaaaaaaaaaaaa44aa44aa44aa44aa91919191ffffff0000119999ffaa01aa")
	b := fe(field, "0x4fffffffffffffffdddddddd8888888888888888000000119a9a9a9a8b8b8b8b")
	e := fe(field, "0x46bd0357810d2d61ef4e4a7fc390f52dc65c76170001a3129aac34358b358d34")
	c := &FieldElement{}
	field.Add(c, a, b)
	if !e.Eq(c) {
		t.Errorf("field element addition fails, have %s, want %s", c.String(), e.String())
	}
}

func TestFieldDoubling(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(field, "0x6aaaaaaaaaaaaaaa44aa44aa44aa44aa91919191ffffff0000119999ffaa01aa")
	e := fe(field, "0x6167ae022bb7d80c561ab14c7fb2b14fcf657f210001a20100233334ff540353")
	c := &FieldElement{}
	field.Double(c, a)
	if !e.Eq(c) {
		t.Errorf("field element doubling fails, have %s, want %s", c.String(), e.String())
	}
}

func TestFieldDoublingZero(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(field, "0x00")
	c := &FieldElement{}
	field.Double(c, a)
	if !a.Eq(c) {
		t.Errorf("field element doubling fails, have %s, want %s", c.String(), a.String())
	}
}

func TestFieldSubtraction(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(field, "0x6aaaaaaaaaaaaaaa44aa44aa44aa44aa91919191ffffff0000119999ffaa01aa")
	b := fe(field, "0x4fffffffffffffffdddddddd8888888888888888000000119a9a9a9a8b8b8b8b")
	e := fe(field, "0x5942fca87ef2d29dcc6d713b4d801be34ab49af8fffe5d109a8900ff8be189e2")
	c := &FieldElement{}
	field.Sub(c, b, a)
	if !e.Eq(c) {
		t.Errorf("field element subtraction fails, have %s, want %s", c.String(), e.String())
	}
	e = fe(field, "0x1aaaaaaaaaaaaaaa66cc66ccbc21bc2209090909fffffeee6576feff741e761f")
	field.Sub(c, a, b)
	if !e.Eq(c) {
		t.Errorf("field element subtraction fails, have %s, want %s", c.String(), e.String())
	}
}

func TestFieldNegate(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(nil, "0x6aaaaaaaaaaaaaaa44aa44aa44aa44aa91919191ffffff0000119999ffaa01aa")
	e := fe(nil, "0x0942fca87ef2d29dee8f935dc4f7935ac22c1270fffe5cfeffee66650055fe57")
	c := &FieldElement{}
	field.Neg(c, a)
	if !e.Eq(c) {
		t.Errorf("field element negation fails, have %s, want %s", c.String(), e.String())
	}
}

func TestMontgomeryReduction(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	T := [8]uint64{
		0x22bbccdd55558888,
		0xaabbccdd55558888,
		0xaabbccdd55558888,
		0x11bbccdd55558888,
		0x22bbccdd55558888,
		0xaabbccdd55558888,
		0xaabbccdd55558888,
		0x22bbccdd55558888,
	}
	e := fe(nil, "0x0ac1b4094057dae42dab79d6693ee71d832ffa2bb7648e3884a7d38f035dceed")
	r := new(FieldElement)
	field.montReduce(r, T)
	if !r.Eq(e) {
		t.Errorf("montgomerry reduction fails, have %s, want %s", r.String(), e.String())
	}
}

func TestFieldModularReduction(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(nil, "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	e := fe(nil, "0x1824b159acc5056f998c4fefecbc4ff55884b7fa0003480200000001fffffffd")
	field.Mul(a, a, field.r1)
	if !a.Eq(e) {
		t.Errorf("modular reduction fails, have %s, want %s", a.String(), e.String())
	}
}

func TestMontgomeryReductionZero(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	T := [8]uint64{0, 0, 0, 0, 0, 0, 0, 0}
	e := &FieldElement{0, 0, 0, 0}
	r := new(FieldElement)
	field.montReduce(r, T)
	if !r.Eq(e) {
		t.Errorf("montgomerry reduction fails, have %s, want %s", r.String(), e.String())
	}
}

func TestMontgomeryReductionOne(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	T := [8]uint64{1, 0, 0, 0, 0, 0, 0, 0}
	e := field.rN1
	r := new(FieldElement)
	field.montReduce(r, T)
	if !r.Eq(e) {
		t.Errorf("montgomerry reduction fails, have %s, want %s", r.String(), e.String())
	}
}

func TestMontgomeryForm(t *testing.T) {
	p, _ := new(big.Int).SetString("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001"[2:], 16)
	field := NewField(p)
	a := fe(field, "0x6aaaaaaaaaaaaaaa44aa44aa44aa44aa91919191ffffff0000119999ffaa01aa")
	e := fe(nil, "0x6aaaaaaaaaaaaaaa44aa44aa44aa44aa91919191ffffff0000119999ffaa01aa")
	one := &FieldElement{1, 0, 0, 0}
	field.Mul(a, a, one)
	if !a.Eq(e) {
		t.Errorf("field element addition fails, have %s, want %s", a.String(), e.String())
	}
}

func TestFieldElementMontgomeryMultiplication(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(field, "0x1aaaaaaaaaaaaaa8ccccccccccccccc8777777777777777dffffffffddddddd0")
	b := fe(field, "0x6cccccccccccccc911111111555555559393939393939393ffffffffeeeeeeef")
	e := fe(field, "0x678f0b264343979944fb3663d336c345b4347ed20629e7de98aa44cc4aa681bb")
	c := new(FieldElement)
	field.Mul(c, a, b)
	if !c.Eq(e) {
		t.Errorf("mont multiplication fails, have %s, want %s", c.String(), e.String())
	}
}

func TestFieldElementMontgomeryMultiplicationOne(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(field, "0x01")
	b := fe(field, "0x6cccccccccccccc911111111555555559393939393939393ffffffffeeeeeeef")
	e := b
	c := new(FieldElement)
	field.Mul(c, a, b)
	if !c.Eq(e) {
		t.Errorf("mont multiplication fails, have %s, want %s", c.String(), e.String())
	}
}

func TestFieldElementMontgomeryMultiplicationZero(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := &FieldElement{0, 0, 0, 0}
	b := fe(field, "0x6cccccccccccccc911111111555555559393939393939393ffffffffeeeeeeef")
	e := a
	c := new(FieldElement)
	field.Mul(c, a, b)
	if !c.Eq(e) {
		t.Errorf("mont multiplication fails, have %s, want %s", c.String(), e.String())
	}
}

func TestFieldElementMontgomerySquare(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(field, "0x1aaaaaaaaaaaaaa8ccccccccccccccc8777777777777777dffffffffddddddd0")
	e := fe(field, "0x39b25f1073641793dcb0007efc52a272b06514bd57562d45292d141976190f71")
	c := new(FieldElement)
	field.Square(c, a)
	if !c.Eq(e) {
		t.Errorf("mont multiplication fails, have %s, want %s", c.String(), e.String())
	}
}

func TestFieldElementMontgomerySquareOne(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(field, "0x01")
	e := field.r1
	c := new(FieldElement)
	field.Square(c, a)
	if !c.Eq(e) {
		t.Errorf("mont squaring one fails, have %s, want %s", c.String(), e.String())
	}
	a = fe(nil, "0x01")
	e = field.rN1
	c = new(FieldElement)
	field.Square(c, a)
	if !c.Eq(e) {
		t.Errorf("mont squaring one fails, have %s, want %s", c.String(), e.String())
	}
}

func TestFieldElementMontgomerySquareZero(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := &FieldElement{0, 0, 0, 0}
	e := fe(field, "0x00")
	c := new(FieldElement)
	field.Square(c, a)
	if !c.Eq(e) {
		t.Errorf("mont squaring one fails, have %s, want %s", c.String(), e.String())
	}
}

func TestFieldInverseEuclid(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(nil, "0x66ffeeeeddddccccffffeeeeddddcccc99aa99aa88bb88bb1919191928282828")
	e := fe(nil, "0x27c28f49cabcf02ec28a6a44d07436e062d004894bffeeefa73ab2abc10f487f")
	inv := new(FieldElement)
	field.InvEEA(inv, a)
	if !inv.Eq(e) {
		t.Errorf("inversion fails (euclid), have %s, want %s", inv.String(), e.String())
	}
}

func TestFieldInverseMontgomeryDown(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(field, "0x66ffeeeeddddccccffffeeeeddddcccc99aa99aa88bb88bb1919191928282828")
	e := fe(nil, "0x27c28f49cabcf02ec28a6a44d07436e062d004894bffeeefa73ab2abc10f487f")
	inv := new(FieldElement)
	field.InvMontDown(inv, a)
	if !inv.Eq(e) {
		t.Errorf("inversion fails (montgomery down), have %s, want %s", inv.String(), e.String())
	}
	// also
	a = fe(nil, "0x66ffeeeeddddccccffffeeeeddddcccc99aa99aa88bb88bb1919191928282828")
	e = fe(field, "0x27c28f49cabcf02ec28a6a44d07436e062d004894bffeeefa73ab2abc10f487f")
	inv = new(FieldElement)
	field.InvMontDown(inv, a)
	if !inv.Eq(e) {
		t.Errorf("inversion fails (montgomery down), have %s, want %s", inv.String(), e.String())
	}
}

func TestFieldInverseMontgomeryUp(t *testing.T) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(field, "0x66ffeeeeddddccccffffeeeeddddcccc99aa99aa88bb88bb1919191928282828")
	e := fe(field, "0x27c28f49cabcf02ec28a6a44d07436e062d004894bffeeefa73ab2abc10f487f")
	inv := new(FieldElement)
	field.InvMontUp(inv, a)
	if !inv.Eq(e) {
		t.Errorf("inversion fails (montgomery up), have %s, want %s", inv.String(), e.String())
	}
}

func BenchmarkFieldAddition(t *testing.B) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(nil, "0x6aaaaaaaaaaaaaaa44aa44aa44aa44aa91919191ffffff0000119999ffaa01aa")
	b := fe(nil, "0x4fffffffffffffffdddddddd8888888888888888000000119a9a9a9a8b8b8b8b")
	c := new(FieldElement)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		field.Add(c, a, b)
	}
}

func BenchmarkFieldSubtraction(t *testing.B) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(nil, "0x6aaaaaaaaaaaaaaa44aa44aa44aa44aa91919191ffffff0000119999ffaa01aa")
	b := fe(nil, "0x4fffffffffffffffdddddddd8888888888888888000000119a9a9a9a8b8b8b8b")
	c := new(FieldElement)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		field.Sub(c, a, b)
	}
}

func BenchmarkFieldMontgomeryReduction(t *testing.B) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	T := [8]uint64{
		0x22bbccdd55558888,
		0xaabbccdd55558888,
		0xaabbccdd55558888,
		0x11bbccdd55558888,
		0x22bbccdd55558888,
		0xaabbccdd55558888,
		0xaabbccdd55558888,
		0x22bbccdd55558888,
	}
	var b = new(FieldElement)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		field.montReduce(b, T)
	}
	_ = b
}

func BenchmarkFieldMontgomeryMultiplication(t *testing.B) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(nil, "0x6aaaaaaaaaaaaaaa44aa44aa44aa44aa91919191ffffff0000119999ffaa01aa")
	b := fe(nil, "0x4fffffffffffffffdddddddd8888888888888888000000119a9a9a9a8b8b8b8b")
	c := new(FieldElement)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		field.Mul(c, a, b)
	}
}

func BenchmarkFieldMontgomerySquaring(t *testing.B) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(nil, "0x6aaaaaaaaaaaaaaa44aa44aa44aa44aa91919191ffffff0000119999ffaa01aa")
	c := new(FieldElement)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		field.Square(c, a)
	}
}

func BenchmarkFieldInverse1(t *testing.B) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(nil, "0x1aaaaaaaaaaaaaa8ccccccccccccccc8777777777777777dffffffffddddddd0")
	inv := new(FieldElement)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		field.InvEEA(inv, a)
	}
}

func BenchmarkFieldInverse2(t *testing.B) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(nil, "0x1aaaaaaaaaaaaaa8ccccccccccccccc8777777777777777dffffffffddddddd0")
	inv := new(FieldElement)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		field.InvMontDown(inv, a)
	}
}

func BenchmarkFieldInverse3(t *testing.B) {
	p := bigFromStr16("0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	field := NewField(p)
	a := fe(nil, "0x1aaaaaaaaaaaaaa8ccccccccccccccc8777777777777777dffffffffddddddd0")
	inv := new(FieldElement)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		field.InvMontUp(inv, a)
	}
}
