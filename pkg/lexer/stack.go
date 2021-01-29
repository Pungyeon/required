package lexer

import "fmt"

type Stack struct {
	index int
	stack []byte
}

func NewStack(size int) *Stack {
	return &Stack{
		stack: make([]byte, size),
	}
}

func (s *Stack) Pop() byte {
	s.index--
	if s.index < 0 {
		fmt.Println("stack underflow")
		return 0
	}
	return s.stack[s.index+1]
}

func (s *Stack) Push(b byte) {
	s.index++
	if s.index >= len(s.stack) {
		s.stack = append(s.stack, b)
	} else {
		s.stack[s.index] = b
	}
}
