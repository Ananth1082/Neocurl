package compiler

type CollectionNode struct {
	Name     string
	Children []childrenNode
}
type childrenNode interface {
	getName() string
}

type Filename string

type HttpVerb string

type SetNode struct {
	Name      string
	Variables map[string]string
}

func (s SetNode) getName() string {
	return s.Name
}

type RequestNode struct {
	Name string
	Verb HttpVerb
	In   []Filename
	Head map[string]string
	Body string
}

func (r RequestNode) getName() string {
	return r.Name
}
