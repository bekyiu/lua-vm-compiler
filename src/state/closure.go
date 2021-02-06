package state

import "write_lua/src/binchunk"
import . "write_lua/src/api"

// 要么是lua闭包, 要么是go闭包
type closure struct {
	proto  *binchunk.Prototype
	goFunc GoFunction
	upvals []*upvalue
}

type upvalue struct {
	val *luaValue
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	c := &closure{proto: proto}
	if n := len(proto.Upvalues); n > 0 {
		c.upvals = make([]*upvalue, n)
	}
	return c
}

func newGoClosure(f GoFunction, nUpvals int) *closure {
	c := &closure{goFunc: f}
	if nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}
