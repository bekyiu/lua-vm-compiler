package state

import . "write_lua/src/api"

// 往指定表里设置值
func (this *luaState) setTable(t, k, v luaValue) {
	if tbl, ok := t.(*luaTable); ok {
		tbl.put(k, v)
		return
	}
	panic("not a table")
}

// 往索引指定的表中设置k, v
func (this *luaState) SetTable(idx int) {
	t := this.stack.get(idx)
	v := this.stack.pop()
	k := this.stack.pop()
	this.setTable(t, k, v)
}

func (this *luaState) SetField(idx int, k string) {
	t := this.stack.get(idx)
	v := this.stack.pop()
	this.setTable(t, k, v)
}

func (this *luaState) SetI(idx int, i int64) {
	t := this.stack.get(idx)
	v := this.stack.pop()
	this.setTable(t, i, v)
}

// 向全局环境写入指定值，值在栈顶
func (this *luaState) SetGlobal(name string) {
	t := this.registry.get(LUA_RIDX_GLOBALS)
	v := this.stack.pop()
	this.setTable(t, name, v)
}

// 向全局环境注册go函数
func (this *luaState) Register(name string, f GoFunction) {
	this.PushGoFunction(f)
	this.SetGlobal(name)
}