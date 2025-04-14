package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"unicode"
)

type Parser struct {
	data []byte
	pos  int
}

func NewParser(input []byte) *Parser {
	return &Parser{data: input, pos: 0}
}

func (p *Parser) peek() byte {
	return p.data[p.pos]
}

func (p *Parser) next() byte {
	ch := p.data[p.pos]
	p.pos++
	return ch
}

func (p *Parser) parse() interface{} {
	ch := p.peek()
	switch ch {
	case 'i':
		return p.parseInt()
	case 'l':
		return p.parseList()
	case 'd':
		return p.parseDict()
	default:
		if unicode.IsDigit(rune(ch)) {
			return p.parseString()
		}
		panic(fmt.Sprintf("unexpected character: %c", ch))
	}
}

func (p *Parser) parseInt() int {
	p.next()
	start := p.pos

	for p.data[p.pos] != 'e' {
		p.pos++
	}
	valStr := string(p.data[start:p.pos])

	length, err := strconv.Atoi(valStr)
	if err != nil {
		panic("invalid integer: " + valStr)
	}
	p.pos++
	return length
}

func (p *Parser) parseString() string {
	start := p.pos
	for p.data[p.pos] != ':' {
		p.pos++
	}
	lenStr := string(p.data[start:p.pos])
	length, err := strconv.Atoi(lenStr)
	if err != nil {
		panic("invalid string length")
	}
	p.pos++ // skip ':'
	// read string content
	end := p.pos + length
	str := string(p.data[p.pos:end])
	p.pos = end
	return str

}

func (p *Parser) parseList() []interface{} {
	// like if I talk about the implementation of the list then it is straight forward like we just have to
	// check that and parse that the other work is done by the parse func
	p.next()
	var result []interface{}
	for p.data[p.pos] != 'e' {
		result = append(result, p.parse())
	}
	p.pos++
	return result
}

func (p *Parser) parseDict() map[string]interface{} {
	p.next()

	result := make(map[string]interface{})
	// if i create my own logic then there is some type of mechansim through which we have to distinguish between the key and values
	for p.peek() != 'e' {
		key := p.parse().(string)
		value := p.parse()
		result[key] = value
	}
	p.pos++
	return result
}

func GeneratePeerID() string {
	return "-GT0001-" + RandString(12)
}

func RandString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func EncodeInfoHash(infoHash [20]byte) string {
	var encoded string
	for _, b := range infoHash {
		encoded += fmt.Sprintf("%%%02X", b) // %XX format
	}
	return encoded
}
