package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var (
	newline = "\n"
	index   = make(map[string]GitObject)
	workDir = "work"
)

type GitObject interface {
	Type() []byte
	DoHash()
	ID() Hashcode
	Store()
	Path() string
	//createFileContent()
}

type Tree struct {
	path     string
	children []GitObject
	content  []byte
	hashCode Hashcode
}

func (t *Tree) Store() {
	t.createFileContent()
	ioutil.WriteFile(workDir+"/"+t.ID().Hex(), t.content, 0644)
}

func (t *Tree) DoHash() {
	code := []byte{}
	for _, child := range t.children {
		code = append(code, child.ID().Array()...)
	}
	t.hashCode = ComputeHash(t.Type(), code)
}

func (b *Tree) ID() Hashcode {
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
	content  []byte
	path     string
	hashCode Hashcode
}

func (b *Blob) createFileContent() []byte {
	filecontent := []byte{}
	filecontent = append(b.Type(), []byte(newline)...)
	filecontent = append(filecontent, b.content...)
	return filecontent
}

func (b Blob) Type() []byte {
	return []byte("blob")
}

func (b *Blob) DoHash() {
	b.hashCode = ComputeHash(b.Type(), b.content)
}

func (b *Blob) Path() string {
	return b.path
}

func (b *Blob) ID() Hashcode {
	return b.hashCode
}

func (b *Blob) Store() {
	filecontent := b.createFileContent()
	ioutil.WriteFile(workDir+"/"+b.ID().Hex(), filecontent, 0644)
}

func main() {

	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		os.Mkdir(workDir, 0755)
	}

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
		obj.Store()
		tree.children = append(tree.children, obj)
	}
	tree.DoHash()
	tree.Store()
	fmt.Printf("Tree : %v\n", tree.ID().Hex())
	fmt.Printf("Tree : %v\n", string(tree.content))
}
