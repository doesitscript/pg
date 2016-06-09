package parser

import (
	"bytes"
	"strconv"
)

type Parser struct {
	b []byte
}

func New(b []byte) *Parser {
	return &Parser{b: b}
}

func NewString(s string) *Parser {
	return &Parser{b: []byte(s)}
}

func (p *Parser) Bytes() []byte {
	return p.b
}

func (p *Parser) Valid() bool {
	return len(p.b) > 0
}

func (p *Parser) Read() byte {
	c := p.b[0]
	p.Skip(c)
	return c
}

func (p *Parser) Peek() byte {
	if p.Valid() {
		return p.b[0]
	}
	return 0
}

func (p *Parser) Advance() {
	p.b = p.b[1:]
}

func (p *Parser) Skip(c byte) bool {
	if p.Peek() == c {
		p.Advance()
		return true
	}
	return false
}

func (p *Parser) SkipString(s string) bool {
	if len(s) > len(p.b) {
		return false
	}
	if !bytes.Equal(p.b[:len(s)], []byte(s)) {
		return false
	}
	p.b = p.b[len(s):]
	return true
}

func (p *Parser) ReadSep(c byte) ([]byte, bool) {
	ind := bytes.IndexByte(p.b, c)
	if ind == -1 {
		b := p.b
		p.b = p.b[len(p.b):]
		return b, false
	}

	b := p.b[:ind]
	p.b = p.b[ind+1:]
	return b, true
}

func (p *Parser) ReadIdentifier() []byte {
	end := len(p.b)
	for i, ch := range p.b {
		if !(isAlnum(ch) || ch == '_') {
			end = i
			break
		}
	}
	if end <= 0 {
		return nil
	}
	b := p.b[:end]
	p.b = p.b[end:]
	return b
}

func (p *Parser) ReadNumber() int {
	end := len(p.b)
	for i, ch := range p.b {
		if !isNum(ch) {
			end = i
			break
		}
	}
	if end <= 0 {
		return 0
	}
	n, _ := strconv.Atoi(string(p.b[:end]))
	p.b = p.b[end:]
	return n
}

func (p *Parser) readSubstring() []byte {
	var b []byte
	for p.Valid() {
		c := p.Read()
		switch c {
		case '\\':
			switch p.Peek() {
			case '\\':
				b = append(b, '\\')
				p.Advance()
			case '"':
				b = append(b, '"')
				p.Advance()
			default:
				b = append(b, c)
			}
		case '\'':
			switch p.Peek() {
			case '\'':
				b = append(b, '\'')
				p.Skip(c)
			default:
				b = append(b, c)
			}
		case '"':
			return b
		default:
			b = append(b, c)
		}
	}
	return b
}
