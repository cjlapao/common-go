package cryptorand

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

type CryptoSource struct{}

func (s CryptoSource) Seed(seed int64) {}

func (s CryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s CryptoSource) Uint64() (v uint64) {
	err := binary.Read(rand.Reader, binary.BigEndian, &v)
	if err != nil {
		v = uint64(time.Now().Unix())
	}
	return v
}
