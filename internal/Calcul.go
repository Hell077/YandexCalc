package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func isOperator(c string) bool {
	operators := map[string]struct{}{"+": {}, "-": {}, "*": {}, "/": {}}
	_, ok := operators[c]
	return ok
}

func isDigit(c string) bool {
	_, err := strconv.ParseFloat(c, 64)
	return err == nil
}

func precedence(op string) int {
	precedences := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}
	return precedences[op]
}

func toRPN(expression string) ([]string, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	outputQueue := []string{}
	operatorStack := []string{}

	for i := 0; i < len(expression); i++ {
		char := string(expression[i])

		if isDigit(char) {
			num := char
			for i+1 < len(expression) && (isDigit(string(expression[i+1])) || expression[i+1] == '.') {
				i++
				num += string(expression[i])
			}
			outputQueue = append(outputQueue, num)
		} else if char == "(" {
			operatorStack = append(operatorStack, char)
		} else if char == ")" {
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != "(" {
				outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) == 0 {
				return nil, errors.New("mismatched parentheses")
			}
			operatorStack = operatorStack[:len(operatorStack)-1]
		} else if isOperator(char) {
			for len(operatorStack) > 0 {
				top := operatorStack[len(operatorStack)-1]
				if isOperator(top) && precedence(top) >= precedence(char) {
					outputQueue = append(outputQueue, top)
					operatorStack = operatorStack[:len(operatorStack)-1]
				} else {
					break
				}
			}
			operatorStack = append(operatorStack, char)
		} else {
			return nil, fmt.Errorf("unknown character: %s", char)
		}
	}

	for len(operatorStack) > 0 {
		top := operatorStack[len(operatorStack)-1]
		if top == "(" {
			return nil, errors.New("mismatched parentheses")
		}
		outputQueue = append(outputQueue, top)
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	return outputQueue, nil
}

func evaluateRPN(rpn []string) (float64, error) {
	stack := []float64{}

	for _, token := range rpn {
		if isDigit(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else if isOperator(token) {
			if len(stack) < 2 {
				return 0, errors.New("invalid expression: not enough operands")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, errors.New("division by zero")
				}
				stack = append(stack, a/b)
			}
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression: incorrect number of operands")
	}

	return stack[0], nil
}

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	rpn, err := toRPN(expression)
	if err != nil {
		return 0, err
	}
	return evaluateRPN(rpn)
}
