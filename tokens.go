package compiler

type TokenKind uint16

const (
	TokenEOF TokenKind = iota
	TokenSpace
	TokenNewLine      // \n
	TokenQuotation    // "
	TokenParenLeft    // (
	TokenParenRight   // )
	TokenBraceLeft    // {
	TokenBraceRight   // }
	TokenBracketLeft  // [
	TokenBracketRight // ]
	TokenComma        // ,
	TokenColon        // :

	TokenPlus    // +
	TokenMinus   // -
	TokenStar    // *
	TokenSlash   // /
	TokenModulus // %

	TokenCompare  // ==
	TokenNotEqual // !=
	TokenAnd      // &&
	TokenOr       // ||
	TokenNot      // !

	TokenBitAnd     // &
	TokenBitOr      // |
	TokenBitXor     // ^
	TokenBitClear   // &^
	TokenShiftLeft  // <<
	TokenShiftRight // >>

	TokenEqual            // =
	TokenAddAssign        // +=
	TokenSubAssign        // -=
	TokenMulAssign        // *=
	TokenDivAssign        // /=
	TokenModAssign        // %=
	TokenBitAndAssign     // &=
	TokenBitOrAssign      // |=
	TokenBitXorAssign     // ^=
	TokenBitClearAssign   // &^=
	TokenShiftLeftAssign  // <<=
	TokenShiftRightAssign // >>=

	TokenIncrement    // ++
	TokenDecrement    // --
	TokenGreater      // >
	TokenGreaterEqual // >=
	TokenLess         // <
	TokenLessEqual    // <=

	TokenArrowLeft // <-

	TokenTab
	TokenPackage
	TokenImport
	TokenConst
	TokenVar
	TokenIf
	TokenElse
	TokenFunc
	TokenDefer
	TokenMake
	TokenGo
	TokenMap
	TokenType
	TokenChan
	TokenStruct
	TokenInterface
	TokenSwitch
	TokenCase
	TokenDefault
	TokenReturn
)

var tokens = [][]byte{
	[]byte(""),
	[]byte(" "),
	[]byte("\n"),
	[]byte("\""),
	[]byte("("),
	[]byte(")"),
	[]byte("{"),
	[]byte("}"),
	[]byte("["),
	[]byte("]"),
	[]byte(","),
	[]byte(":"),

	[]byte("+"),
	[]byte("-"),
	[]byte("*"),
	[]byte("/"),
	[]byte("%"),

	[]byte("=="),
	[]byte("!="),
	[]byte("&&"),
	[]byte("||"),
	[]byte("!"),

	[]byte("&"),
	[]byte("|"),
	[]byte("^"),
	[]byte("&^"),
	[]byte("<<"),
	[]byte(">>"),

	[]byte("="),
	[]byte("+="),
	[]byte("-="),
	[]byte("*="),
	[]byte("/="),
	[]byte("%="),
	[]byte("&="),
	[]byte("|="),
	[]byte("^="),
	[]byte("&^="),
	[]byte("<<="),
	[]byte(">>="),

	[]byte("++"),
	[]byte("--"),
	[]byte(">"),
	[]byte(">="),
	[]byte("<"),
	[]byte("<="),

	[]byte("<-"),

	[]byte("    "),
	[]byte("package"),
	[]byte("import"),
	[]byte("const"),
	[]byte("var"),
	[]byte("if"),
	[]byte("else"),
	[]byte("func"),
	[]byte("defer"),
	[]byte("make"),
	[]byte("go"),
	[]byte("map"),
	[]byte("type"),
	[]byte("chan"),
	[]byte("struct"),
	[]byte("interface"),
	[]byte("switch"),
	[]byte("case"),
	[]byte("default"),
	[]byte("return"),
}

func (t TokenKind) Bytes() []byte {
	if int(t) >= len(tokens) {
		return nil
	}
	return tokens[t]
}

func (t TokenKind) Len() int {
	if int(t) >= len(tokens) {
		return 0
	}
	return len(tokens[t])
}
