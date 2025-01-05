package cg

import (
	"context"
	"github.com/nicolerobin/zrpc/log"
	"strings"
)

type RawString string

func (s RawString) Build() string {
	log.Info(context.Background(), "entrance")
	return string(s)
}

func (s RawString) String() string {
	log.Info(context.Background(), "entrance")
	return s.Build()
}

func (s RawString) Attr(attr string) RawString {
	return S(string(s) + "." + attr)
}

func (s RawString) Call(items ...string) RawString {
	return S(string(s) + "(" + strings.Join(items, ", ") + ")")
}

func (s RawString) NewLine() Builder {
	return S(string(s) + "\n")
}

// Addr builds a raw string builder from string
func Addr(t string) RawString {
	return S("&" + t)
}

// S builds a raw string builder from string
func S(t string) RawString {
	return RawString(t)
}

// B builds a bool value
func B(v bool) RawString {
	s := "false"
	if v {
		s = "true"
	}
	return S(s)
}

func Paren(d Builder) RawString {
	return S("(" + d.Build() + ")")
}

func Assign(x RawString, y Builder) RawString {
	return S(string(x) + " = " + y.Build())
}

func Assert(x, t Builder) RawString {
	return S(x.Build() + ".(" + t.Build() + ")")
}

func Attr(t string) RawString {
	return S("&" + t)
}
