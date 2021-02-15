package ast

// chunk ::= block
// type Chunk *Block

// block ::= {stat} [retstat]
// retstat ::= return [explist] [‘;’]
// explist ::= exp {‘,’ exp}
type Block struct {
	LastLine int    // 代码块末尾的行号
	Stats    []Stat // 0个或多个语句
	RetExps  []Exp  // 0个人或一个return表达式
}
