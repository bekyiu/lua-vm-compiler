package state

// 访问指定索引处的值, 取其长度, 推入栈顶
func (this *luaState) Len(idx int) {
	val := this.stack.get(idx)
	if s, ok := val.(string); ok {
		this.stack.push(int64(len(s)))
	} else if result, ok := callMetamethod(val, val, "__len", this); ok {
		this.stack.push(result)
	} else if t, ok := val.(*luaTable); ok {
		this.stack.push(int64(t.len()))
	} else {
		panic("Len err")
	}
}

// 弹出n个值, 拼接, 放回栈顶
func (this *luaState) Concat(n int) {
	if n == 0 {
		this.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if this.IsString(-1) && this.IsString(-2) {
				s2 := this.ToString(-1)
				s1 := this.ToString(-2)
				this.stack.pop()
				this.stack.pop()
				this.stack.push(s1 + s2)
				continue
			}
			// 不是字符串/数字, 尝试调用元方法
			b := this.stack.pop()
			a := this.stack.pop()
			if result, ok := callMetamethod(a, b, "__concat", this); ok {
				this.stack.push(result)
				continue
			}
			panic("Concat err")
		}
	}
}

// 弹出栈顶的键, 表由索引指定, 压入下一个键值对
func (this *luaState) Next(idx int) bool {
	val := this.stack.get(idx)
	if t, ok := val.(*luaTable); ok {
		key := this.stack.pop()
		if nextKey := t.nextKey(key); nextKey != nil {
			this.stack.push(nextKey)
			this.stack.push(t.get(nextKey))
			return true
		}
		return false
	}
	panic("Next error, 指定索引处不是表")
}

// 从栈顶弹出一个值, 作为异常抛出
func (this *luaState) Error() int {
	err := this.stack.pop()
	panic(err)
}