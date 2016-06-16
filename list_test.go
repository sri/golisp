package main

import (
	"testing"
)

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
