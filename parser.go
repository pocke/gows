package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

const (
	SPACE = ' '
	TAB   = '\t'
	NL    = '\n'

	EOF = 0
)

const NO_OPERAND = 0

type Parser struct {
	sc  *bufio.Scanner
	eof bool
	cs  []Command
}

func NewParser(r io.Reader) *Parser {
	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanRunes)
	p := &Parser{
		sc:  sc,
		eof: false,
		cs:  make([]Command, 0),
	}
	return p
}

func (p *Parser) Parse() ([]Command, error) {
L:
	for !p.eof {
		c := p.nextc()
		switch c {
		case SPACE:
			p.parseStack()
		case NL:
			p.parseFlow()
		case TAB:
			c := p.nextc()
			switch c {
			case SPACE:
				p.parseCalc()
			case TAB:
				p.parseHeap()
			case NL:
				p.parseIO()
			default:
				return nil, errors.Errorf("%s is unknown character", c)
			}
		case EOF:
			break L
		default:
			return nil, errors.Errorf("%s is unknown character", c)
		}
	}
	return p.cs, nil
}

func (p *Parser) parseStack() error {
	c := p.nextc()
	switch c {
	case SPACE:
		p.cs = append(p.cs, Command{C_STACK_PUSH, p.nextint()})
	case NL:
		c := p.nextc()
		switch c {
		case SPACE:
			p.cs = append(p.cs, Command{C_STACK_DUP, NO_OPERAND})
		case TAB:
			p.cs = append(p.cs, Command{C_STACK_SWAP, NO_OPERAND})
		case NL:
			p.cs = append(p.cs, Command{C_STACK_POP, NO_OPERAND})
		default:
			return errors.Errorf("%s is unknown character", c)
		}
	default:
		return errors.Errorf("%s is unknown character", c)
	}

	return nil
}

func (p *Parser) parseFlow() error {
	c := p.nextc()
	switch c {
	case SPACE:
		c := p.nextc()
		switch c {
		case SPACE:
			p.cs = append(p.cs, Command{C_FLOW_DEF, p.nextlabel()})
		case TAB:
			p.cs = append(p.cs, Command{C_FLOW_CALL, p.nextlabel()})
		case NL:
			p.cs = append(p.cs, Command{C_FLOW_JUMP, p.nextlabel()})
		default:
			return errors.Errorf("%s is unknown character", c)
		}
	case TAB:
		c := p.nextc()
		switch c {
		case SPACE:
			p.cs = append(p.cs, Command{C_FLOW_JUMP_IF_ZERO, p.nextlabel()})
		case TAB:
			p.cs = append(p.cs, Command{C_FLOW_JUMP_IF_NEG, p.nextlabel()})
		case NL:
			p.cs = append(p.cs, Command{C_FLOW_END, NO_OPERAND})
		default:
			return errors.Errorf("%s is unknown character", c)
		}
	case NL:
		c := p.nextc()
		switch c {
		case NL:
			p.cs = append(p.cs, Command{C_FLOW_EXIT, NO_OPERAND})
		default:
			return errors.Errorf("%s is unknown character", c)
		}
	default:
		return errors.Errorf("%s is unknown character", c)
	}

	return nil
}

func (p *Parser) parseCalc() error {
	c := p.nextc()
	switch c {
	case SPACE:
		c := p.nextc()
		switch c {
		case SPACE:
			p.cs = append(p.cs, Command{C_CALC_ADD, NO_OPERAND})
		case TAB:
			p.cs = append(p.cs, Command{C_CALC_SUB, NO_OPERAND})
		case NL:
			p.cs = append(p.cs, Command{C_CALC_MULTI, NO_OPERAND})
		default:
			return errors.Errorf("%s is unknown character", c)
		}
	case TAB:
		c := p.nextc()
		switch c {
		case SPACE:
			p.cs = append(p.cs, Command{C_CALC_DIV, NO_OPERAND})
		case TAB:
			p.cs = append(p.cs, Command{C_CALC_MOD, NO_OPERAND})
		default:
			return errors.Errorf("%s is unknown character", c)
		}
	default:
		return errors.Errorf("%s is unknown character", c)
	}

	return nil
}

func (p *Parser) parseHeap() error {
	c := p.nextc()
	switch c {
	case SPACE:
		p.cs = append(p.cs, Command{C_HEAP_SAVE, NO_OPERAND})
	case TAB:
		p.cs = append(p.cs, Command{C_HEAP_LOAD, NO_OPERAND})
	default:
		return errors.Errorf("%s is unknown character", c)
	}

	return nil
}

func (p *Parser) parseIO() error {
	c := p.nextc()
	switch c {
	case SPACE:
		c := p.nextc()
		switch c {
		case SPACE:
			p.cs = append(p.cs, Command{C_IO_WRITE_CH, NO_OPERAND})
		case TAB:
			p.cs = append(p.cs, Command{C_IO_WRITE_NUM, NO_OPERAND})
		default:
			return errors.Errorf("%s is unknown character", c)
		}
	case TAB:
		c := p.nextc()
		switch c {
		case SPACE:
			p.cs = append(p.cs, Command{C_IO_READ_CH, NO_OPERAND})
		case TAB:
			p.cs = append(p.cs, Command{C_IO_READ_NUM, NO_OPERAND})
		default:
			return errors.Errorf("%s is unknown character", c)
		}
	default:
		return errors.Errorf("%s is unknown character", c)
	}

	return nil
}

// returns SPACE, TAB or NL
func (p *Parser) nextc() byte {
	for p.sc.Scan() {
		bs := p.sc.Bytes()
		if len(bs) != 1 {
			continue
		}
		b := bs[0]
		if b != SPACE && b != TAB && b != NL {
			continue
		}

		return b
	}

	p.eof = true
	return EOF
}

func (p *Parser) nextint() int {
	sign := p.nextc()
	n := 0
L:
	for {
		c := p.nextc()
		switch c {
		case SPACE:
			// plus 0
		case TAB:
			n += 1
		case NL:
			break L
		default:
			panic(fmt.Sprintf("Unknown: %+v", c))
		}
		n *= 2
	}
	n /= 2
	if sign == TAB {
		n = -n
	}
	return n
}

func (p *Parser) nextlabel() int {
	label := 0
L:
	for {
		c := p.nextc()
		switch c {
		case TAB:
			label += 1
		case SPACE:
			// plus 0
		case NL:
			break L
		default:
			panic(fmt.Sprintf("Unknown: %+v", c))
		}
		label = label * 2
	}
	return label
}
