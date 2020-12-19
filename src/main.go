package main

import (
	"fmt"
	"io/ioutil"
	"write_lua/src/binchunk"
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
		// 序号, 行号, 16进制的指令
		fmt.Printf("\t%d\t[%s]\t0x%08X\n", pc+1, line, code)
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
