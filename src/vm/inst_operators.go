package vm

import . "write_lua/src/api"

// R(A) := RK(B) op RK(C)
func _binaryArith(ins Instruction, vm LuaVM, op ArithOp) {
	a, b, c := ins.ABC()
	vm.GetRK(b)
	vm.GetRK(c)
	vm.Arith(op)
	vm.Replace(a + 1)
}

// R(A) := op R(B)
func _unaryArith(ins Instruction, vm LuaVM, op ArithOp) {
	a, b, _ := ins.ABC()
	vm.PushValue(b + 1)
	vm.Arith(op)
	vm.Replace(a + 1)
}

// R(A) := length of R(B)
func length(ins Instruction, vm LuaVM) {
	a, b, _ := ins.ABC()
	vm.Len(b + 1)
	vm.Replace(a + 1)
}

// R(A) := R(B).. ... ..R(C)
func concat(ins Instruction, vm LuaVM) {
	a, b, c := ins.ABC()
	a += 1
	b += 1
	c += 1
	n := c - b + 1
	vm.CheckStack(n)
	for i := b; i <= c; i++ {
		vm.PushValue(i)
	}
	vm.Concat(n)
	vm.Replace(a)
}

// if ((RK(B) op RK(C)) ~= A) then pc++
func _compare(ins Instruction, vm LuaVM, op CompareOp) {
	a, b, c := ins.ABC()
	vm.GetRK(b)
	vm.GetRK(c)
	// (a != 0) 相当于把数字转换为bool
	if vm.Compare(-1, -2, op) != (a != 0) {
		// 条件满足
		// pc++ 是为了跳过下一条jmp指令以执行body
		vm.AddPC(1)
	}
	vm.Pop(2)
}

// R(A) := not R(B)
func not(ins Instruction, vm LuaVM) {
	a, b, _ := ins.ABC()
	vm.PushBoolean(!vm.ToBoolean(b + 1))
	vm.Replace(a + 1)
}

// if (bool(R(B)) == C) then R(A) := R(B) else pc++
func testSet(ins Instruction, vm LuaVM) {
	a, b, c := ins.ABC()
	a += 1
	b += 1
	if vm.ToBoolean(b) == (c != 0) {
		vm.Copy(b, a)
	} else {
		vm.AddPC(1)
	}
}

// if (bool(R(A)) != C) then pc++
func test(ins Instruction, vm LuaVM) {
	a, _, c := ins.ABC()
	if vm.ToBoolean(a) != (c != 0) {
		vm.AddPC(1)
	}
}

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

func add(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPADD) }  // +
func sub(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSUB) }  // -
func mul(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMUL) }  // *
func mod(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMOD) }  // %
func pow(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPPOW) }  // ^
func div(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPDIV) }  // /
func idiv(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPIDIV) } // //
func band(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBAND) } // &
func bor(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPBOR) }  // |
func bxor(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBXOR) } // ~
func shl(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHL) }  // <<
func shr(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHR) }  // >>
func unm(i Instruction, vm LuaVM)  { _unaryArith(i, vm, LUA_OPUNM) }   // -
func bnot(i Instruction, vm LuaVM) { _unaryArith(i, vm, LUA_OPBNOT) }  // ~

func eq(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPEQ) } // ==
func lt(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLT) } // <
func le(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLE) } // <=
