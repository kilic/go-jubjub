package jubjub

import (
	"encoding/hex"
	"math/big"
)

var (
	big1   = new(big.Int).SetUint64(1)
	big2   = new(big.Int).SetUint64(2)
	big3   = new(big.Int).SetUint64(3)
	big64  = new(big.Int).SetUint64(64)
	big256 = new(big.Int).SetUint64(256)
)

// var zeroBigInt = new(big.Int).SetInt64(0)
// var oneBigInt = new(big.Int).SetInt64(0)
// var d256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)

func bn() *big.Int {
	return new(big.Int)
}

func bigFromStr10(s string) *big.Int {
	n, _ := new(big.Int).SetString(s, 10)
	return n
}

func bigFromStr16(s string) *big.Int {
	if s[:2] == "0x" {
		s = s[2:]
	}
	n, _ := new(big.Int).SetString(s, 16)
	return n
}

func bigFromInt64(i int64) *big.Int {
	return new(big.Int).SetInt64(i)
}

func toBytes(s string) []byte {
	h, _ := hex.DecodeString(s)
	return h
}

func toUint(s string) [4]uint64 {
	var i int64
	var bigTwo = bigFromInt64(2)
	value := bigFromStr16(s[2:])
	b := bn().Exp(bigTwo, bigFromInt64(64), nil)
	digits := [4]uint64{0, 0, 0, 0}
	for i < 4 {
		digits[i] = bn().Mod(value, b).Uint64()
		value.Div(value, b)
		i++
	}
	return digits
}
