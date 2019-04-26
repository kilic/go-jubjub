package jubjub

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestAffinePointMontgomeryDecode(t *testing.T) {
	curve := NewJubjub()
	field := curve.field
	g1 := &AffinePoint{
		field: field,
		x:     fe(field, "0x11dafe5d23e1218086a365b99fbf3d3be72f6afd7d1f72623e6b071492d1122b"),
		y:     fe(field, "0x1d523cf1ddab1a1793132e78c866c0c33e26ba5cc220fed7cc3f870e59d292aa"),
	}
	x, _ := hex.DecodeString("11dafe5d23e1218086a365b99fbf3d3be72f6afd7d1f72623e6b071492d1122b")
	y, _ := hex.DecodeString("1d523cf1ddab1a1793132e78c866c0c33e26ba5cc220fed7cc3f870e59d292aa")
	out := make([]byte, 64)
	g1.MontgomeryDecodeUncompressed(out)

	if !bytes.Equal(x, out[:32]) {
		t.Errorf("bad x coordinate, have: %x, want: %x", out[:32], x)
	}
	if !bytes.Equal(x, out[:32]) {
		t.Errorf("bad y coordinate, have: %x, want: %x", out[32:], y)
	}

}

func TestAffinePointFromBytes(t *testing.T) {
	curve := NewJubjub()
	in := make([]byte, 64)
	x, _ := hex.DecodeString("11dafe5d23e1218086a365b99fbf3d3be72f6afd7d1f72623e6b071492d1122b")
	y, _ := hex.DecodeString("1d523cf1ddab1a1793132e78c866c0c33e26ba5cc220fed7cc3f870e59d292aa")
	copy(in[:32], x[:])
	copy(in[32:64], y[:])
	_, err := curve.NewAffinePointFromUncompressed(in)
	if err != nil {
		t.Errorf("cannot unmarshal new point")
	}
}

func TestAffinePointFromBytes2(t *testing.T) {
	curve := NewJubjub()
	in := make([]byte, 64)
	out := make([]byte, 64)
	x, _ := hex.DecodeString("11dafe5d23e1218086a365b99fbf3d3be72f6afd7d1f72623e6b071492d1122b")
	y, _ := hex.DecodeString("1d523cf1ddab1a1793132e78c866c0c33e26ba5cc220fed7cc3f870e59d292aa")
	copy(in[:32], x[:])
	copy(in[32:64], y[:])
	point, _ := curve.NewAffinePointFromUncompressed(in)
	point.MontgomeryDecodeUncompressed(out)
	if !bytes.Equal(in, out) {
		t.Errorf("cannot decode coorditates have: %x, want: %x", out, in)
	}
}
