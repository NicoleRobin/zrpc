package cg

import (
	"context"
	"github.com/nicolerobin/zrpc/log"
	"strings"
)

type IfBuilder struct {
	ifBlock   BlockBuilder
	elseBlock BlockBuilder
}

func (i IfBuilder) Body(items ...Builder) IfBuilder {
	i.ifBlock = i.ifBlock.Body(items...)
	return i
}

func (i IfBuilder) AppendBody(items ...Builder) IfBuilder {
	i.ifBlock = i.ifBlock.AppendBody(items...)
	return i
}

func (i IfBuilder) Else(bodies ...Builder) IfBuilder {
	i.elseBlock = BlockBuilder{
		blockType: "else",
		enclosed:  true,
		bodies:    bodies,
	}
	return i
}

func (i IfBuilder) Build() string {
	log.Info(context.Background(), "entrance")
	b := i.ifBlock.Build()

	if i.elseBlock.blockType == "" {
		return b
	}
	return strings.TrimRight(b, "\n") + i.elseBlock.Build()
}

func (i IfBuilder) String() string {
	log.Info(context.Background(), "entrance")
	return i.Build()
}

func If(cond Builder) IfBuilder {
	return IfBuilder{
		ifBlock: BlockBuilder{
			blockType: "if",
			enclosed:  true,
			title:     cond,
		},
	}
}
