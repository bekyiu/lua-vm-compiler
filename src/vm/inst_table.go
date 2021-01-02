package vm

import . "write_lua/src/api"

// R(A) := {} (size = B, C)
// 创建空表
func newTable(ins Instruction, vm LuaVM) {
	a, b, c := ins.ABC()
	a += 1
	// 二进制chunk中b, c采用了浮点字节的编码
	vm.CreateTable(Fb2int(b), Fb2int(c))
	vm.Replace(a)
}

// R(A) := R(B)[RK(C)]
// 根据键从表里取值
func getTable(ins Instruction, vm LuaVM) {
	a, b, c := ins.ABC()
	a += 1
	b += 1
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

// R(A)[RK(B)] := RK(C)
// 根据键往表里设置值
func setTable(ins Instruction, vm LuaVM) {
	a, b, c := ins.ABC()
	a += 1
	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(a)
}

const LFIELDS_PRE_FLUSH = 50

// R(A)[(C - 1) * FPF + i] := R(A + i), 1 <= i <= B
// 只针对数组, 把存在寄存器里的值批量复制到表中
func setList(ins Instruction, vm LuaVM) {
	a, b, c := ins.ABC()
	a += 1
	if c > 0 {
		c -= 1
	} else {
		// 获取批次数
		c = Instruction(vm.Fetch()).Ax()
	}
	// 数组起始索引
	idx := int64(c * LFIELDS_PRE_FLUSH)
	for i := 1; i <= b; i++ {
		vm.PushValue(a + i)
		vm.SetI(a, idx+int64(i))
	}
}
