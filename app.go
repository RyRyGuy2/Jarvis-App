package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// calculator evaluates expressions such as:
// 5+5
// (15*10)/10
// ((5*10)/3)+5
func (a *App) Calculator(calculation string) (float64, error) {
	return Calculate(calculation)
}

func Calculate(calculation string) (float64, error) {
	tokens, err := tokenize(calculation)
	if err != nil {
		return 0, err
	}

	result, nextIndex, err := evaluateExpression(tokens, 0)
	if err != nil {
		return 0, err
	}

	if nextIndex != len(tokens) {
		return 0, errors.New("invalid expression")
	}

	return result, nil
}

// Converts:
//
// ((5 * 10) / 3) + 5
//
// into:
//
// ["(", "(", "5", "*", "10", ")", "/", "3", ")", "+", "5"]
func tokenize(calculation string) ([]string, error) {
	var tokens []string
	var currentNumber strings.Builder

	flushNumber := func() {
		if currentNumber.Len() > 0 {
			tokens = append(tokens, currentNumber.String())
			currentNumber.Reset()
		}
	}

	for _, char := range calculation {

		if unicode.IsSpace(char) {
			// Whitespace can separate tokens, but it must not silently join
			// two numbers (for example, "1 2" must not become "12").
			flushNumber()
			continue
		}

		if unicode.IsDigit(char) || char == '.' {
			currentNumber.WriteRune(char)
			continue
		}

		flushNumber()

		switch char {
		case '+', '-', '*', '/', '(', ')':
			tokens = append(tokens, string(char))

		case '=':
			// Allow something like:
			// 5 + 5 =
			continue

		default:
			return nil, fmt.Errorf("invalid character: %c", char)
		}
	}

	flushNumber()

	if len(tokens) == 0 {
		return nil, errors.New("empty calculation")
	}

	return tokens, nil
}

// expression handles:
// + and -
func evaluateExpression(tokens []string, index int) (float64, int, error) {
	left, index, err := evaluateTerm(tokens, index)
	if err != nil {
		return 0, index, err
	}

	for index < len(tokens) {

		operator := tokens[index]

		if operator != "+" && operator != "-" {
			break
		}

		index++

		right, nextIndex, err := evaluateTerm(tokens, index)
		if err != nil {
			return 0, index, err
		}

		if operator == "+" {
			left += right
		} else {
			left -= right
		}

		index = nextIndex
	}

	return left, index, nil
}

// term handles:
// * and /
func evaluateTerm(tokens []string, index int) (float64, int, error) {
	left, index, err := evaluateFactor(tokens, index)
	if err != nil {
		return 0, index, err
	}

	for index < len(tokens) {

		operator := tokens[index]

		if operator != "*" && operator != "/" {
			break
		}

		index++

		right, nextIndex, err := evaluateFactor(tokens, index)
		if err != nil {
			return 0, index, err
		}

		if operator == "*" {
			left *= right
		} else {

			if right == 0 {
				return 0, index, errors.New("cannot divide by zero")
			}

			left /= right
		}

		index = nextIndex
	}

	return left, index, nil
}

// factor handles:
// numbers
// parentheses
// negative numbers
func evaluateFactor(tokens []string, index int) (float64, int, error) {

	if index >= len(tokens) {
		return 0, index, errors.New("expected number")
	}

	token := tokens[index]

	// Handle negative numbers:
	// -5
	// -(5+2)
	if token == "-" {

		value, nextIndex, err := evaluateFactor(tokens, index+1)
		if err != nil {
			return 0, index, err
		}

		return -value, nextIndex, nil
	}

	// Handle positive prefix:
	// +5
	if token == "+" {

		value, nextIndex, err := evaluateFactor(tokens, index+1)
		if err != nil {
			return 0, index, err
		}

		return value, nextIndex, nil
	}

	// Parentheses
	if token == "(" {

		value, nextIndex, err := evaluateExpression(tokens, index+1)
		if err != nil {
			return 0, index, err
		}

		if nextIndex >= len(tokens) || tokens[nextIndex] != ")" {
			return 0, nextIndex, errors.New("missing closing parenthesis")
		}

		return value, nextIndex + 1, nil
	}

	// A stray closing parenthesis means the current
	// expression is finished.
	if token == ")" {
		return 0, index, errors.New("unexpected closing parenthesis")
	}

	value, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return 0, index, fmt.Errorf("invalid number: %s", token)
	}

	return value, index + 1, nil
}

// Greet returns a greeting for the given name
func (a *App) Greet(command string) string {

	if strings.EqualFold(command, "Add") {

		return strconv.Itoa(Add(2, 2))

	} else if strings.EqualFold(command, "Subtract") {

		return strconv.Itoa(Subtract(2, 2))

	} else if strings.EqualFold(command, "Multiply") {

		return strconv.Itoa(Multiply(2, 2))

	} else if strings.EqualFold(command, "Divide") {

		return strconv.FormatFloat(Divide(2, 3), 'f', 4, 64)

	} else if strings.EqualFold(command, "Factorial") {

		result, err := Factorial(4)

		if err != nil {
			return err.Error()
		}

		return strconv.Itoa(result)
	}

	return "Provide A Valid Command from the options: \n Add \n Subtract \n Multiply \n Divide"
}

func FormatInput(input string) (string, int, int) {
	words := strings.Fields(input)

	integers := []int{}
	var command string

	for _, word := range words {

		number, err := strconv.Atoi(word)

		if err != nil {
			command = word
			continue
		}

		integers = append(integers, number)
	}

	// Prevent index out of range crashes.
	if len(integers) < 2 {
		return command, 0, 0
	}

	return command, integers[0], integers[1]
}

func Add(a int, b int) int {
	return a + b
}

func Subtract(a int, b int) int {
	return a - b
}

func Multiply(a int, b int) int {
	return a * b
}

func Divide(a int, b int) float64 {
	return float64(a) / float64(b)
}

func Factorial(a int) (int, error) {

	if a < 1 {
		return 0, errors.New("factorial requires value of at least 1")
	}

	total := 1

	for i := 2; i <= a; i++ {
		total *= i
	}

	return total, nil
}
