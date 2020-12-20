package vm

type Instruction uint32

const MAXARG_Bx = (1 << 18) - 1   // 2^18 - 1
const MAXARG_sBx = MAXARG_Bx >> 1 // (2^18 - 1) / 2

// 从指令中提取操作码
func (this Instruction) Opcode() int {
	//0b 0000 0000 0000 0000 0000 0011 1111
	return int(this & 0x3F)
}

// 从iABC模式的指令中提取参数
func (this Instruction) ABC() (a, b, c int) {
	a = int(this >> 6 & 0xFF)
	b = int(this >> 14 & 0x1FF)
	c = int(this >> 23 & 0x1FF)
	return
}

// 从iABx模式的命令中提出参数
func (this Instruction) ABx() (a, bx int) {
	a = int(this >> 6 & 0xFF)
	bx = int(this >> 14 & 0x3FFFF)
	return
}

// iAsBx
func (this Instruction) AsBx() (a, sbx int) {
	a, bx := this.ABx()
	// bx是有符号数, 减去一个偏移
	return a, bx - MAXARG_sBx
}

func (this Instruction) Ax() int {
	return int(this >> 6)
}

func (this Instruction) OpName() string {
	return opcodes[this.Opcode()].name
}
func (this Instruction) OpMode() byte {
	return opcodes[this.Opcode()].opMode
}
func (this Instruction) BMode() byte {
	return opcodes[this.Opcode()].argBMode
}
func (this Instruction) CMode() byte {
	return opcodes[this.Opcode()].argCMode
}