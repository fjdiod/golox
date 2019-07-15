package scanner


type Token struct {
	Type TokenType
	Lexeme string
	Literal interface{}
}
