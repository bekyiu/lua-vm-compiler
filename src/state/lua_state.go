package state

import . "write_lua/src/api"

// 封装了lua解释器的状态
type luaState struct {
	stack    *luaStack // lua栈
	registry *luaTable // 注册表, 提供给用户
}

func New() *luaState {
	registry := newLuaTable(0, 0)
	// 全局环境
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(0, 0))
	ls := &luaState{registry: registry}
	ls.pushLuaStack(newLuaStack(LUA_MINSTACK, ls))
	return ls
}

// 压入一个调用栈帧, this.stack 永远指向栈顶的栈帧, 是链表的头部
func (this *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = this.stack
	this.stack = stack
}

// 弹出一个调用帧
func (this *luaState) popLuaStack() {
	stack := this.stack
	this.stack = stack.prev
	stack.prev = nil
}
