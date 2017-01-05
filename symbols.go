package main

type LispSymbol struct {
	name string
}

var SYMBOLS = make(map[string]LispSymbol)

func InitSymbols() {
	syms := []string{
		"backquote",
		"def",
		"if",
		"lambda",
		"let",
		"macro",
		"quote",
		"unquote",
		"unquote-splice",
	}

	for _, sym := range syms {
		SYMBOLS[sym] = LispSymbol{sym}
	}
}
