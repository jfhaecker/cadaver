package lib

type GitObject interface {
	Type() []byte
	DoHash()
	ID() Hashcode
	Store(workDir string)
	Path() string
	//createFileContent()
}
