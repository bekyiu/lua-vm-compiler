package vm

import . "write_lua/src/api"

// R(A) -= R(A + 2); pc += sbx
func forPrep(ins Instruction, vm LuaVM) {
	a, sbx := ins.AsBx()
	a += 1
	// R(A) -= R(A + 2)
	vm.PushValue(a)
	vm.PushValue(a + 2)
	vm.Arith(LUA_OPSUB)
	vm.Replace(a)
	// pc += sbx
	vm.AddPC(sbx)
}

/*
A: start index, A + 1: limit, A + 2: step, A + 3: i

R(A) += R(A + 2)

步长>=0, <?=的意思是 <=, 否则为 >=

if R(A) <?= R(A + 1) then {
	pc += sbx
	R(A + 3) = R(A)
}
*/
func forLoop(ins Instruction, vm LuaVM) {
	a, sbx := ins.AsBx()
	a += 1
	// R(A) += R(A + 2)
	vm.PushValue(a + 2)
	vm.PushValue(a)
	vm.Arith(LUA_OPADD)
	vm.Replace(a)
	// 判断步长是否为正数
	isPositive := vm.ToNumber(a+2) >= 0
	// R(A) <?= R(A + 1)
	if (isPositive && vm.Compare(a, a+1, LUA_OPLE)) ||
		(!isPositive && vm.Compare(a+1, a, LUA_OPLE)) {
		vm.AddPC(sbx)
		vm.Copy(a, a+3)
	}

}

func tForLoop(ins Instruction, vm LuaVM) {
	a, sbx := ins.AsBx()
	a += 1
	if !vm.IsNil(a + 1) {
		vm.Copy(a+1, a)
		vm.AddPC(sbx)
	}
}