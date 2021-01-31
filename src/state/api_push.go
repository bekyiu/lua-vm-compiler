package state

import . "write_lua/src/api"

func (this *luaState) PushNil() {
	this.stack.push(nil)
}

func (this *luaState) PushBoolean(b bool) {
	this.stack.push(b)
}
func (this *luaState) PushInteger(n int64) {
	this.stack.push(n)
}
func (this *luaState) PushNumber(n float64) {
	this.stack.push(n)
}
func (this *luaState) PushString(s string) {
	this.stack.push(s)
}
// 把go函数转换为go闭包, 然后压栈
func (this *luaState) PushGoFunction(f GoFunction) {
	this.stack.push(newGoClosure(f))
}
// 把全局环境压栈
func (this *luaState) PushGlobalTable() {
	global := this.registry.get(LUA_RIDX_GLOBALS)
	this.stack.push(global)
}