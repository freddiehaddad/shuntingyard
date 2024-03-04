package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Token IDs
const (
	SPACE   = ' '
	PLUS    = '+'
	MINUS   = '-'
	LPAREN  = '('
	RPAREN  = ')'
	INTEGER = 'I'
	NEGATE  = 'N'
)

// Tokenized value from input.
type token struct {
	id    byte
	value int
}

// Simple stack data structure. Supports, push, pop, top, and empty operations.
type stack[T any] struct {
	data []T
}

// push the value v onto the stack.
func (s *stack[T]) push(v T) {
	s.data = append(s.data, v)
}

// pop a value from the stack.
func (s *stack[T]) pop() T {
	v := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return v
}

// top returns the value at the top of the stack without modifying it.
func (s *stack[T]) top() T {
	return s.data[len(s.data)-1]
}

// empty returns true if the stack is empty, false otherwise.
func (s *stack[T]) empty() bool {
	return len(s.data) == 0
}

// evalOperators performs arithmetic operations on the values based
// determined by the operators. The function assumes Reverse Polish
// notation.  For example, if '+' is at the top of the operator stack
// and 10 and 5 are the top two elements on the values stack, then
// the three values are popped, added together, and the result, 15,
// is pushed to the values stack. The operation continues until the
// operators stack is empty or a left parenthesis is encountered.
func evalOperators(operators *stack[byte], values *stack[int]) {
	var left, right, result int
	for !operators.empty() && operators.top() != LPAREN {
		right = values.pop()

		switch operators.pop() {
		case NEGATE:
			result = -right
		case MINUS:
			left = values.pop()
			result = left - right
		case PLUS:
			left = values.pop()
			result = left + right
		}
		values.push(result)
	}
}

// parse reads the input s and returns a slice of tokens with spaces being
// ignored. Valid input is expected as no error checking is performed.
func parse(s string) []*token {
	var tokens []*token
	var sb strings.Builder

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case SPACE:
			continue
		case MINUS:
			tokens = append(tokens, &token{MINUS, MINUS})
		case PLUS:
			tokens = append(tokens, &token{PLUS, PLUS})
		case LPAREN:
			tokens = append(tokens, &token{LPAREN, LPAREN})
		case RPAREN:
			tokens = append(tokens, &token{RPAREN, RPAREN})
		default:
			for i < len(s) && '0' <= s[i] && s[i] <= '9' {
				sb.WriteByte(s[i])
				i++
			}
			value, _ := strconv.ParseInt(sb.String(), 10, 32)
			tokens = append(tokens, &token{INTEGER, int(value)})
			sb.Reset()
			i--
		}
	}

	return tokens
}

// calculate evaluates the input s and returns the result.  Valid operations
// are addition and subtraction of positive and negative integers.  Evaluation
// is left-associative.  Precedence is given to expressions surrounded by
// parentheses.
func calculate(s string) int {
	values := &stack[int]{}
	operators := &stack[byte]{}

	tokens := parse(s)

	waitingForValue := true
	for _, token := range tokens {
		switch token.id {
		case INTEGER:
			values.push(token.value)
			waitingForValue = false
		case PLUS:
			evalOperators(operators, values)
			operators.push(token.id)
			waitingForValue = true
		case MINUS:
			// unary operator
			if waitingForValue {
				operators.push(NEGATE)
			} else {
				evalOperators(operators, values)
				operators.push(token.id)
				waitingForValue = true
			}
		case LPAREN:
			operators.push(token.id)
		case RPAREN:
			evalOperators(operators, values)
			// discard the "("
			operators.pop()
		}
	}

	evalOperators(operators, values)
	return values.pop()
}

func main() {
	s := "(1+(4+5+2)-3)+(6+8)"
	result := calculate(s)
	fmt.Println(s, "=", result)
}
