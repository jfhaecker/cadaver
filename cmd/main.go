package main

import (
	"bytes"
	"cadaver/internal/lib"
	"fmt"
	"os"
)

var (
	index   = make(map[string]lib.GitObject)
	workDir = "work"
)

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
	b := lib.Blob{Content: content, FilePath: path}
	index[path] = &b
}

func commit() {
	fmt.Println("===========Commit")
	tree := &lib.Tree{}
	tree.Children = make([]lib.GitObject, 0)
	for _, obj := range index {
		obj.DoHash()
		file, _ := os.Open(workDir + "/" + obj.ID().Hex())
		obj.Store(file)
		tree.Children = append(tree.Children, obj)
	}
	tree.DoHash()
	file, _ := os.Open(workDir + "/" + tree.ID().Hex())
	tree.Store(file)

	fmt.Printf("Tree ID=%v\n", tree.ID().Hex())
	var b bytes.Buffer
	tree.Store(&b)
	fmt.Printf(b.String())

}
