package main

type LispSymbol struct {
	name string
}

var SYMBOLS = make(map[string]LispSymbol)

func InitSymbols() {
	syms := []string{
		"if",
		"lambda",
		"let",
		"macro",
		"quote",
	}

	for _, sym := range syms {
		SYMBOLS[sym] = LispSymbol{sym}
	}
}
