package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TokenKind int

const (
	TK_RESERVED TokenKind = iota
	TK_IDENT
	TK_NUM
	TK_STRING
	TK_EOF
)

type Token struct {
	next     *Token
	kind     TokenKind
	val      int
	str      string
	len      int
	contents string
}

func errorAt(loc string, format string, a ...string) {
	fmt.Println(loc)
	fmt.Printf(format, a)
	fmt.Println()
	os.Exit(1)
}

func errorToken(tok *Token, format string, a ...string) {
	fmt.Printf(format, a)
	fmt.Println()
	os.Exit(1)
}

func NewToken(kind TokenKind, cur *Token, str string, len int) *Token {
	tok := &Token{
		kind: kind,
		str:  str,
		len:  len,
	}
	cur.next = tok
	return tok
}

func isWhiteSpace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

func isLetter(r rune) bool {
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z') || r == '_'
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isAlNum(r rune) bool {
	return isDigit(r) || isLetter(r)
}

func isPunct(r rune) bool {
	return strings.ContainsRune("+-*/=(){}[]<>!;:,.&", r)
}

func (p *Parser) consume(op string) bool {
	if p.token.kind != TK_RESERVED || p.token.str != op {
		return false
	}
	p.token = p.token.next
	return true
}

func (p *Parser) peek(s string) bool {
	return p.token.str == s
}

func (p *Parser) consumeIdent() *Token {
	if p.token.kind != TK_IDENT {
		return nil
	}
	token := p.token
	p.token = token.next
	return token
}

func (p *Parser) expect(op string) {
	if p.token.kind != TK_RESERVED || p.token.str != op {
		errorAt(p.token.str, "expected '%s'", op)
	}
	p.token = p.token.next
}

func (p *Parser) expectNumber() int {
	if p.token.kind != TK_NUM {
		errorAt(p.token.str, "expect a number")
	}
	val := p.token.val
	p.token = p.token.next
	return val
}

func (p *Parser) expectIdent() string {
	if p.token.kind != TK_IDENT {
		errorAt(p.token.str, "expected an identifier")
	}
	s := p.token.str
	p.token = p.token.next
	return s
}

func (t *Token) AtEOF() bool {
	return t.kind == TK_EOF
}

func startsWithReserved(s string) (string, bool) {
	keywords := []string{
		"return", "if", "else", "while", "for",
		"int", "char", "sizeof", "struct",
	}
	for _, v := range keywords {
		if strings.HasPrefix(s, v) {
			return v, true
		}
	}

	ops := []string{"==", "!=", "<=", ">="}
	for _, v := range ops {
		if strings.HasPrefix(s, v) {
			return v, true
		}
	}

	return "", false
}

func getEscapeChar(r rune) rune {
	switch r {
	case 'a':
		return '\a'
	case 'b':
		return '\b'
	case 't':
		return '\t'
	case 'n':
		return '\n'
	case 'v':
		return '\v'
	case 'f':
		return '\f'
	case 'r':
		return '\r'
	case 'e':
		return 27
	case '0':
		return 0
	default:
		return r
	}
}

func readStringLiteral(cur *Token, start string) *Token {
	str := ""
	i := 1
	for {
		if i == len(start) {
			errorAt("", "unclosed string lieteral")
		}

		if start[i] == '"' {
			break
		}

		if start[i] == '\\' {
			i++
			str += string(getEscapeChar(rune(start[i])))
			i++
		} else {
			str += string(start[i])
			i++
		}
	}

	tok := NewToken(TK_STRING, cur, start[:i], i+1)
	tok.contents = str
	return tok
}

func Tokenize(input string) *Token {
	head := &Token{}
	cur := head
	i := 0
	for i = 0; i < len(input); {
		if isWhiteSpace(rune(input[i])) {
			i++
			continue
		}
		if strings.HasPrefix(input[i:], "//") {
			i += 2
			for input[i] != '\n' {
				i++
			}
			continue
		}
		if strings.HasPrefix(input[i:], "/*") {
			i += 2
			index := strings.Index(input[i:], "*/")
			if index == -1 {
				errorAt("", "unclosed block comment")
			}
			i += index + 2
			continue
		}
		if keyword, ok := startsWithReserved(input[i:]); ok {
			cur = NewToken(TK_RESERVED, cur, keyword, len(keyword))
			i += len(keyword)
			continue
		}
		if isLetter(rune(input[i])) {
			pos := i
			i++
			for ; i < len(input) && isAlNum(rune(input[i])); i++ {
			}
			cur = NewToken(TK_IDENT, cur, input[pos:i], i-pos)
			continue
		}
		if isPunct(rune(input[i])) {
			cur = NewToken(TK_RESERVED, cur, input[i:i+1], 1)
			i++
			continue
		}
		if rune(input[i]) == '"' {
			cur = readStringLiteral(cur, input[i:])
			i += cur.len
			continue
		}
		if isDigit(rune(input[i])) {
			pos := i
			for ; i < len(input) && isDigit(rune(input[i])); i++ {
			}
			cur = NewToken(TK_NUM, cur, input[pos:i], i-pos)
			val, _ := strconv.ParseInt(input[pos:i], 10, 32)
			cur.val = int(val)
			continue
		}

		errorAt(input[i:], "invalid token")
	}

	NewToken(TK_EOF, cur, input[i:], 0)
	return head.next
}
