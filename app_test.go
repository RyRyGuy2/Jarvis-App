package main

import (
	"math"
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       float64
	}{
		{name: "addition", expression: "5+5", want: 10},
		{name: "precedence", expression: "2+3*4", want: 14},
		{name: "parentheses", expression: "((5*10)/3)+5", want: 65.0 / 3.0},
		{name: "unary minus", expression: "-(5+2)", want: -7},
		{name: "unary plus", expression: "+4", want: 4},
		{name: "decimal", expression: "1.5 * 2", want: 3},
		{name: "trailing equals", expression: "5 + 5 =", want: 10},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Calculate(test.expression)
			if err != nil {
				t.Fatalf("Calculate(%q) returned error: %v", test.expression, err)
			}
			if math.Abs(got-test.want) > 1e-12 {
				t.Fatalf("Calculate(%q) = %v, want %v", test.expression, got, test.want)
			}
		})
	}
}

func TestCalculateRejectsInvalidExpressions(t *testing.T) {
	tests := []struct {
		expression string
		wantError  string
	}{
		{expression: "", wantError: "empty calculation"},
		{expression: "1/0", wantError: "cannot divide by zero"},
		{expression: "(1+2", wantError: "missing closing parenthesis"},
		{expression: "1 2", wantError: "invalid expression"},
		{expression: "2+a", wantError: "invalid character"},
		{expression: "1..2", wantError: "invalid number"},
	}

	for _, test := range tests {
		t.Run(test.expression, func(t *testing.T) {
			_, err := Calculate(test.expression)
			if err == nil {
				t.Fatalf("Calculate(%q) unexpectedly succeeded", test.expression)
			}
			if !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("Calculate(%q) error = %q, want it to contain %q", test.expression, err, test.wantError)
			}
		})
	}
}
