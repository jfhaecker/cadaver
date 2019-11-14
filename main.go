package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
)

type Blob struct {
	content   []byte
	ofilename string
	nfilename string
	hashCode  []byte
}

type Tree struct {
	ofilename string
	nfilename string
	blobs     []Blob
	subtrees  []Tree
}

func (b *Blob) blob() {
	h := sha256.New()
	h.Write([]byte(b.content))
	b.hashCode = h.Sum(nil)
	b.content = append([]byte("blob\n"), b.content...)
	b.nfilename = fmt.Sprintf("%x", b.hashCode)
}

func main() {
	hashObject([]byte("hallo"), "1.txt")
	hashObject([]byte("tolles"), "2.txt")
	hashObject([]byte("programm"), "3.txt")
}

func hashObject(content []byte, filename string) {

	b := Blob{content: content, ofilename: filename}
	b.blob()

	fmt.Printf("%v\n", b.nfilename)
	ioutil.WriteFile(b.nfilename, b.content, 0644)

}
