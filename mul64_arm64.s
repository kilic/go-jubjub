// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !pure_go, arm64

#include "textflag.h"

// This file provides fast assembly versions for the elementary
// arithmetic operations on vectors implemented in arith.go.

// func mul64(x, y uint64) (z1, z0 uint64)
TEXT ·mul64(SB),NOSPLIT,$0
	MOVD	x+0(FP), R0
	MOVD	y+8(FP), R1
	MUL	R0, R1, R2
	UMULH	R0, R1, R3
	MOVD	R3, z1+16(FP)
	MOVD	R2, z0+24(FP)
	RET