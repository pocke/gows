package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type CommandType int

const (
	C_STACK_PUSH CommandType = iota + 1
	C_STACK_DUP
	C_STACK_SWAP
	C_STACK_POP

	C_CALC_ADD
	C_CALC_SUB
	C_CALC_MULTI
	C_CALC_DIV
	C_CALC_MOD

	C_HEAP_SAVE
	C_HEAP_LOAD

	C_FLOW_DEF
	C_FLOW_CALL
	C_FLOW_JUMP
	C_FLOW_JUMP_IF_ZERO
	C_FLOW_JUMP_IF_NEG
	C_FLOW_END
	C_FLOW_EXIT

	C_IO_WRITE_CH
	C_IO_WRITE_NUM
	C_IO_READ_CH
	C_IO_READ_NUM
)

func main() {
	err := Main(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func Main(args []string) error {
	if len(args) != 2 {
		return errors.New("Usage: gows PATH_TO_WS")
	}
	wsPath := args[1]
	f, err := os.Open(wsPath)
	if err != nil {
		return err
	}
	defer f.Close()
	p := NewParser(f)
	cs, err := p.Parse()
	if err != nil {
		return err
	}

	vm := NewVM(cs)
	err = vm.Eval(os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	return nil
}

type Command struct {
	Type    CommandType
	Operand int
}
