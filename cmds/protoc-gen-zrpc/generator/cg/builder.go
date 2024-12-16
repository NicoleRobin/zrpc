package cg

type Builder interface {
	Build() string
	String() string
}

type RawString string

func (s RawString) Build() string {
	return string(s)
}

func (s RawString) String() string {
	return s.Build()
}

type ComposeBuilder []Builder

func (b ComposeBuilder) Build() string {
	var s string
	for _, builder := range b {
		s += builder.Build() + "\n"
	}

	return s
}

func (b ComposeBuilder) String() string {
	return b.Build()
}

type BlockBuilder struct {
	blockType  string
	enclosed   bool
	encloser   string
	blockStart string
	delimiter  string
	title      Builder
	bodies     []Builder
}

func Const(items ...Builder) BlockBuilder {
	return BlockBuilder{
		blockType: "const",
		encloser:  "()",
		enclosed:  true,
		bodies:    items,
	}
}

func S(t string) RawString {
	return RawString(t)
}

func Array(name string) BlockBuilder {
	return BlockBuilder{
		title:     S(name),
		enclosed:  true,
		blockType: "[]",
		delimiter: ";",
	}
}
