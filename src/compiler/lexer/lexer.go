package lexer

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 匹配注释左侧
var reOpeningLongBracket = regexp.MustCompile(`^\[=*\[`)
var reNewLine = regexp.MustCompile("\r\n|\n\r|\n|\r")
var reShortStr = regexp.MustCompile(`(?s)(^'(\\\\|\\'|\\\n|\\z\s*|[^'\n])*')|(^"(\\\\|\\"|\\\n|\\z\s*|[^"\n])*")`)

var reDecEscapeSeq = regexp.MustCompile(`^\\[0-9]{1,3}`)
var reHexEscapeSeq = regexp.MustCompile(`^\\x[0-9a-fA-F]{2}`)
var reUnicodeEscapeSeq = regexp.MustCompile(`^\\u\{[0-9a-fA-F]+\}`)

var reNumber = regexp.MustCompile(`^0[xX][0-9a-fA-F]*(\.[0-9a-fA-F]*)?([pP][+\-]?[0-9]+)?|^[0-9]*(\.[0-9]*)?([eE][+\-]?[0-9]+)?`)
var reIdentifier = regexp.MustCompile(`^[_\d\w]+`)

type Lexer struct {
	chunk     string // 源代码
	chunkName string // 原文件名
	line      int    // 当前行号

	nextToken string
	nextTokenKind int
	nextTokenLine int
}

func NewLexer(chunk, chunkName string) *Lexer {
	return &Lexer{chunk, chunkName, 1, "", 0, 0}
}

