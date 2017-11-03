package spriter

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/big"
	mathrand "math/rand"
)

type Flipper interface {
	Next() bool
	Int64() int64
}

type SeededFlipper struct {
	*mathrand.Rand
}

func NewSeededFlipper(seed int64) Flipper {
	return &SeededFlipper{mathrand.New(mathrand.NewSource(seed))}
}

func (f *SeededFlipper) Next() bool {
	return f.Float64() > 0.5
}

func (f *SeededFlipper) Int64() int64 {
	return f.Int63()
}

type BytesFlipper struct {
	b big.Int
	i int
}

func NewBytesFlipper(b []byte) Flipper {
	f := &BytesFlipper{}
	f.b.SetBytes(b)
	f.i = 0
	return f
}

func (f *BytesFlipper) Next() bool {
	if f.i >= f.b.BitLen() {
		f.i = 0
	}
	result := f.b.Bit(f.i) > 0
	f.i++
	return result
}

func (f *BytesFlipper) Int64() int64 {
	return int64(binary.BigEndian.Uint64(f.b.Bytes()[:8]))
}

type RandomFlipper struct {
	*BytesFlipper
}

func NewRandomFlipper() Flipper {
	f := NewBytesFlipper(randomBytes(32)).(*BytesFlipper)
	return &RandomFlipper{f}
}

func (f *RandomFlipper) Next() bool {
	if f.BytesFlipper.i >= f.b.BitLen() {
		f.BytesFlipper.b.SetBytes(randomBytes(32))
	}
	return f.BytesFlipper.Next()
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := cryptorand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}
