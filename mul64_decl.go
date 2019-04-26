// +build !pure_go, amd64 arm64

package jubjub

// implemented in mul64_$GOARCH.s
func mul64(x, y uint64) (z1, z0 uint64)
