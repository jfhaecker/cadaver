package lib

import "io"

type GitObject interface {
	Type() []byte
	DoHash()
	ID() Hashcode
	Store(io.Writer)
	Path() string
}
