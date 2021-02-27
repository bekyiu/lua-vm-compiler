package parser

import . "write_lua/src/compiler/ast"
import . "write_lua/src/compiler/lexer"

// parser 入口
func Parse(chunk, chunkName string) *Block {
	lexer := NewLexer(chunk, chunkName)
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_EOF)
	return block
}
