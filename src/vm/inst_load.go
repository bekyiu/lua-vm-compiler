package vm

import . "write_lua/src/api"

// 给连续n个局部变脸(寄存器)设置nil
func loadNil(ins Instruction, vm LuaVM) {
	a, b, _ := ins.ABC()
	// 寄存器索引
	a += 1
	vm.PushNil()
	for i := a; i <= a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

// 给单个寄存器设置bool值
// R(A) := (bool) B; if(C) pc++
func loadBool(ins Instruction, vm LuaVM) {
	a, b, c := ins.ABC()
	vm.PushBoolean(b != 0)
	vm.Replace(a + 1)
	if c != 0 {
		vm.AddPC(1)
	}
}

// 将常量表的某个常量加载到指定寄存器
// R(A) := Kst(Bx)
// bx 18bit, 所以常量表的最大索引就是2^18 - 1
func loadK(ins Instruction, vm LuaVM) {
	a, bx := ins.ABx()
	vm.GetConst(bx)
	vm.Replace(a + 1)
}

// 结合extraage指令, 拓展常量表的索引为2^26 - 1
func loadKx(ins Instruction, vm LuaVM) {
	a, _ := ins.ABx()
	ax := Instruction(vm.Fetch()).Ax()
	vm.GetConst(ax)
	vm.Replace(a + 1)
}