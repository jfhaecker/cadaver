package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"io/ioutil"
)

var (
	newline = "\n"
	//index   = make([]GitObject, 0)
	index = make(map[string]GitObject)
)

type Hash [sha1.Size]byte

func (h Hash) Hex() string {
	s := hex.EncodeToString(h[:])
	return s
}

func (h Hash) Array() []byte {
	return h[:]
}

type Hasher struct {
	hash.Hash
}

func NewHasher() Hasher {
	h := Hasher{sha1.New()}
	return h
}

func (h Hasher) Write(content []byte) {
	h.Hash.Write(content)
}

func (h Hasher) Sum() (hash Hash) {
	copy(hash[:], h.Hash.Sum(nil))
	return
}

type GitObject interface {
	Type() []byte
	DoHash()
	ID() Hash
	Store()
	Path() string
	//createFileContent()
}

type Tree struct {
	path     string
	children []GitObject
	content  []byte
	hashCode Hash
}

func (t *Tree) Store() {
	t.createFileContent()
	ioutil.WriteFile(t.ID().Hex(), t.content, 0644)
}

func (t *Tree) DoHash() {
	h := NewHasher()
	for _, child := range t.children {
		h.Write(child.ID().Array())
	}
	t.hashCode = h.Sum()
}

func (b *Tree) ID() Hash {
	return b.hashCode
}

func (t Tree) Type() []byte {
	return []byte("tree")
}

func (t *Tree) createFileContent() {
	header := []byte{}
	header = append(t.Type(), []byte(newline)...)
	t.content = append(t.content, header...)
	for _, child := range t.children {
		c := fmt.Sprintf("%v %x %v \n",
			string(child.Type()),
			child.ID(),
			child.Path())
		t.content = append(t.content, []byte(c)...)
	}
}

type Blob struct {
	content []byte
	path    string
	//nfilename   string
	hashCode Hash
	//filecontent []byte
}

func (b *Blob) createFileContent() []byte {
	filecontent := []byte{}
	filecontent = append(b.Type(), []byte(newline)...)
	filecontent = append(filecontent, b.content...)
	//b.nfilename = fmt.Sprintf("%x", b.ID().Hex())
	return filecontent
}

func (b Blob) Type() []byte {
	return []byte("blob")
}

func (b *Blob) DoHash() {
	h := NewHasher()
	h.Write([]byte(b.content))
	b.hashCode = h.Sum()
}

func (b *Blob) Path() string {
	return b.path
}

func (b *Blob) ID() Hash {
	return b.hashCode
}

func (b *Blob) Store() {
	//fmt.Printf("---> %v Storing: [%v]\n", b.nfilename, string(b.content))
	filecontent := b.createFileContent()
	ioutil.WriteFile(b.ID().Hex(), filecontent, 0644)

}

func main() {
	add([]byte("hallo1"), "1.txt")
	add([]byte("hola1"), "2.txt")
	add([]byte("servus1"), "3.txt")
	fmt.Printf("\nIndex: %v\n", len(index))
	commit()
	add([]byte("servus2"), "3.txt")
	fmt.Printf("\nIndex: %v\n", len(index))
	commit()
}

func add(content []byte, path string) {
	b := Blob{content: content, path: path}
	index[path] = &b
}

func commit() {
	fmt.Println("===========Commit")
	tree := &Tree{}
	tree.children = make([]GitObject, 0)
	for _, obj := range index {
		obj.DoHash()
		//fmt.Printf("Hallo:%#v\n", obj)
		obj.Store()
		tree.children = append(tree.children, obj)
		//hc := fmt.Sprintf("%x", obj.GetHashCode)
		//fmt.Printf("%v:%x\n", string(obj.Type()), obj.GetHashCode())
	}
	tree.DoHash()
	fmt.Printf("Tree : %v\n", tree.ID().Hex())
	tree.Store()
}
