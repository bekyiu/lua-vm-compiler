package vm

import (
	. "write_lua/src/api"
)

// 把当前闭包的某个upvalue值拷贝到目标寄存器中
// R(A) := Upvalue[B]
func getUpval(ins Instruction, vm LuaVM) {
	a, b, _ := ins.ABC()
	a += 1
	b += 1
	vm.Copy(LuaUpvalueIndex(b), a)
}

// 使用寄存器中的值给当前闭包的upvalue赋值
// Upvalue[B] := R(A)
func setUpval(ins Instruction, vm LuaVM) {
	a, b, _ := ins.ABC()
	a += 1
	b += 1
	vm.Copy(a, LuaUpvalueIndex(b))
}

// 如果当前闭包的某个upvalue是表, 则根据键从该表里取值
// R(A) := Upvalue[B][RK(C)]
func getTabUp(ins Instruction, vm LuaVM) {
	a, b, c := ins.ABC()
	a += 1
	b += 1
	vm.GetRK(c)
	vm.GetTable(LuaUpvalueIndex(b))
	vm.Replace(a)
}

// 如果当前闭包的某个upvalue是表, 则根据键往表里写值
// Upvalue[A][RK(B)] := RK(C)
func setTabUp(ins Instruction, vm LuaVM) {
	a, b, c := ins.ABC()
	a += 1
	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(LuaUpvalueIndex(a))
}
