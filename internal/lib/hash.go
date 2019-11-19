package lib

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"strconv"
)

type Hashcode [sha1.Size]byte

func (h Hashcode) Hex() string {
	s := hex.EncodeToString(h[:])
	return s
}

func (h Hashcode) Array() []byte {
	return h[:]
}

type Hasher struct {
	hash.Hash
}

func ComputeHash(objtype []byte, content []byte) (h Hashcode) {
	hcalc := NewHasher(objtype, content, int64(len(content)))
	copy(h[:], hcalc.Hash.Sum(nil))
	return
}

func NewHasher(objtype []byte, content []byte, size int64) Hasher {
	h := Hasher{sha1.New()}
	h.Write(objtype)
	h.Write([]byte(" "))
	h.Write([]byte(strconv.FormatInt(size, 10)))
	h.Write([]byte{0})
	h.Write(content)
	return h
}
