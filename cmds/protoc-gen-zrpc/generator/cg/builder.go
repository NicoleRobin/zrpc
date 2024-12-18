package cg

import (
	"strconv"
)

type Builder interface {
	Build() string
	String() string
}

// P builds a pointer
func P(t string) Builder {
	return S("*" + t)
}

func DefineAssign(x string, y Builder) Builder {
	return S(x + " := " + y.Build())
}

func Ne(x, y string) Builder {
	return S(x + " != " + y)
}

func Eq(x, y string) Builder {
	return S(x + " == " + y)
}

func Break() Builder {
	return S("break")
}

func BreakLabel(name string) Builder {
	return S("break" + name)
}

func Return(s ...Builder) Builder {
	r := "return "

	for i, v := range s {
		r += v.Build()

		if i != len(s)-1 {
			r += ", "
		}
	}
	return S(r)
}

func Variadic(name string) Builder {
	return S("..." + name)
}

func KV(key string, val Builder) Builder {
	return S(key + ": " + val.Build())
}

func Comment(msg string) Builder {
	return S("// " + msg)
}

func Index(n string, i int) Builder {
	return S(n + "[" + strconv.Itoa(i) + "]")
}

func Recv(r string, b Builder) Builder {
	return S(r + " <- " + b.Build())
}

func Make(t Builder, args ...Builder) Builder {
	s := "make(" + t.Build()

	for _, a := range args {
		s += ", " + a.Build()
	}
	s += ")"

	return S(s)
}

func Chan(t Builder) Builder {
	return S("chan " + t.Build())
}

func Go(t Builder) Builder {
	return S("go " + t.Build())
}

func TypeAlias(name string, tp Builder) Builder {
	return S("type " + name + " " + tp.Build())
}
