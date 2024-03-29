package lib

import (
	"fmt"
	"io"
)

type Tree struct {
	path     string
	Children []GitObject
	//Content  []byte
	hashCode Hashcode
}

func (t *Tree) Store(writer io.Writer) {
	writer.Write(t.createFileContent())
	//t.createFileContent()
	//ioutil.WriteFile(workDir+"/"+t.ID().Hex(), t.Content, 0644)
}

func (t *Tree) DoHash() {
	code := []byte{}
	for _, child := range t.Children {
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

func (t *Tree) createFileContent() (header []byte) {
	header = []byte{}
	header = append(t.Type(), []byte("\n")...)
	//t.Content = append(t.Content, header...)
	for _, child := range t.Children {
		c := fmt.Sprintf("%v %x %v \n",
			string(child.Type()),
			child.ID(),
			child.Path())
		header = append(header, []byte(c)...)
	}
	return
}
