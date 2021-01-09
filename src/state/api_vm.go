package state

// 返回当前pc
func (this *luaState) PC() int {
	return this.stack.pc
}

// 修改pc
func (this *luaState) AddPC(n int) {
	this.stack.pc += n
}

// 取指令
func (this *luaState) Fetch() uint32 {
	ins := this.stack.closure.proto.Codes[this.stack.pc]
	this.stack.pc++
	return ins
}

// 将常量表索引处的常量推入栈顶
func (this *luaState) GetConst(idx int) {
	c := this.stack.closure.proto.Constants[idx]
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

// 返回当前lua函数的寄存器数量
func (this *luaState) RegisterCount() int {
	return int(this.stack.closure.proto.MaxStackSize)
}

// 把传递给当前lua函数的变长参数推入栈顶 多退少补
func (this *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(this.stack.varargs)
	}
	this.stack.check(n)
	this.stack.pushN(this.stack.varargs, n)
}

// 把当前lua函数的子函数原型实例化为闭包推入栈顶
func (this *luaState) LoadProto(idx int) {
	proto := this.stack.closure.proto.Protos[idx]
	c := newLuaClosure(proto)
	this.stack.push(c)
}
