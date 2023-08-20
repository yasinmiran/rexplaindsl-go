package api

const (
	Anchor ConstructType = iota
	Set
	Group
	Literal
	Numeric
	Operator
	Repetition
)

type Template interface {
	ToRegex() string
}

type Expression struct {
	Construct ConstructType
	Instance  interface{}
	ToRegex   func() string
}

type ConstructType int
