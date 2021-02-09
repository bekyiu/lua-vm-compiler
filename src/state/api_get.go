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

// 从table里取值, 压入栈, 返回值的类型, raw 表示是否要忽略元方法
func (this *luaState) getTable(table, key luaValue, raw bool) LuaType {
	if tbl, ok := table.(*luaTable); ok {
		val := tbl.get(key)
		if raw || val != nil || !tbl.hasMetafield("__index") {
			this.stack.push(val)
			return typeOf(val)
		}
	}
	if !raw {
		if mf := getMetafield(table, "__index", this); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				// 可能会继续触发__index
				return this.getTable(x, key, false)
			case *closure:
				this.stack.push(mf)
				this.stack.push(table)
				this.stack.push(key)
				this.Call(2, 1)
				v := this.stack.get(-1)
				return typeOf(v)
			}
		}
	}
	panic("not a table")
}

// 从指定索引的表中根据 栈顶的key 取出value, 压栈, 返回值的类型
func (this *luaState) GetTable(idx int) LuaType {
	t := this.stack.get(idx)
	k := this.stack.pop()
	return this.getTable(t, k, false)
}

// 从指定索引的表中根据 给定的key 取出value, 压栈, 返回值的类型
func (this *luaState) GetField(idx int, k string) LuaType {
	t := this.stack.get(idx)
	return this.getTable(t, k, false)
}

func (this *luaState) GetI(idx int, i int64) LuaType {
	t := this.stack.get(idx)
	return this.getTable(t, i, false)
}

// 把全局变量中的指定字段压入栈顶
func (this *luaState) GetGlobal(name string) LuaType {
	t := this.registry.get(LUA_RIDX_GLOBALS)
	return this.getTable(t, name, false)
}

// 看指定索引处的值是否有元表, 如果有, 压入栈顶并返回true
func (this *luaState) GetMetatable(idx int) bool {
	val := this.stack.get(idx)
	if mt := getMetatable(val, this); mt != nil {
		this.stack.push(mt)
		return true
	}
	return false
}

func (this *luaState) RawGet(idx int) LuaType {
	t := this.stack.get(idx)
	k := this.stack.pop()
	return this.getTable(t, k, true)
}

func (this *luaState) RawGetI(idx int, i int64) LuaType {
	t := this.stack.get(idx)
	return this.getTable(t, i, true)
}