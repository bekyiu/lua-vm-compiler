package state

// 返回当前pc
func (this *luaState) PC() int {
	return this.pc
}

// 修改pc
func (this *luaState) AddPC(n int) {
	this.pc += n
}

// 取指令
func (this *luaState) Fetch() uint32 {
	ins := this.proto.Codes[this.pc]
	this.pc++
	return ins
}

// 将常量表索引处的常量推入栈顶
func (this *luaState) GetConst(idx int) {
	c := this.proto.Constants[idx]
	this.stack.push(c)
}

// 将指定常量或栈值推入栈顶
func (this *luaState) GetRK(rk int) {
	// 常量表索引
	if rk > 0xFF {
		// 抹去最高位得到常量表的索引
		this.GetConst(rk & 0xFF)
	} else { // 寄存器索引
		// lua虚拟机寄存器索引从0开始
		this.PushValue(rk + 1)
	}
}
