package main

type LispSymbol struct {
	name string
}

var SYMBOLS = make(map[string]LispSymbol)

func InitSymbols() {
	SYMBOLS["if"] = LispSymbol{"if"}
	SYMBOLS["quote"] = LispSymbol{"quote"}
	SYMBOLS["lambda"] = LispSymbol{"lambda"}
}
