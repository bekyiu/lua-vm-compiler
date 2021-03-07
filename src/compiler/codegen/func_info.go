package codegen

import "write_lua/src/vm"

// 函数编译结果的内部结构
type funcInfo struct {
	constants map[interface{}]int    // 常量表; k: 常量值, v: 常量在表中的索引
	usedRegs  int                    // 已分配寄存器的数量
	maxRegs   int                    // 寄存器个数
	scopeLv   int                    // 当前作用域层次, 从0开始
	locVars   []*locVarInfo          // 顺序记录函数内部声明的局部变量
	locNames  map[string]*locVarInfo // 当前生效的局部变量
	breaks    [][]int                // 记录循环块内待处理的跳转指令
	parent    *funcInfo              // 方便定位到外围函数
	upvalues  map[string]upvalInfo   // 存放upvalue表
	insts     []uint32               // 编码后的指令
	subFuncs  []*funcInfo            // 存放子函数信息
	numParams int                    // 参数个数
	isVararg  bool                   // 是否是可变参数
}

// 返回常量在表中的索引, 索引从0开始
func (this *funcInfo) indexOfConstant(k interface{}) int {
	if idx, found := this.constants[k]; found {
		return idx
	}
	idx := len(this.constants)
	this.constants[k] = idx
	return idx
}

// 返回一个空闲的寄存器索引, 索引从0开始
func (this *funcInfo) allocReg() int {
	this.usedRegs++
	if this.usedRegs >= 255 {
		panic("函数或表达式使用寄存器数量达到上限")
	}
	// 更新最大寄存器数量
	if this.usedRegs > this.maxRegs {
		this.maxRegs = this.usedRegs
	}
	return this.usedRegs - 1
}

// 回收最近分配的寄存器
func (this *funcInfo) freeReg() {
	this.usedRegs--
}

// 连续分配n个寄存器, 返回第一个寄存器的索引
func (this *funcInfo) allocRegs(n int) int {
	for i := 0; i < n; i++ {
		this.allocReg()
	}
	return this.usedRegs - n
}

// 回收最近分配的n个寄存器
func (this *funcInfo) freeRegs(n int) {
	for i := 0; i < n; i++ {
		this.freeReg()
	}
}

// 局部变量信息
type locVarInfo struct {
	prev     *locVarInfo // 前驱节点, 串联同名局部变量
	name     string      // 局部变量名
	scopeLv  int         // 局部变量所在作用域层次
	slot     int         // 局部变量名绑定的寄存器索引
	captured bool        // 是否被闭包捕获
}

// 闭包捕获的局部变量
type upvalInfo struct {
	locVarSlot int // 如果闭包捕获的是直接外围函数的局部变量, 该字段表示 该局部变量所占寄存器索引
	upvalIndex int // 如果捕获的是外围函数的upvalue, 该字段表示 该upvalue在外围函数upvalue表中的索引
	index      int // 记录upvalue在表中的索引
}

// 返回upvalue在表中的索引
func (this *funcInfo) indexOfUpval(name string) int {
	if upval, ok := this.upvalues[name]; ok {
		return upval.index
	}
	// 尝试绑定name和upvalue
	if this.parent != nil {
		if locVar, found := this.parent.locNames[name]; found {
			idx := len(this.upvalues)
			this.upvalues[name] = upvalInfo{
				locVarSlot: locVar.slot,
				upvalIndex: -1,
				index:      idx,
			}
			locVar.captured = true
			return idx
		}
		if uvIdx := this.parent.indexOfUpval(name); uvIdx >= 0 {
			idx := len(this.upvalues)
			this.upvalues[name] = upvalInfo{
				locVarSlot: -1,
				upvalIndex: uvIdx,
				index:      idx,
			}
			return idx
		}
	}
	return -1
}

// 进入新的作用域
func (this *funcInfo) enterScope(breakable bool) {
	this.scopeLv++
	// breaks 长度就是块的深度
	if breakable {
		// 循环块
		this.breaks = append(this.breaks, []int{})
	} else {
		this.breaks = append(this.breaks, nil)
	}
}

