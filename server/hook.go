package server

import (
	"context"
	"github.com/nicolerobin/zrpc/core"
)

// BeforeCall execute before api handler called
type BeforeCall func(context.Context, interface{}) (context.Context, error)

// AfterCall execute after api handler called
type AfterCall func(context.Context, interface{}) (context.Context, error)

const (
	Hook core.HookName = "server"
)

var hookTypes = []interface{}{
	BeforeCall(nil),
	AfterCall(nil),
}

var (
	beforeCallHoos []BeforeCall
	afterCallHooks []AfterCall
)

func registerHook(hook interface{}) {
	switch v := hook.(type) {
	case BeforeCall:
		beforeCallHoos = append(beforeCallHoos, v)
	case AfterCall:
		afterCallHooks = append(afterCallHooks, v)
	}
}

func init() {
	core.AddHookRegister(Hook, core.NamedHookRegister(hookTypes, registerHook))
}
