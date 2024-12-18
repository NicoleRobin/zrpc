package cg

import "strings"

type BlockBuilder struct {
	blockType  string
	enclosed   bool
	encloser   string
	blockStart string
	delimiter  string
	title      Builder
	bodies     []Builder
}

func (b BlockBuilder) Body(items ...Builder) BlockBuilder {
	b.bodies = items
	return b
}

func (b BlockBuilder) AppendBody(items ...Builder) BlockBuilder {
	bs := make([]Builder, len(b.bodies)+len(items))
	n := copy(bs, b.bodies)
	copy(bs[n:], items)
	b.bodies = bs
	return b
}

func (b BlockBuilder) Build() string {
	s := b.blockType
	if b.title != nil {
		s += " " + b.title.Build()
	}
	if b.blockStart != "" {
		s += " " + b.blockStart
	}

	encloser := "{}"
	if b.encloser != "" {
		encloser = b.encloser
	}
	if b.enclosed {
		s += encloser[0:1]
	}

	delimiter := "\n"
	if b.delimiter != "" {
		delimiter = b.delimiter
	}
	s += "\n"

	for _, item := range b.bodies {
		s += item.Build() + delimiter
	}

	if b.enclosed {
		s = strings.TrimRight(s, "\n")
		s += "\n" + encloser[1:]
	}

	return s
}

func (b BlockBuilder) String() string {
	return b.Build()
}

func Const(items ...Builder) BlockBuilder {
	return BlockBuilder{
		blockType: "const",
		encloser:  "()",
		enclosed:  true,
		bodies:    items,
	}
}

func Array(name string) BlockBuilder {
	return BlockBuilder{
		title:     S(name),
		enclosed:  true,
		blockType: "[]",
		delimiter: ";",
	}
}

func For(cond string) BlockBuilder {
	return BlockBuilder{
		blockType: "for",
		enclosed:  true,
		title:     S(cond),
	}
}

func Label(label string) Builder {
	return S(label + ":")
}

func Select() BlockBuilder {
	return BlockBuilder{
		blockType: "select",
		enclosed:  true,
		title:     S(""),
	}
}

func Case(selector Builder) BlockBuilder {
	return BlockBuilder{
		blockType:  "case",
		title:      selector,
		blockStart: ":",
	}
}

func Default(selector Builder) BlockBuilder {
	return BlockBuilder{
		blockType:  "default",
		title:      S(""),
		blockStart: ":",
	}
}

func Struct(name string) BlockBuilder {
	return BlockBuilder{
		blockType:  "type",
		blockStart: "struct",
		enclosed:   true,
		title:      S(name),
	}
}

func Interface(name string) BlockBuilder {
	b := Struct(name)
	b.blockStart = "interface"
	return b
}

func StructLiteral(name string) BlockBuilder {
	b := Struct(name)
	b.delimiter = ",\n"
	b.blockStart = ""
	b.blockType = ""
	return b
}

func StructPointerLiteral(name string) BlockBuilder {
	b := StructLiteral(name)
	b.blockType = "&"
	return b
}
