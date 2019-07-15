package scanner

import (
	"strconv"
)

type Scanner struct {
	Source string
	Current int
	Start int
	Tokens []Token
	Line int
	KeyWords map[string]TokenType
}

func NewScanner(source string) *Scanner {
	scanner := Scanner{source, 0, 0, []Token{}, 0, nil}
	scanner.KeyWords = map[string]TokenType{
		"if": IF, "else": ELSE, "elif": ELIF, "while": WHILE,  "break": BREAK,
		"continue": CONTINUE, "for": FOR, "fun": FUN, "nil": NIL, "var": VAR,
		"true": TRUE, "FALSE": FALSE, "print": PRINT, "return": RETURN,
		"or": OR, "and": AND,
	}
	return &scanner
}

func (s *Scanner) scan() {
	for !s.isAtEnd() {
		s.Start = s.Current
		c := s.advance()
		s.addToken(c)
	}
	s.put(Token{EOF, "", ""})
}

func (s * Scanner) addToken(c string) {
	switch c{
	case "(":
		s.put(Token{LP, "(", nil,})
	case ")":
		s.put(Token{RP, ")", nil,})
	case "{":
		s.put(Token{LB, "{", nil,})
	case "}":
		s.put(Token{RB, "}", nil,})
	case ".":
		s.put(Token{DOT, ".", nil,})
	case ",":
		s.put(Token{COMMA, ",", nil,})
	case ";":
		s.put(Token{SEMICOLON, ";", nil,})
	case ":":
		s.put(Token{COLON, ":", nil,})
	case "+":
		if s.match("=") {
			s.put(Token{PLUS_EQUAL, "+=", nil})
		} else {
			s.put(Token{PLUS, "+", nil})
		}
	case "-":
		if s.match("=") {
			s.put(Token{MINUS_EQUAL, "-=", nil})
		} else {
			s.put(Token{MINUS, "-", nil})
		}
	case "*":
		if s.match("=") {
			s.put(Token{STAR_EQUAL, "*=", nil})
		} else {
			s.put(Token{STAR, "*", nil})
		}
	case "/":
		if s.match("=") {
			s.put(Token{SLASH_EQUAL, "/=", nil})
		} else {
			s.put(Token{SLASH, "/", nil})
		}
	case "?":
		s.put(Token{STAR, "?", nil,})
	case "!":
		if s.match("=") {
			s.put(Token{BANG_EQUAL, "!=", nil})
		} else {
			s.put(Token{BANG, "!", nil})
		}
	case ">":
		if s.match("=") {
			s.put(Token{GREATER_EQUAL, ">=", nil})
		} else {
			s.put(Token{GRATER, ">", nil})
		}
	case "<":
		if s.match("=") {
			s.put(Token{LESS_EQUAL, "<=", nil})
		} else {
			s.put(Token{LESS, "<", nil})
		}
	case "=":
		if s.match("=") {
			s.put(Token{EQUAL_EQUAL, "==", nil})
		} else {
			s.put(Token{EQUAL, "=", nil})
		}
	case " ":
		break
	case "\t":
		break
	case "\r":
		break
	case "\n":
		s.Line++
	case `"`:
		s.string()
	default:
		if isNumeric(c) {
			s.number();
		}
		if isAlpha(c) {
			s.identifier()
		}
	}
}

func (s * Scanner) advance() string {
	s.Current++
	return string(s.Source[s.Current-1])
}

func (s * Scanner) match(expected string) bool {
	if s.isAtEnd() {
		return false
	}
	if string(s.Source[s.Current]) != expected {
		return false
	}
	s.Current++
	return true
}

func (s * Scanner) put(token Token) {
	s.Tokens = append(s.Tokens, token)
}

func (s * Scanner) string() {
	for !s.isAtEnd() && s.peek() != `"` {
		if (s.peek() == "\n") {
			s.Line++
		}
		s.advance()
	}
	s.advance()
	lexeme := s.Source[s.Start:s.Current]
	value := s.Source[s.Start+1:s.Current-1]
	s.put(Token{STRING, lexeme, value})
}

func (s * Scanner) number() {
	for isNumeric(s.peek()) {
		s.advance()
	}
	if s.peek() == "." && isNumeric(s.peekNext()) {
		s.advance()
		for isNumeric(s.peek()) {
			s.advance()
		}
	}
	lexeme := s.Source[s.Start:s.Current]
	f, _ := strconv.ParseFloat(lexeme, 64)
	s.put(Token{NUMBER, lexeme, f})
}

func (s * Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	lexeme := s.Source[s.Start:s.Current]
	if tp, ok := s.KeyWords[lexeme]; ok {
		s.put(Token{tp, lexeme, nil})
		return
	}
	s.put(Token{STRING, lexeme, nil})
}

func (s * Scanner) peek() string {
	if s.isAtEnd() {
		return "\x00"
	}
	return string(s.Source[s.Current])
}

func (s * Scanner) peekNext() string {
	if s.Current + 1 >= len(s.Source) {
		return "\x00"
	}
	return string(s.Source[s.Current + 1])
}

func (s * Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func isAlpha(c string) bool {
	return (c >= "a" && c <= "z") || (c >= "A" && c <= "Z") || c == "_"
}

func isNumeric(c string) bool {
	return c >= "0" && c <= "9"
}

func isAlphaNumeric(c string) bool {
	return isAlpha(c) || isAlpha(c)
}