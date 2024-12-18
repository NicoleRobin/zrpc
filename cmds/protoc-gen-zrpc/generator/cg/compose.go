package cg

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
