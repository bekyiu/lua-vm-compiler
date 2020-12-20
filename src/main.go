package main

import (
	"fmt"
	"io/ioutil"
	"write_lua/src/binchunk"
	"write_lua/src/vm"
)

func main() {
	data, _ := ioutil.ReadFile("D:\\lua\\lua_code\\ch02\\luac.out")
	proto := binchunk.Undump(data)
	list(proto)
}

func list(f *binchunk.Prototype) {
	if f != nil {
		printHeader(f)
		printCode(f)
		printDetail(f)

		for _, proto := range f.Protos {
			list(proto)
		}
	}
}

func printHeader(f *binchunk.Prototype) {
	funcType := "main"
	if f.LineDefined > 0 {
		funcType = "function"
	}
	varargFlag := ""
	if f.IsVararg > 0 {
		varargFlag = "+"
	}
	fmt.Printf(" \n%s <%s:%d,%d>(%d instructions)\n", funcType, f.Source, f.LineDefined, f.LastLineDefined, len(f.Codes))
	fmt.Printf("%d%s params, %d slots, %d upvalues, ", f.NumParams, varargFlag, f.MaxStackSize, len(f.Upvalues))
	fmt.Printf("%d locals, %d constants, %d functions\n", len(f.LocVars), len(f.Constants), len(f.Protos))
}

func printCode(f *binchunk.Prototype) {
	for pc, code := range f.Codes {
		line := "-"
		if len(f.LineInfos) > 0 {
			line = fmt.Sprintf("%d", f.LineInfos[pc])
		}
		ins := vm.Instruction(code)
		// 序号, 行号, 指令名
		fmt.Printf("\t%d\t[%s]\t%s\t", pc+1, line, ins.OpName())
		printOperands(ins)
		fmt.Println()
	}
}

func printOperands(ins vm.Instruction) {
	switch ins.OpMode() {
	case vm.IABC:
		a, b, c := ins.ABC()

		fmt.Printf("%d", a)
		if ins.BMode() != vm.OpArgN {
			// 第9位是1, 常量表索引
			if b > 0xFF {
				// b&0xff得到正数的索引
				fmt.Printf(" %d", -1-(b&0xFF))
			} else {
				fmt.Printf(" %d", b)
			}
		}
		if ins.CMode() != vm.OpArgN {
			if c > 0xFF {
				fmt.Printf(" %d", -1-(c&0xFF))
			} else {
				fmt.Printf(" %d", c)
			}
		}
	case vm.IABx:
		a, bx := ins.ABx()

		fmt.Printf("%d", a)
		if ins.BMode() == vm.OpArgK {
			fmt.Printf(" %d", -1-bx)
		} else if ins.BMode() == vm.OpArgU {
			fmt.Printf(" %d", bx)
		}
	case vm.IAsBx:
		a, sbx := ins.AsBx()
		fmt.Printf("%d %d", a, sbx)
	case vm.IAx:
		ax := ins.Ax()
		fmt.Printf("%d", -1-ax)
	}
}

func printDetail(f *binchunk.Prototype) {
	fmt.Printf("constants (%d):\n", len(f.Constants))
	for i, constant := range f.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(constant))
	}
	fmt.Printf("locals (%d):\n", len(f.LocVars))
	for i, local := range f.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, local.VarName, local.StartPC+1, local.EndPC+1)
	}
	fmt.Printf("upvalues (%d):\n", len(f.Upvalues))
	for i, upvalue := range f.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i,
			upvalueName(f, i), upvalue.Instack, upvalue.Idx)
	}
}

func upvalueName(f *binchunk.Prototype, i int) string {
	if len(f.UpvalueNames) > 0 {
		return f.UpvalueNames[i]
	}
	return "-"
}

func constantToString(k interface{}) interface{} {
	switch k.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", k)
	case float64:
		return fmt.Sprintf("%g", k)
	case int64:
		return fmt.Sprintf("%d", k)
	case string:
		return fmt.Sprintf("%q", k)
	default:
		return "?"
	}
}
