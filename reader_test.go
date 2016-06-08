package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	expectations := []TestExpectation{
		{"\n\t", "nil"},
		{"\"hello world\"", "\"hello world\""},
		{"()", "nil"}, // empty list
		{"1", "1"},
		{"'(1 2 3)", "(quote (1 2 3))"},
		// {"'   \r\n(1     2    3  )", "(1 2 3)"},
		{"(1 2 (3.1 3.2 3.3) 4 5 \"six\")", "(1 2 (3.1 3.2 3.3) 4 5 \"six\")"},
		{"1.34", "1.34"},
		{"(if 10 20 30)", "(if 10 20 30)"},
	}

	for _, exp := range expectations {
		reader := bufio.NewReader(strings.NewReader(exp.arg))
		result := Read(reader)
		actual := LispObject2String(result)
		if actual != exp.expected {
			t.Errorf("TestRead: Expected (eval %v) => %s, but got %s",
				exp.arg, exp.expected, actual)
		}

	}

}
