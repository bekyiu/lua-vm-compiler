package state

import "write_lua/src/binchunk"
import . "write_lua/src/api"

// 要么是lua闭包, 要么是go闭包
type closure struct {
	proto  *binchunk.Prototype
	goFunc GoFunction
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	return &closure{proto: proto}
}

func newGoClosure(f GoFunction) *closure {
	return &closure{goFunc: f}
}
