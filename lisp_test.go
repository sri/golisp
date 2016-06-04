package main

import (
	"testing"
	"bufio"
	"strings"
)

func TestReadList(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("(1 2 3)"))
	result := LispObject2String(ReadList(reader))

	if result != "(1 2 3)" {
		t.Errorf("%s != (1 2 3)", result)
	}
}
