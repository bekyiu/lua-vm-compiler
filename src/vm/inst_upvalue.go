package vm

import . "write_lua/src/api"

// 把某个全局变量放入指定寄存器
func getTabUp(ins Instruction, vm LuaVM) {
	a, _, c := ins.ABC()
	a += 1

	vm.PushGlobalTable()
	vm.GetRK(c)
	vm.GetTable(-2)
	vm.Replace(a)
	vm.Pop(1)
}
