package main

import (
	"bufio"
	"fmt"
	"os"
	"io"
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
		lispObj, err := Read(reader)

		if err != nil {
			if err == io.EOF {
				fmt.Println()
				os.Exit(0)
			} else {
				LispError(err)
			}
		}

		result := Eval(lispObj, env)
		fmt.Println(LispObject2String(result))

	}
}

func main() {
	InitSymbols()
	Repl()
}
