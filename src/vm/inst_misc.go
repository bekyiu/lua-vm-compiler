package vm

import . "write_lua/src/api"

// R(A) := R(B)
func move(ins Instruction, vm LuaVM) {
	a, b, _ := ins.ABC()
	vm.Copy(b + 1, a + 1)
}

// 无条件跳转
func jmp(ins Instruction, vm LuaVM) {
	a, sbx := ins.AsBx()
	vm.AddPC(sbx)
	if a != 0 {
		panic("todo...")
	}
}