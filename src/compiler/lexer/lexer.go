package lexer

type Lexer struct {
	chunk     string // 源代码
	chunkName string // 原文件名
	line      int    // 当前行号
}

func NewLexer(chunk, chunkName string) *Lexer {
	return &Lexer{chunk, chunkName, 1}
}

// 返回下一个token, 类型, 行号
//func (this *Lexer) NextToken() (line, kind int, token string) {
//
//}