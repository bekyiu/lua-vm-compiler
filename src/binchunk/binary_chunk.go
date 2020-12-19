package binchunk

type binaryChunk struct {
	header                  // 头部
	sizeUpvalues byte       // 主函数upvalue数量
	mainFunc     *Prototype // 主函数原型
}

type header struct {
	// 校验用
	signature [4]byte // 魔数 0x1b4c7561
	version   byte    // 二进制chunk的版本号 5.3.4对应0x53
	format    byte    // 格式号 0
	luacData  [6]byte // 常量 0x19930d0a1a0a

	// 整数和lua虚拟机指令的宽度
	cintSize        byte // 4
	sizetSize       byte // 8
	instructionSize byte // 4
	luaIntegerSize  byte // 8
	luaNumberSize   byte // 8

	// 校验大小端和浮点数格式
	luacInt int64   // 0x5678
	luacNum float64 // 370.5
}

// header 中会用到的常量
const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

// 常量表 tag值
const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type Prototype struct {
	// 函数基本信息
	Source          string // 函数来源
	LineDefined     uint32 // 起始行号
	LastLienDefined uint32 // 终止行号
	NumParams       byte   // 固定参数个数
	IsVararg        byte   // 是否是可变参数
	MaxStackSize    byte   // 函数执行期间需要的寄存器个数

	Codes      []uint32      // 指令表
	Constants []interface{} // 常量表
	Upvalues  []Upvalue     // upvalue表
	Protos    []*Prototype  // 子函数原型

	// 调试信息
	LineInfos     []uint32 // 行号表
	LocVars      []LocVar // 局部变量表
	UpvalueNames []string // upvalue名列表
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

type Upvalue struct {
	Instack byte
	Idx     byte
}

// 解析二进制chunk
func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()
	reader.readByte()
	return reader.readProto("")
}