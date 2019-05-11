package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/pkg/errors"
)

type VM struct {
	stack     []int
	heap      map[int]int
	callStack []int
	index     int
	labels    map[int]int
	cs        []Command
}

func NewVM(cs []Command) *VM {
	vm := &VM{
		stack:     []int{},
		heap:      map[int]int{},
		callStack: []int{},
		index:     0,
		labels:    map[int]int{},
		cs:        cs,
	}
	return vm
}

func (v *VM) Eval(r io.Reader, w io.Writer) error {
	v.prepareLabels()
	sc := bufio.NewScanner(r)

	for {
		c := v.cs[v.index]

		// fmt.Println(v.stack)
		// fmt.Println(v.heap)
		// fmt.Printf("index: %d\n", v.index)
		// fmt.Println("--------------------")
		// fmt.Printf("%+v\n", c)

		switch c.Type {
		case C_STACK_PUSH:
			v.stack = append(v.stack, c.Operand)
		case C_STACK_DUP:
			v.stack = append(v.stack, v.stack[len(v.stack)-1])
		case C_STACK_SWAP:
			a := v.stack[len(v.stack)-1]
			b := v.stack[len(v.stack)-2]
			v.stack[len(v.stack)-1] = b
			v.stack[len(v.stack)-2] = a
		case C_STACK_POP:
			v.stack = v.stack[:len(v.stack)-1]
		case C_CALC_ADD:
			v.stack[len(v.stack)-2] = v.stack[len(v.stack)-2] + v.stack[len(v.stack)-1]
			v.stack = v.stack[:len(v.stack)-1]
		case C_CALC_SUB:
			v.stack[len(v.stack)-2] = v.stack[len(v.stack)-2] - v.stack[len(v.stack)-1]
			v.stack = v.stack[:len(v.stack)-1]
		case C_CALC_MULTI:
			v.stack[len(v.stack)-2] = v.stack[len(v.stack)-2] * v.stack[len(v.stack)-1]
			v.stack = v.stack[:len(v.stack)-1]
		case C_CALC_DIV:
			v.stack[len(v.stack)-2] = div(v.stack[len(v.stack)-2], v.stack[len(v.stack)-1])
			v.stack = v.stack[:len(v.stack)-1]
		case C_CALC_MOD:
			v.stack[len(v.stack)-2] = abs(v.stack[len(v.stack)-2]) % abs(v.stack[len(v.stack)-1])
			v.stack = v.stack[:len(v.stack)-1]
		case C_HEAP_SAVE:
			val := v.stack[len(v.stack)-1]
			addr := v.stack[len(v.stack)-2]
			v.heap[addr] = val
			v.stack = v.stack[:len(v.stack)-2]
		case C_HEAP_LOAD:
			addr := v.stack[len(v.stack)-1]
			v.stack[len(v.stack)-1] = v.heap[addr]
		case C_FLOW_DEF:
			// skip
		case C_FLOW_CALL:
			v.callStack = append(v.callStack, v.index)
			v.index = v.labels[c.Operand]
		case C_FLOW_JUMP:
			v.index = v.labels[c.Operand]
		case C_FLOW_JUMP_IF_ZERO:
			if v.stack[len(v.stack)-1] == 0 {
				v.index = v.labels[c.Operand]
			}
			v.stack = v.stack[:len(v.stack)-1]
		case C_FLOW_JUMP_IF_NEG:
			if v.stack[len(v.stack)-1] < 0 {
				v.index = v.labels[c.Operand]
			}
			v.stack = v.stack[:len(v.stack)-1]
		case C_FLOW_END:
			v.index = v.callStack[len(v.callStack)-1]
			v.callStack = v.callStack[:len(v.callStack)-1]
		case C_FLOW_EXIT:
			return nil
		case C_IO_WRITE_CH:
			b := v.stack[len(v.stack)-1]
			fmt.Fprintf(w, "%c", b)
			v.stack = v.stack[:len(v.stack)-1]
		case C_IO_WRITE_NUM:
			i := v.stack[len(v.stack)-1]
			fmt.Fprint(w, i)
			v.stack = v.stack[:len(v.stack)-1]
		case C_IO_READ_CH:
			addr := v.stack[len(v.stack)-1]
			sc.Split(bufio.ScanBytes)
			if !sc.Scan() {
				return errors.New("Bad input")
			}
			v.heap[addr] = int(sc.Bytes()[0])
			v.stack = v.stack[:len(v.stack)-1]
		case C_IO_READ_NUM:
			addr := v.stack[len(v.stack)-1]
			sc.Split(bufio.ScanLines)
			if !sc.Scan() {
				return errors.New("Bad input")
			}
			i, err := strconv.Atoi(sc.Text())
			if err != nil {
				return err
			}
			v.heap[addr] = i
			v.stack = v.stack[:len(v.stack)-1]
		default:
			return errors.New("Unreachable")
		}

		v.index += 1
	}
}

func (v *VM) prepareLabels() {
	for idx, c := range v.cs {
		if c.Type == C_FLOW_DEF {
			v.labels[c.Operand] = idx
		}
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func div(i, j int) int {
	return int(math.Floor(float64(i) / float64(j)))
}
