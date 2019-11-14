package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
)

var (
	newline = "\n"
	//index   = make([]GitObject, 0)
	index = make(map[string]GitObject)
)

type GitObject interface {
	Type() []byte
	DoHash()
	GetHashCode() []byte
	Store()
	GetOFilename() string
	CreateFileContent()
}

type Tree struct {
	ofilename string
	nfilename string
	children  []GitObject
	content   []byte
	hashCode  []byte
}

func (t *Tree) Store() {
	ioutil.WriteFile(t.nfilename, t.content, 0644)
}

func (t *Tree) DoHash() {
	h := sha256.New()
	for _, child := range t.children {
		h.Write(child.GetHashCode())
	}
	t.hashCode = h.Sum(nil)
	t.nfilename = fmt.Sprintf("%x", t.hashCode)

}

func (t Tree) Type() []byte {
	return []byte("tree")
}

type Blob struct {
	content     []byte
	ofilename   string
	nfilename   string
	hashCode    []byte
	filecontent []byte
}

func (b *Blob) CreateFileContent() {
	b.filecontent = []byte{}
	header := []byte{}
	header = append(b.Type(), []byte(newline)...)
	b.filecontent = append(header, b.content...)
	b.nfilename = fmt.Sprintf("%x", b.hashCode)
}

func (b Blob) Type() []byte {
	return []byte("blob")
}

func (b *Blob) DoHash() {
	h := sha256.New()
	h.Write([]byte(b.content))
	b.hashCode = h.Sum(nil)
}

func (b *Blob) GetOFilename() string {
	return b.ofilename
}

func (b *Blob) GetHashCode() []byte {
	return b.hashCode
}

func (b *Blob) Store() {
	//fmt.Printf("---> %v Storing: [%v]\n", b.nfilename, string(b.content))
	ioutil.WriteFile(b.nfilename, b.filecontent, 0644)

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

func add(content []byte, filename string) {
	b := Blob{content: content, ofilename: filename}
	index[filename] = &b
	//index = append(index, &b)
	//fmt.Printf("%v\n", b.nfilename)
}

func commit() {
	fmt.Println("===========Commit")
	tree := &Tree{}
	tree.children = make([]GitObject, 0)
	for _, obj := range index {
		obj.DoHash()
		//fmt.Printf("Hallo:%v\n", obj)
		obj.Store()
		tree.children = append(tree.children, obj)
		//hc := fmt.Sprintf("%x", obj.GetHashCode)
		//fmt.Printf("%v:%x\n", string(obj.Type()), obj.GetHashCode())
	}
	header := []byte{}
	header = append(tree.Type(), []byte(newline)...)
	tree.content = append(tree.content, header...)
	for _, child := range tree.children {
		c := fmt.Sprintf("%v %x %v \n",
			string(child.Type()),
			child.GetHashCode(),
			child.GetOFilename())
		tree.content = append(tree.content, []byte(c)...)
		//fmt.Printf("%v", c)

	}
	tree.DoHash()
	fmt.Printf("Tree : %x\n", string(tree.hashCode))
	fmt.Printf("%v\n", string(tree.content))
	tree.Store()
}
