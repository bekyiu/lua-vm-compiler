package state

import (
	"fmt"
	"write_lua/src/binchunk"
	"write_lua/src/vm"
)

// 加载二进制chunk或lua脚本
// 将chunk解析为闭包压入栈顶
// mode: b表示chunk是二进制, t表示chunk是lua脚本, bt表示两者都可
func (this *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk)
	c := newLuaClosure(proto)
	this.stack.push(c)
	// 0表示加载成功
	return 0
}

// 函数调用, nArgs指定被调函数参数的数量, nResults指定被调函数返回值的数量
// 被调函数和他的参数都在栈顶
func (this *luaState) Call(nArgs, nResults int) {
	val := this.stack.get(-(nArgs + 1))
	if c, ok := val.(*closure); ok {
		fmt.Printf("call %s<%d, %d>\n", c.proto.Source,
			c.proto.LineDefined, c.proto.LastLineDefined)
		this.callLuaClosure(nArgs, nResults, c)
	} else {
		panic("not function")
	}
}

func (this *luaState) callLuaClosure(nArgs int, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1
	// 创建调用帧
	newStack := newLuaStack(nRegs + 20)
	newStack.closure = c

	funcAndArgs := this.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	// top以下都是寄存器
	newStack.top = nRegs
	// 记录可变参数
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[1+nParams:]
	}

	this.pushLuaStack(newStack)
	this.runLuaClosure()
	this.popLuaStack()

	// 在被掉帧中获取返回值, 压入当前帧
	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		this.stack.check(len(results))
		this.stack.pushN(results, nResults)
	}
}

// 函数执行
func (this *luaState) runLuaClosure() {
	for {
		ins := vm.Instruction(this.Fetch())
		ins.Execute(this)
		if ins.Opcode() == vm.OP_RETURN {
			break
		}
	}
}
