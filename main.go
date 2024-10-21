package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	res, err := Calc("1+2*6/23")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	outputQueue := []string{}
	operatorStack := []string{}

	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	isOperator := func(c string) bool {
		_, ok := precedence[c]
		return ok
	}

	isDigit := func(c string) bool {
		_, err := strconv.ParseFloat(c, 64)
		return err == nil
	}

	var appendOperatorToOutput = func(op string) {
		for len(operatorStack) > 0 {
			top := operatorStack[len(operatorStack)-1]
			if isOperator(top) && precedence[top] >= precedence[op] {
				outputQueue = append(outputQueue, top)
				operatorStack = operatorStack[:len(operatorStack)-1]
			} else {
				break
			}
		}
		operatorStack = append(operatorStack, op)
	}

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
				return 0, errors.New("неверное выражение: отсутствующая открывающая скобка")
			}
			operatorStack = operatorStack[:len(operatorStack)-1]
		} else if isOperator(char) {
			appendOperatorToOutput(char)
		} else {
			return 0, fmt.Errorf("неизвестный символ: %s", char)
		}
	}

	for len(operatorStack) > 0 {
		top := operatorStack[len(operatorStack)-1]
		if top == "(" {
			return 0, errors.New("неверное выражение: отсутствующая закрывающая скобка")
		}
		outputQueue = append(outputQueue, top)
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	stack := []float64{}

	for _, token := range outputQueue {
		if isDigit(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else if isOperator(token) {
			if len(stack) < 2 {
				return 0, errors.New("неверное выражение: недостаточно операндов")
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
					return 0, errors.New("деление на ноль")
				}
				stack = append(stack, a/b)
			}
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("неверное выражение: недопустимое количество операндов")
	}

	return stack[0], nil
}
