package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestEval(t *testing.T) {
	expectations := []TestExpectation{
		{"nil", "nil"},
		{"()", "nil"},
		// If
		{"(if 10 20 30)", "20"},                // if true
		{"(if nil 20 30)", "30"},               // if false
		{"(if (if 10 nil 20) 400 500)", "500"}, // cond expr return nil
		{"(if (if 10 20 20) 400 500)", "400"},  // cond expr return true

		// Quote
		{"'(if 10 20 30)", "(if 10 20 30)"},

		// Lambda
		{"((lambda (a b c) (+ a b c)) 1 2 3)", "6"},
		// Hiding variable
		{"((if 1 (lambda (a) (+ a 2)) (lambda (a) (+ a 100))) 100)", "102"},
	}

	for _, exp := range expectations {
		reader := bufio.NewReader(strings.NewReader(exp.arg))
		result := Eval(Read(reader), GlobalEnv())
		actual := LispObject2String(result)
		if actual != exp.expected {
			t.Errorf("Test Eval: Expected (eval %s) => %s, but got %s",
				exp.arg, exp.expected, actual)
		}

	}
}