// 看下一个token, 但是不跳过(从使用上来说不跳过, 但实现上相当于还是跳过了)
// 缓存在结构体里
func (this *Lexer) LookAhead() int {
	if this.nextTokenLine > 0 {
		return this.nextTokenKind
	}
	currentLine := this.line
	line, kind, token := this.NextToken()
	this.line = currentLine
	this.nextTokenLine = line
	this.nextTokenKind = kind
	this.nextToken = token
	return kind
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

	// 跳过所有换行
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

// 处理字符串中的转义字符
func (this *Lexer) escape(str string) string {
	var buf bytes.Buffer
	for len(str) > 0 {
		// 如果不是转移字符
		if str[0] != '\\' {
			buf.WriteByte(str[0])
			str = str[1:]
			continue
		}
		if len(str) == 1 {
			this.error("unfinished string")
		}
		switch str[1] {
		case 'a':
			buf.WriteByte('\a')
			str = str[2:]
			continue
		case 'b':
			buf.WriteByte('\b')
			str = str[2:]
			continue
		case 'f':
			buf.WriteByte('\f')
			str = str[2:]
			continue
		case 'n', '\n':
			buf.WriteByte('\n')
			str = str[2:]
			continue
		case 'r':
			buf.WriteByte('\r')
			str = str[2:]
			continue
		case 't':
			buf.WriteByte('\t')
			str = str[2:]
			continue
		case 'v':
			buf.WriteByte('\v')
			str = str[2:]
			continue
		case '"':
			buf.WriteByte('"')
			str = str[2:]
			continue
		case '\'':
			buf.WriteByte('\'')
			str = str[2:]
			continue
		case '\\':
			buf.WriteByte('\\')
			str = str[2:]
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // \ddd, 插入任意ascii
			if found := reDecEscapeSeq.FindString(str); found != "" {
				d, _ := strconv.ParseInt(found[1:], 10, 32)
				if d < 0xFF {
					buf.WriteByte(byte(d))
					str = str[len(found):]
					continue
				}
				this.error("decimal escape too large near '%s'", found)
			}
		case 'x': // \xXX, XX是16进制, 也是插入任意ascii
			if found := reHexEscapeSeq.FindString(str); found != "" {
				d, _ := strconv.ParseInt(found[2:], 16, 32)
				buf.WriteByte(byte(d))
				str = str[len(found):]
				continue
			}
		case 'u': // \u{XXX...} X是16进制 用于插入unicode
			if found := reUnicodeEscapeSeq.FindString(str); found != "" {
				d, err := strconv.ParseInt(found[3:len(found) - 1], 16, 32)
				if err == nil && d <= 0x10FFFF {
					buf.WriteRune(rune(d))
					str = str[len(found):]
					continue
				}
				this.error("utf-8 value too large near '%s'", found)
			}
		case 'z': // 用于跳过紧跟其后的空白字符
			str = str[2:]
			for len(str) > 0 && isWhiteSpace(str[0]) {
				str = str[1:]
			}
			continue
		}
	}
	return buf.String()
}

// 提取短字符串
func (this *Lexer) scanShortString() string {
	if str := reShortStr.FindString(this.chunk); str != "" {
		this.next(len(str))
		// 跳过引号
		str = str[1 : len(str)-1]
		// 处理转义字符
		if strings.Index(str, `\`) >= 0 {
			this.line += len(reNewLine.FindAllString(str, -1))
			str = this.escape(str)
		}
		return str
	}
	this.error("unfinished string")
	return ""
}

// 根据正则表达式提取字符串
func (this *Lexer) scan(re *regexp.Regexp) string {
	if token := re.FindString(this.chunk); token != "" {
		this.next(len(token))
		return token
	}
	panic("unreachable!")
}

// 提取数字
func (this *Lexer) scanNumber() string {
	return this.scan(reNumber)
}

// 提取标识符
func (this *Lexer) scanIdentifier() string {
	return this.scan(reIdentifier)
}
// 返回下一个token, 类型, 行号
func (this *Lexer) NextToken() (line, kind int, token string) {
	// 从缓存里那
	if this.nextTokenLine > 0 {
		line = this.nextTokenLine
		kind = this.nextTokenKind
		token = this.nextToken
		this.line = this.nextTokenLine
		this.nextTokenLine = 0
		return
	}

	this.skipWhiteSpaces()
	if len(this.chunk) == 0 {
		return this.line, TOKEN_EOF, "EOF"
	}

	switch this.chunk[0] {
	case ';':
		this.next(1)
		return this.line, TOKEN_SEP_SEMI, ";"
	case ',':
		this.next(1)
		return this.line, TOKEN_SEP_COMMA, ","
	case '(':
		this.next(1)
		return this.line, TOKEN_SEP_LPAREN, "("
	case ')':
		this.next(1)
		return this.line, TOKEN_SEP_RPAREN, ")"
	case ']':
		this.next(1)
		return this.line, TOKEN_SEP_RBRACK, "]"
	case '{':
		this.next(1)
		return this.line, TOKEN_SEP_LCURLY, "{"
	case '}':
		this.next(1)
		return this.line, TOKEN_SEP_RCURLY, "}"
	case '+':
		this.next(1)
		return this.line, TOKEN_OP_ADD, "+"
	case '-':
		this.next(1)
		return this.line, TOKEN_OP_MINUS, "-"
	case '*':
		this.next(1)
		return this.line, TOKEN_OP_MUL, "*"
	case '^':
		this.next(1)
		return this.line, TOKEN_OP_POW, "^"
	case '%':
		this.next(1)
		return this.line, TOKEN_OP_MOD, "%"
	case '&':
		this.next(1)
		return this.line, TOKEN_OP_BAND, "&"
	case '|':
		this.next(1)
		return this.line, TOKEN_OP_BOR, "|"
	case '#':
		this.next(1)
		return this.line, TOKEN_OP_LEN, "#"
	case ':':
		if this.test("::") {
			this.next(2)
			return this.line, TOKEN_SEP_LABEL, "::"
		} else {
			this.next(1)
			return this.line, TOKEN_SEP_COLON, ":"
		}
	case '/':
		if this.test("//") {
			this.next(2)
			return this.line, TOKEN_OP_IDIV, "//"
		} else {
			this.next(1)
			return this.line, TOKEN_OP_DIV, "/"
		}
	case '~':
		if this.test("~=") {
			this.next(2)
			return this.line, TOKEN_OP_NE, "~="
		} else {
			this.next(1)
			return this.line, TOKEN_OP_WAVE, "~"
		}
	case '=':
		if this.test("==") {
			this.next(2)
			return this.line, TOKEN_OP_EQ, "=="
		} else {
			this.next(1)
			return this.line, TOKEN_OP_ASSIGN, "="
		}
	case '<':
		if this.test("<<") {
			this.next(2)
			return this.line, TOKEN_OP_SHL, "<<"
		} else if this.test("<=") {
			this.next(2)
			return this.line, TOKEN_OP_LE, "<="
		} else {
			this.next(1)
			return this.line, TOKEN_OP_LT, "<"
		}
	case '>':
		if this.test(">>") {
			this.next(2)
			return this.line, TOKEN_OP_SHR, ">>"
		} else if this.test(">=") {
			this.next(2)
			return this.line, TOKEN_OP_GE, ">="
		} else {
			this.next(1)
			return this.line, TOKEN_OP_GT, ">"
		}
	case '.':
		if this.test("...") {
			this.next(3)
			return this.line, TOKEN_VARARG, "..."
		} else if this.test("..") {
			this.next(2)
			return this.line, TOKEN_OP_CONCAT, ".."
		} else if len(this.chunk) == 1 || !isDigit(this.chunk[1]) {
			this.next(1)
			return this.line, TOKEN_SEP_DOT, "."
		}
	case '[':
		if this.test("[[") || this.test("[=") {
			return this.line, TOKEN_STRING, this.scanLongString()
		} else {
			this.next(1)
			return this.line, TOKEN_SEP_LBRACK, "["
		}
	case '\'', '"':
		return this.line, TOKEN_STRING, this.scanShortString()
	}

	c := this.chunk[0]
	if c == '.' || isDigit(c) {
		token := this.scanNumber()
		return this.line, TOKEN_NUMBER, token
	}
	if c == '_' || isLetter(c) {
		token := this.scanIdentifier()
		if kind, found := keywords[token]; found {
			return this.line, kind, token // keyword
		} else {
			return this.line, TOKEN_IDENTIFIER, token
		}
	}

	this.error("unexpected symbol near %q", c)
	return
}

// 提取指定类型的token
func (this *Lexer) NextTokenOfKind(kind int) (line int, token string) {
	line, _kind, token := this.NextToken()
	if kind != _kind {
		this.error("syntax error near '%s'", token)
	}
	return line, token
}

// 提取标识符
func (this *Lexer) NextIdentifier() (line int, token string) {
	return this.NextTokenOfKind(TOKEN_IDENTIFIER)
}
// 返回当前行号
func (this *Lexer) Line() int {
	return this.line
}

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
