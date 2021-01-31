package state

import . "write_lua/src/api"

// 创建表, 压栈
func (this *luaState) CreateTable(nArr, nRec int) {
	table := newLuaTable(nArr, nRec)
	this.stack.push(table)
}

func (this *luaState) NewTable() {
	this.CreateTable(0, 0)
}

// 从table里取值, 压入栈, 返回值的类型
func (this *luaState) getTable(table, key luaValue) LuaType {
	if tbl, ok := table.(*luaTable); ok {
		val := tbl.get(key)
		this.stack.push(val)
		return typeOf(val)
	}
	panic("not a table")
}
// 从指定索引的表中根据 栈顶的key 取出value, 压栈, 返回值的类型
func (this *luaState) GetTable(idx int) LuaType {
	t := this.stack.get(idx)
	k := this.stack.pop()
	return this.getTable(t, k)
}

// 从指定索引的表中根据 给定的key 取出value, 压栈, 返回值的类型
func (this *luaState) GetField(idx int, k string) LuaType {
	t := this.stack.get(idx)
	return this.getTable(t, k)
}

func (this *luaState) GetI(idx int, i int64) LuaType {
	t := this.stack.get(idx)
	return this.getTable(t, i)
}

// 把全局变量中的指定字段压入栈顶
func (this *luaState) GetGlobal(name string) LuaType {
	t := this.registry.get(LUA_RIDX_GLOBALS)
	return this.getTable(t, name)
}