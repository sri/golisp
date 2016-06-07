package main

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func readTest(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	return LispObject2String(Read(reader))
}

func TestReadList(t *testing.T) {
	result := readTest("   (  1 2     3   )")
	if result != "(1 2 3)" {
		t.Errorf("%s != (1 2 3)", result)
	}
}

func TestReadString(t *testing.T) {
	result := readTest("\"hello world\"")
	if result != "\"hello world\"" {
		t.Errorf("test read string: %s != \"hello world\"", result)
	}
}

func TestEvalEmptyList(t *testing.T) {
	result := LispObject2String(Eval(NewList(), new(LispEnv)))
	if result != "nil" {
		t.Errorf("eval empty list: %v != nil", result)
	}
}

func TestListToString(t *testing.T) {
	list := NewList(1, 2, NewList(3.1, 3.2, 3.3), 4, 5, "six")
	result := list.String()
	if result != "(1 2 (3.1 3.2 3.3) 4 5 \"six\")" {
		t.Errorf("list to string: %v != (1 2 (3.1 3.2 3.3) 4 5 \"six\")", result)
	}
}

func TestEvalIfTrue(t *testing.T) {
	list := NewList(SYMBOLS["if"], 2, 3, 4)
	result := LispObject2String(Eval(list, new(LispEnv)))
	if result != "3" {
		t.Error("eval if true: doesn't eval to 3", list.String())
	}
}

func TestEvalIfFalse(t *testing.T) {
	list := NewList(SYMBOLS["if"], NIL, 3, 4)
	result := LispObject2String(Eval(list, new(LispEnv)))
	if result != "4" {
		t.Error("eval if false: doesn't eval to 4", list.String(), "=>", result)
	}
}

func TestEvalReadIfTrueExpr(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("(if (if 10 20 20) 400 500)"))
	result := Eval(Read(reader), new(LispEnv))
	s := LispObject2String(result)
	if s != "400" {
		t.Error("test eval read if true expr: doesn't eval to 400", "evals to =>", s)
	}
}

func TestEvalReadIfFalseExpr(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("(if (if 10 nil 20) 400 500)"))
	result := Eval(Read(reader), new(LispEnv))
	s := LispObject2String(result)
	if s != "500" {
		t.Error("test eval read if false expr: doesn't eval to 400", "evals to =>", s)
	}
}

func TestMain(m *testing.M) {
	InitSymbols()
	os.Exit(m.Run())
}
