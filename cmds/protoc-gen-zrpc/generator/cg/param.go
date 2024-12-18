package cg

// ParamBuilder is used to build parameters
type ParamBuilder struct {
	name      string
	paramType Builder

	value string
	isVar bool
}

var NoParam = ParamBuilder{}

func Param(name string, t Builder) ParamBuilder {
	return ParamBuilder{
		name:      name,
		paramType: t,
	}
}

func (p ParamBuilder) Type(t Builder) ParamBuilder {
	p.paramType = t
	return p
}

func (p ParamBuilder) Name() string {
	return p.name
}

func (p ParamBuilder) Value(m Builder) ParamBuilder {
	p.value = m.Build()
	return p
}

func (p ParamBuilder) Build() string {
	var s string

	if p.isVar {
		s += "var "
	}
	s += p.name
	if p.name != "" {
		s += " "
	}
	if p.paramType != nil {
		s += p.paramType.Build()
	}
	if p.value != "" {
		s += " = " + p.value
	}
	if p.isVar {
		s += "\n"
	}
	return s
}

func (p ParamBuilder) String() string {
	return p.Build()
}
func Var(name string) ParamBuilder {
	return ParamBuilder{
		name:  name,
		isVar: true,
	}
}
