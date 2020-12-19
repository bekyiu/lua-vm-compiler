package binchunk

import (
	"encoding/binary"
	"fmt"
)

type reader struct {
	// 二进制chunk
	data []byte
}

// 读一个字节
func (this *reader) readByte() byte {
	b := this.data[0]
	this.data = this.data[1:]
	return b
}

// 读一个cint, data是小端格式
func (this *reader) readUint32() uint32 {
	u := binary.LittleEndian.Uint32(this.data)
	this.data = this.data[4:]
	return u
}

// 读一个size_t, data是小端格式
func (this *reader) readUint64() uint64 {
	u := binary.LittleEndian.Uint64(this.data)
	this.data = this.data[8:]
	return u
}

func (this *reader) readLuaInteger() int64 {
	return int64(this.readUint64())
}

func (this *reader) readLuaNumber() float64 {
	return float64(this.readUint64())
}

func (this *reader) readString() string {
	size := uint(this.readByte())
	// null
	if size == 0 {
		return ""
	}
	// 长字符串
	if size == 0xFF {
		size = uint(this.readUint64())
	}
	bytes := this.readBytes(size - 1)
	return string(bytes)
}

// 从字节流里读n个字节
func (this *reader) readBytes(n uint) []byte {
	bytes := this.data[:n]
	this.data = this.data[n:]
	return bytes
}

func (this *reader) checkHeader() {
	if string(this.readBytes(4)) != LUA_SIGNATURE {
		panic("二进制chunk签名错误")
	}
	if this.readByte() != LUAC_VERSION {
		panic("二进制chunk版本号不匹配")
	}
	if this.readByte() != LUAC_FORMAT {
		panic("二进制chunk format 不匹配")
	}
	if string(this.readBytes(6)) != LUAC_DATA {
		panic("二进制chunk luac data 不匹配")
	}
	if this.readByte() != CINT_SIZE {
		panic("cint大小不匹配")
	}
	if this.readByte() != CSIZET_SIZE {
		panic("size_t大小不匹配")
	}
	if this.readByte() != INSTRUCTION_SIZE {
		panic("lua指令宽度不匹配")
	}
	if this.readByte() != LUA_INTEGER_SIZE {
		panic("lua整数大小不匹配")
	}
	if this.readByte() != LUA_NUMBER_SIZE {
		panic("lua小数大小不匹配")
	}
	if this.readLuaInteger() != LUAC_INT {
		panic("需要二进制chunk为小端字节序")
	}
	if this.readLuaNumber() != LUAC_NUM {
		panic("需要二进制chunk浮点数遵守IEEE754规范")
	}
}

// 从字节流中读取函数原型
func (this *reader) readProto(parentSource string) *Prototype {
	// 以@开头, 来自.lua文件
	// 以=开头, 来自特定输入
	// 无, 从程序提供的字符串编译而来
	source := this.readString()
	// 子函数
	if source == "" {
		source = parentSource
	}

	return &Prototype{
		Source:          source,
		LineDefined:     this.readUint32(),
		LastLienDefined: this.readUint32(),
		NumParams:       this.readByte(),
		IsVararg:        this.readByte(),
		MaxStackSize:    this.readByte(),
		Codes:            this.readCodes(),
		Constants:       this.readConstants(),
		Upvalues:        this.readUpvalues(),
		Protos:          this.readProtos(source),
		LineInfos:        this.readLineInfos(),
		LocVars:         this.readLocVars(),
		UpvalueNames:    this.readUpvalueNames(),
	}
}

func (this *reader) readCodes() []uint32 {
	// 第一个4个字节是有多少条指令
	insNum := this.readUint32()
	codes := make([]uint32, insNum)
	for i, _ := range codes {
		codes[i] = this.readUint32()
	}
	return codes
}

func (this *reader) readConstants() []interface{} {
	constants := make([]interface{}, this.readUint32())
	for i, _ := range constants {
		constants[i] = this.readConstant()
	}
	return constants
}

func (this *reader) readConstant() interface{} {
	tag := this.readByte()
	switch tag {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return this.readByte() != 0
	case TAG_NUMBER:
		return this.readLuaNumber()
	case TAG_INTEGER:
		return this.readLuaInteger()
	case TAG_SHORT_STR:
	case TAG_LONG_STR:
		return this.readString()
	default:
		panic("未知的常量tag: " + fmt.Sprintf("0x%x", tag))
	}
	return nil
}

func (this *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, this.readUint32())
	for i, _ := range upvalues {
		upvalues[i] = Upvalue{
			Instack: this.readByte(),
			Idx: this.readByte(),
		}
	}
	return upvalues
}

func (this *reader) readProtos(source string) []*Prototype {
	protos := make([]*Prototype, this.readUint32())
	for i, _ := range protos {
		protos[i] = this.readProto(source)
	}
	return protos
}

func (this *reader) readLineInfos() []uint32 {
	// 有多少行
	lineCount := this.readUint32()
	lineInfos := make([]uint32, lineCount)
	for i, _ := range lineInfos {
		lineInfos[i] = this.readUint32()
	}
	return lineInfos
}

func (this *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, this.readUint32())
	for i, _ := range locVars {
		locVars[i] = LocVar{
			VarName: this.readString(),
			StartPC: this.readUint32(),
			EndPC: this.readUint32(),
		}
	}
	return locVars
}

func (this *reader) readUpvalueNames() []string {
	upvalueNames := make([]string, this.readUint32())
	for i, _ := range upvalueNames {
		upvalueNames[i] = this.readString()
	}
	return upvalueNames
}
