package lib

import "io/ioutil"

type Blob struct {
	Content  []byte
	FilePath string
	hashCode Hashcode
}

func (b *Blob) createFileContent() []byte {
	filecontent := []byte{}
	filecontent = append(b.Type(), []byte("\n")...)
	filecontent = append(filecontent, b.Content...)
	return filecontent
}

func (b Blob) Type() []byte {
	return []byte("blob")
}

func (b *Blob) DoHash() {
	b.hashCode = ComputeHash(b.Type(), b.Content)
}

func (b *Blob) Path() string {
	return b.FilePath
}

func (b *Blob) ID() Hashcode {
	return b.hashCode
}

func (b *Blob) Store(workDir string) {
	filecontent := b.createFileContent()
	ioutil.WriteFile(workDir+"/"+b.ID().Hex(), filecontent, 0644)
}
