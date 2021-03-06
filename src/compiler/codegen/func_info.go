package codegen

// 函数编译结果的内部结构
type funcInfo struct {
	constants map[interface{}]int    // 常量表; k: 常量值, v: 常量在表中的索引
	usedRegs  int                    // 已分配寄存器的数量
	maxRegs   int                    // 寄存器个数
	scopeLv   int                    // 当前作用域层次, 从0开始
	locVars   []*locVarInfo          // 顺序记录函数内部声明的局部变量
	locNames  map[string]*locVarInfo // 当前生效的局部变量

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

// 进入新的作用域
func (this *funcInfo) enterScope() {
	this.scopeLv++
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