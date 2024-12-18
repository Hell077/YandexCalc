package internal

import (
	"errors"
	"fmt"
	"testing"
)

func TestCalc(t *testing.T) {
	expressions := []string{
		"1+2/2",
		"5*5+5",
		"1+1",
	}
	expectedResults := []float64{
		2,  // 1 + (2 / 2) = 2.5
		30, // (5 * 5) + 5 = 30
		2,  // 1 + 1 = 2
	}

	for i, expr := range expressions {
		result, err := Calc(expr)
		if err != nil {
			t.Errorf("unexpected error for expression '%s': %v", expr, err)
			continue
		}

		if result != expectedResults[i] {
			t.Errorf("unexpected result for expression '%s': got %f, expected %f", expr, result, expectedResults[i])
		}
	}
}

func TestEvaluateRPN(t *testing.T) {
	tests := []struct {
		name     string
		rpn      []string
		expected float64
		err      error
	}{
		{
			name:     "simple addition",
			rpn:      []string{"2", "3", "+"},
			expected: 5,
			err:      nil,
		},
		{
			name:     "simple subtraction",
			rpn:      []string{"5", "3", "-"},
			expected: 2,
			err:      nil,
		},
		{
			name:     "simple multiplication",
			rpn:      []string{"2", "3", "*"},
			expected: 6,
			err:      nil,
		},
		{
			name:     "simple division",
			rpn:      []string{"6", "3", "/"},
			expected: 2,
			err:      nil,
		},
		{
			name:     "division by zero",
			rpn:      []string{"6", "0", "/"},
			expected: 0,
			err:      errors.New("division by zero"),
		},
		{
			name:     "invalid expression (not enough operands)",
			rpn:      []string{"2", "+"},
			expected: 0,
			err:      errors.New("invalid expression: not enough operands"),
		},

		{
			name:     "complex expression",
			rpn:      []string{"3", "4", "+", "2", "*", "7", "/"},
			expected: 2,
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluateRPN(tt.rpn)
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("expected error %v, got %v", tt.err, err)
			}
			if result != tt.expected {
				t.Errorf("expected result %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestToRPN(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   []string
		err        error
	}{
		{
			name:       "simple addition",
			expression: "2 + 3",
			expected:   []string{"2", "3", "+"},
			err:        nil,
		},
		{
			name:       "simple subtraction",
			expression: "5 - 3",
			expected:   []string{"5", "3", "-"},
			err:        nil,
		},
		{
			name:       "simple multiplication",
			expression: "2 * 3",
			expected:   []string{"2", "3", "*"},
			err:        nil,
		},
		{
			name:       "simple division",
			expression: "6 / 3",
			expected:   []string{"6", "3", "/"},
			err:        nil,
		},
		{
			name:       "expression with parentheses",
			expression: "(2 + 3) * 4",
			expected:   []string{"2", "3", "+", "4", "*"},
			err:        nil,
		},
		{
			name:       "expression with nested parentheses",
			expression: "(2 + 3) * (4 + 5)",
			expected:   []string{"2", "3", "+", "4", "5", "+", "*"},
			err:        nil,
		},
		{
			name:       "mismatched parentheses",
			expression: "(2 + 3 * 4",
			expected:   nil,
			err:        errors.New("mismatched parentheses"),
		},
		{
			name:       "unknown character",
			expression: "2 + 3 $ 4",
			expected:   nil,
			err:        fmt.Errorf("unknown character: $"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toRPN(tt.expression)
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("expected error %v, got %v", tt.err, err)
			}
			if !equal(result, tt.expected) {
				t.Errorf("expected result %v, got %v", tt.expected, result)
			}
		})
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestIsOperator(t *testing.T) {
	tests := []struct {
		operator string
		expected bool
	}{
		{"+", true},
		{"-", true},
		{"*", true},
		{"/", true},
		{"^", false}, // недопустимый оператор
		{"a", false}, // недопустимый оператор
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("isOperator(%s)", tt.operator), func(t *testing.T) {
			result := isOperator(tt.operator)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsDigit(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},   // валидное число
		{"3.14", true},  // валидное число с плавающей точкой
		{"abc", false},  // не число
		{"", false},     // пустая строка
		{"-123", true},  // отрицательное число
		{"+45.6", true}, // число с плюсом
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("isDigit(%s)", tt.input), func(t *testing.T) {
			result := isDigit(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestPrecedence(t *testing.T) {
	tests := []struct {
		operator string
		expected int
	}{
		{"+", 1},
		{"-", 1},
		{"*", 2},
		{"/", 2},
		{"^", 0}, // оператора ^ нет в мапе, по умолчанию возвращаем 0
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("precedence(%s)", tt.operator), func(t *testing.T) {
			result := precedence(tt.operator)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
