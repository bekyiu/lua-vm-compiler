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
	curStack := this.stack
	curClosure := curStack.closure
	subProto := curClosure.proto.Protos[idx]
	subClosure := newLuaClosure(subProto)
	curStack.push(subClosure)

	// 根据函数原型的Upvalue表来初始化子闭包的upvalue
	for i, uvInfo := range subProto.Upvalues {

		uvIdx := int(uvInfo.Idx)
		// 等于1说明子闭包的upvalue是直接捕获的当前闭包的局部变量
		if uvInfo.Instack == 1 {
			if curStack.openuvs == nil {
				curStack.openuvs = map[int]*upvalue{}
			}
			if openuv, found := curStack.openuvs[uvIdx]; found {
				subClosure.upvals[i] = openuv
			} else {
				// 此时子闭包upvalue项对应的uvIdx是 当前闭包中被捕获的局部变量在寄存器中的索引
				subClosure.upvals[i] = &upvalue{&curStack.slots[uvIdx]}
				curStack.openuvs[uvIdx] = subClosure.upvals[i]
			}
		} else {
			// 等于0，说明是来源于更外围函数的局部变量，存在于当前闭包的upvalue表中
			// 此时子闭包upvalue项对应的uvIdx是 该upvalue项在当前闭包upvalue表的索引
			subClosure.upvals[i] = curClosure.upvals[uvIdx]
		}
	}
}

func (this *luaState) CloseUpvalues(a int) {
	for i, openuv := range this.stack.openuvs {
		if i >= a - 1 {
			// val是个新的变量, 地址是不一样的
			val := *(openuv.val)
			// 不在指向栈中的局部变量
			openuv.val = &val
			delete(this.stack.openuvs, i)
		}
	}
}