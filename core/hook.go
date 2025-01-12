package core

import "reflect"

type HookName string

type HookRegister func(interface{})

type HookRegisterInfo struct {
	types        []interface{}
	hookRegister HookRegister
}

func (h HookRegisterInfo) register() HookRegister {
	return h.hookRegister
}

func (h HookRegisterInfo) typeMap() map[string]reflect.Type {
	d := make(map[string]reflect.Type)
	for _, v := range h.types {
		typeOf := reflect.TypeOf(v)
		name := typeOf.Name()
		d[name] = typeOf
	}
	return d
}

type registerEntry struct {
	register HookRegister
	typeMap  map[string]reflect.Type
}

type hookRegisterItem interface {
	register() HookRegister
	typeMap() map[string]reflect.Type
}

func NamedHookRegister(hookTypes []interface{}, register HookRegister) HookRegisterInfo {
	return HookRegisterInfo{
		types:        hookTypes,
		hookRegister: register,
	}
}

var (
	hookRegisterMap = make(map[HookName]registerEntry)
)

// AddHookRegister add hook registers
func AddHookRegister(name HookName, item hookRegisterItem) {
	entry := registerEntry{
		register: item.register(),
		typeMap:  item.typeMap(),
	}
	hookRegisterMap[name] = entry
}

func AddHook(name HookName, hook interface{}) {
	register(hook)
}

func register(hook interface{}) {
}

func init() {
	for hookName, itemRegisterEntry := range hookRegisterMap {
		_ = hookName
		_ = itemRegisterEntry
	}
}
