package lexer

import (
	"fmt"
	"regexp"
	"strings"
)
// 匹配注释左侧
var reOpeningLongBracket = regexp.MustCompile(`^\[=*\[`)
var reNewLine = regexp.MustCompile("\r\n|\n\r|\n|\r")

type Lexer struct {
	chunk     string // 源代码
	chunkName string // 原文件名
	line      int    // 当前行号
}

func NewLexer(chunk, chunkName string) *Lexer {
	return &Lexer{chunk, chunkName, 1}
}

// 判断源码是否以字符串s开头
func (this *Lexer) test(s string) bool {
	return strings.HasPrefix(this.chunk, s)
}

// 跳过n个字符
func (this *Lexer) next(n int) {
	this.chunk = this.chunk[n:]
}

// 跳过注释
// -- xxx, --> xxx
// --[[xxx]]
// --[===[xxx]===], '='的数量可以任意个
func (this *Lexer) skipComment() {
	// skip --
	this.next(2)
	if this.test("[") {
		if reOpeningLongBracket.FindString(this.chunk) != "" {
			this.scanLongString()
			return
		}
	}
	// short comment
	for len(this.chunk) > 0 && !isNewLine(this.chunk[0]) {
		this.next(1)
	}
}

// 跳过空白字符和注释, 增加行号
func (this *Lexer) skipWhiteSpaces() {
	for len(this.chunk) > 0 {
		if this.test("--") {
			this.skipComment()
		} else if this.test("\r\n") || this.test("\n\r") {
			this.next(2)
			this.line += 1
		} else if isNewLine(this.chunk[0]) {
			this.next(1)
			this.line += 1
		} else if isWhiteSpace(this.chunk[0]) {
			this.next(1)
		} else {
			break
		}
	}
}

// 提取左右长方括号中的内容
func (this *Lexer) scanLongString() string {
	// 左长方括号
	openingLongBracket := reOpeningLongBracket.FindString(this.chunk)
	if openingLongBracket == "" {
		this.error("invalid long string delimiter near '%s'", this.chunk[0:2])
	}
	// 对应的右长方括号
	closingLongBracket := strings.Replace(openingLongBracket, "[", "]", -1)
	closingLongBracketIdx := strings.Index(this.chunk, closingLongBracket)
	if closingLongBracketIdx < 0 {
		this.error("unfinished long string comment!")
	}
	// 提取注释文字
	str := this.chunk[len(openingLongBracket):closingLongBracketIdx]
	// 跳过注释
	this.next(closingLongBracketIdx + len(closingLongBracket))

	str = reNewLine.ReplaceAllString(str, "\n")
	this.line += strings.Count(str, "\n")
	if len(str) > 0 && str[0] == '\n' {
		str = str[1:]
	}
	return str
}

func (this *Lexer) error(f string, a ...interface{}) {
	err := fmt.Sprintf(f, a...)
	err = fmt.Sprintf("%s:%d: %s", this.chunkName, this.line, err)
	panic(err)
}

// 返回下一个token, 类型, 行号
//func (this *Lexer) NextToken() (line, kind int, token string) {
//
//}





func isWhiteSpace(c byte) bool {
	switch c {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}
	return false
}

func isNewLine(c byte) bool {
	return c == '\r' || c == '\n'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}