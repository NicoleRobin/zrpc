package cg

import "strings"

// FuncBuilder defines a builder for function.
type FuncBuilder struct {
	receiverName string
	receiverType Builder

	block BlockBuilder

	funcName string

	returnType []Builder
	params     []ParamBuilder
	call       []string

	isType      bool
	isInterface bool
	isMethod    bool
}

func (m FuncBuilder) Attr(n string) RawString {
	return S(m.receiverName + "." + n)
}

func (m FuncBuilder) Name(n string) FuncBuilder {
	m.funcName = n
	return m
}

func (m FuncBuilder) Body(items ...Builder) FuncBuilder {
	m.block = m.block.Body(items...)
	return m
}

func (m FuncBuilder) Param(items ...ParamBuilder) FuncBuilder {
	m.params = items
	return m
}

func (m FuncBuilder) Return(t ...Builder) FuncBuilder {
	m.returnType = t
	return m
}

func (m FuncBuilder) Call(items ...string) FuncBuilder {
	if items == nil {
		items = []string{}
	}
	m.call = items
	return m
}

func (m FuncBuilder) AsType() FuncBuilder {
	m.isType = true
	return m
}

func (m FuncBuilder) Build() string {
	var s string
	if !m.isInterface {
		s += "func "
	}
	if m.isMethod {
		s += "(" + m.receiverName + " " + m.receiverType.Build() + ") "
	}
	s += m.funcName + "("
	for _, p := range m.params {
		if p == NoParam {
			continue
		}

		s += p.Build() + ","
	}
	s += ")"

	if n := len(m.returnType); n != 0 {
		s += " "
		if n > 1 {
			s += "("
		}
		for i, t := range m.returnType {
			s += t.Build()
			if i != n-1 {
				s += ","
			}
		}

		if n > 1 {
			s += ")"
		}
	}

	if !m.isType {
		s += m.block.Build()
		if m.call != nil {
			s += "(" + strings.Join(m.call, ", ") + ")"
		}
	}

	if !m.isType && m.funcName != "" {
		s += "\n"
	}
	return s
}

func (m FuncBuilder) String() string {
	return m.Build()
}

func Func(name string) FuncBuilder {
	return FuncBuilder{
		funcName: name,
		block: BlockBuilder{
			enclosed: true,
		},
	}
}

func Method(receiverName string, receiverType Builder) FuncBuilder {
	b := Func("")
	b.receiverName = receiverName
	b.receiverType = receiverType
	b.isMethod = true
	return b
}

func InterfaceAPI(name string) FuncBuilder {
	b := Func(name)
	b.isInterface = true
	b.isType = true
	return b
}

func Defer(b Builder) Builder {
	return S("defer " + b.Build() + "\n")
}

func New(b Builder) Builder {
	return S("new() " + b.Build() + ")")
}
