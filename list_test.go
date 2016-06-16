package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestReverseList(t *testing.T) {
	expectations := []TestExpectation{
		{"()", "nil"},
		{"(1)", "(1)"},
		{"(1 2 3)", "(3 2 1)"},
	}

	for _, exp := range expectations {
		reader := bufio.NewReader(strings.NewReader(exp.arg))
		result := Read(reader).(*LispList)
		actual := LispObject2String(ReverseList(result))
		if actual != exp.expected {
			t.Errorf("TestReverseList: Expected %v => %s, but got %s",
				exp.arg, exp.expected, actual)
		}

	}
}

func TestNth(t *testing.T) {
	a := NewList(1, 2, 3)

	if a.Nth(0) != NIL {
		t.Errorf("Nth(0)")
	}

	if a.Nth(1) != 1 {
		t.Errorf("Nth(1)")
	}

	if a.Nth(3) != 3 {
		t.Errorf("Nth(3)")
	}

	if a.Nth(100) != NIL {
		t.Errorf("Nth(100)")
	}
}
