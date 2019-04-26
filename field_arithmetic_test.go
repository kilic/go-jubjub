package jubjub

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"
)

type hexstr string

func (s hexstr) uints4() [4]uint64 {
	if s[:2] == "0x" {
		s = s[2:]
	}
	r := new([4]uint64)
	for i := 0; i < 4; i++ {
		bs, _ := hex.DecodeString(string(s[(3-i)*16 : (4-i)*16]))
		r[i] = binary.BigEndian.Uint64(bs)
	}
	return *r
}

func (s hexstr) uints8() [8]uint64 {
	if s[:2] == "0x" {
		s = s[2:]
	}
	r := new([8]uint64)
	for i := 0; i < 8; i++ {
		bs, _ := hex.DecodeString(string(s[(7-i)*16 : (8-i)*16]))
		r[i] = binary.BigEndian.Uint64(bs)
	}
	return *r
}

func hexString(e []uint64) (r string) {
	for i := len(e); i > 0; i-- {
		r = r + fmt.Sprintf("%16.16x", e[i-1])
	}
	return
}

func TestSinglePrecisionMultiplication(t *testing.T) {
	var a uint64 = 0xaabbccdd1111111c
	var b uint64 = 0xff000000aaaaaaab
	var ehi uint64 = 0xaa111110a5d2889e
	var elo uint64 = 0x7d9f4fb1b05b05b4
	hi, lo := mul64(a, b)
	if ehi != hi {
		t.Errorf("single precision multiplication fails for carry limb\n, have %x, want %x", hi, ehi)
	}
	if elo != lo {
		t.Errorf("single precision multiplication fails for first limb\n, have %x, want %x", lo, elo)
	}
}

func TestMultiPrecisionMultiplication(t *testing.T) {
	var a hexstr = "0x6aaaaaaaaaaaaaaa44aa44aa44aa44aa91919191ffffff0000119999ffaa01aa"
	var b hexstr = "0x4fffffffffffffffdddddddd8888888888888888000000119a9a9a9a8b8b8b8b"
	var e hexstr = "0x2155555555555555273c51e6d9036e58e950083f02ccfc1de185217fc929459fc0467f564e5ba377b716c90e4c730d1a2078bc4c1a3544fdf1ed68e95584354e"
	var c [8]uint64
	mul256(&c, a.uints4(), b.uints4())
	if c != e.uints8() {
		t.Errorf("multi precision multiplication fails for first limb\n, have %s, want %s", hexString(c[:]), e)
	}
}

func TestSinglePrecisionSquaring(t *testing.T) {
	var a uint64 = 0xaabbccdd1111111c
	var ehi uint64 = 0x71ddf5da8992abc1
	var elo uint64 = 0xdac76fc0fedcbb10
	hi, lo := square64(a)
	if ehi != hi {
		t.Errorf("single precision multiplication fails for carry limb\n, have %x, want %x", hi, ehi)
	}
	if elo != lo {
		t.Errorf("single precision multiplication fails for first limb\n, have %x, want %x", lo, elo)
	}
}

func TestMultiPrecisionSquaring(t *testing.T) {
	var a hexstr = "0x6aaaaaaaaaaaaaaa44aa44aa44aa44aa91919191ffffff0000119999ffaa01aa"
	var e hexstr = "0x2c71c71c71c71c71721c1cc6c771721c85812e25b2d02ff41c33f3811d3fa1bce41719ece0c67ff282b914d2dc9b6186872d471e6e9a10d1071b516ae1cac4e4"
	b := [8]uint64{}
	square256(&b, a.uints4())
	if e.uints8() != b {
		t.Errorf("multi precision multiplication fails for carry limb\n, have %s, want %s", hexString(b[:]), e)
	}
}

func BenchmarkSinglePrecisionMultiplication(t *testing.B) {
	var a uint64 = 0xaabbccdd1111111c
	var b uint64 = 0xff777777aaaaaaab
	var c, d uint64
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		c, d = mul64(a, b)
	}
	_, _ = c, d
}

func BenchmarkMultiPrecisionMultiplication1(t *testing.B) {
	var a hexstr = "0x6aaaaaaaaaaaaaaa44aa44aa44aa44aa91919191ffffff0000119999ffaa01aa"
	var b hexstr = "0x4fffffffffffffffdddddddd8888888888888888000000119a9a9a9a8b8b8b8b"
	ua := a.uints4()
	ub := b.uints4()
	w := new([8]uint64)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		mul256(w, ua, ub)
	}
	_, _ = ua, ub
}

func BenchmarkSquaring(t *testing.B) {
	var a hexstr = "0x6aaaaaaaaaaaaaaa44aa44aa44aa44aa91919191ffffff0000119999ffaa01aa"
	ua := a.uints4()
	w := [8]uint64{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		square256(&w, ua)
	}
	_ = w
}
