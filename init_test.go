package main

import (
	"os"
	"testing"
)

type TestExpectation struct {
	arg, expected string
}

func TestMain(m *testing.M) {
	InitSymbols()
	os.Exit(m.Run())
}
