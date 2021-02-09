package state

import . "write_lua/src/api"

// 往指定表里设置值, raw表示是否会忽略元方法
func (this *luaState) setTable(t, k, v luaValue, raw bool) {
	if tbl, ok := t.(*luaTable); ok {
		if raw || tbl.get(k) != nil || !tbl.hasMetafield("__newindex") {
			tbl.put(k, v)
			return
		}
	}
	if !raw {
		if mf := getMetafield(t, "__newindex", this); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				this.setTable(x, k, v, false)
				return
			case *closure:
				this.stack.push(mf)
				this.stack.push(t)
				this.stack.push(k)
				this.stack.push(v)
				this.Call(3, 0)
				return
			}
		}
	}
	panic("not a table")
}

// 往索引指定的表中设置k, v
func (this *luaState) SetTable(idx int) {
	t := this.stack.get(idx)
	v := this.stack.pop()
	k := this.stack.pop()
	this.setTable(t, k, v, false)
}

func (this *luaState) SetField(idx int, k string) {
	t := this.stack.get(idx)
	v := this.stack.pop()
	this.setTable(t, k, v, false)
}

func (this *luaState) SetI(idx int, i int64) {
	t := this.stack.get(idx)
	v := this.stack.pop()
	this.setTable(t, i, v, false)
}

// 向全局环境写入指定值，值在栈顶
func (this *luaState) SetGlobal(name string) {
	t := this.registry.get(LUA_RIDX_GLOBALS)
	v := this.stack.pop()
	this.setTable(t, name, v, false)
}

// 向全局环境注册go函数
func (this *luaState) Register(name string, f GoFunction) {
	this.PushGoFunction(f)
	this.SetGlobal(name)
}

// 从栈顶弹出一个表, 并把指定索引处的值的元表设置成该表
func (this *luaState) SetMetatable(idx int) {
	val := this.stack.get(idx)
	mtVal := this.stack.pop()
	if mtVal == nil {
		setMetatable(val, nil, this)
	} else if mt, ok := mtVal.(*luaTable); ok {
		setMetatable(val, mt, this)
	} else {
		panic("SetMetatable error!")
	}
}

func (this *luaState) RawSet(idx int) {
	t := this.stack.get(idx)
	v := this.stack.pop()
	k := this.stack.pop()
	this.setTable(t, k, v, true)
}

func (this *luaState) RawSetI(idx int, i int64) {
	t := this.stack.get(idx)
	v := this.stack.pop()
	this.setTable(t, i, v, true)
}