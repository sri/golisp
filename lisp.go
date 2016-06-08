package main

import (
	"bufio"
	"fmt"
	"os"
)

type LispObject interface{}

func LispObject2String(obj LispObject) string {
	switch o := obj.(type) {
	case string:
		return fmt.Sprintf("%#v", obj)
	case LispSymbol:
		return o.name
	case *LispList:
		return o.String()
	default:
		return fmt.Sprintf("%v", obj)
	}

}

func Repl() {
	reader := bufio.NewReader(os.Stdin)
	env := GlobalEnv()

	for {
		fmt.Print("LISP> ")
		lispObj := Read(reader)
		result := Eval(lispObj, env)
		fmt.Println(LispObject2String(result))
	}
}

func main() {
	InitSymbols()
	Repl()
}