// 把break语句对应的跳转指令添加到最近的循环块中
func (this *funcInfo) addBreakJmp(pc int) {
	for i := this.scopeLv; i >= 0; i-- {
		if this.breaks[i] != nil {
			this.breaks[i] = append(this.breaks[i], pc)
			return
		}
	}
	panic("break error")
}

// 向当前作用域添加一个局部变量, 返回分配的寄存器索引
func (this *funcInfo) addLocVar(name string) int {
	newVar := &locVarInfo{
		prev:     this.locNames[name],
		name:     name,
		scopeLv:  this.scopeLv,
		slot:     this.allocReg(),
		captured: false,
	}
	this.locVars = append(this.locVars, newVar)
	this.locNames[name] = newVar
	return newVar.slot
}

// 返回局部变量绑定寄存器的索引
func (this *funcInfo) slotOfLocVar(name string) int {
	if locVar, found := this.locNames[name]; found {
		return locVar.slot
	}
	return -1
}

// 解绑局部变量, 回收寄存器
func (this *funcInfo) removeLocVar(locVar *locVarInfo) {
	this.freeReg()
	if locVar.prev == nil {
		delete(this.locNames, locVar.name)
	} else if locVar.prev.scopeLv == locVar.scopeLv {
		// 同一作用域内存在同名局部变量
		this.removeLocVar(locVar.prev)
	} else {
		// 同名局部变量在更外层的作用域
		// 重新与外层的局部变量绑定
		this.locNames[locVar.name] = locVar.prev
	}
}

// 退出作用域
func (this *funcInfo) exitScope() {
	// 待修正的跳转指令
	pendingBreakJmps := this.breaks[len(this.breaks)-1]
	this.breaks = this.breaks[:len(this.breaks)-1]
	a := this.getJmpArgA()
	for _, pc := range pendingBreakJmps {
		sbx := this.pc() - pc
		// 修正指令
		i := (sbx+vm.MAXARG_sBx)<<14 | a<<6 | vm.OP_JMP
		this.insts[pc] = uint32(i)
	}
	this.scopeLv--
	for _, locVar := range this.locNames {
		if locVar.scopeLv > this.scopeLv {
			this.removeLocVar(locVar)
		}
	}
}

//func (self *funcInfo) getJmpArgA() int {
//	hasCapturedLocVars := false
//	minSlotOfLocVars := self.maxRegs
//	for _, locVar := range self.locNames {
//		if locVar.scopeLv == self.scopeLv {
//			for v := locVar; v != nil && v.scopeLv == self.scopeLv; v = v.prev {
//				if v.captured {
//					hasCapturedLocVars = true
//				}
//				if v.slot < minSlotOfLocVars && v.name[0] != '(' {
//					minSlotOfLocVars = v.slot
//				}
//			}
//		}
//	}
//	if hasCapturedLocVars {
//		return minSlotOfLocVars + 1
//	} else {
//		return 0
//	}
//}

// 用于生成四种模式的指令
func (this *funcInfo) emitABC(opcode, a, b, c int) {
	i := b<<23 | c<<14 | a<<6 | opcode
	this.insts = append(this.insts, uint32(i))
}

func (this *funcInfo) emitABx(opcode, a, bx int) {
	i := bx<<14 | a<<6 | opcode
	this.insts = append(this.insts, uint32(i))
}

func (this *funcInfo) emitAsBx(opcode, a, b int) {
	i := (b+vm.MAXARG_sBx)<<14 | a<<6 | opcode
	this.insts = append(this.insts, uint32(i))
}

func (this *funcInfo) emitAx(opcode, ax int) {
	i := ax<<6 | opcode
	this.insts = append(this.insts, uint32(i))
}

// 返回已经生成的最后一条指令的索引
func (this *funcInfo) pc() int {
	return len(this.insts) - 1
}

func (this *funcInfo) fixSbx(pc, sBx int) {
	i := this.insts[pc]
	i = i << 18 >> 18                     // 清除sbx操作数
	i = i | uint32(sBx+vm.MAXARG_sBx)<<14 // 重置sbx操作数
	this.insts[pc] = i
}
