package main

import (
	"fmt"
)

type LispGoFn func(*LispList) LispObject

func LispFn_Add(args *LispList) LispObject {
	sum := float64(0)
	for {
		if args == NIL {
			break
		}
		switch n := args.First().(type) {
		case float64:
			sum += n
		case int64:
			sum += float64(n)
		default:
			panic("Invalid arg to +: " + LispObject2String(n))
		}
		args = args.Rest()
	}

	return LispObject(sum)
}

func LispFn_Print(args ...LispObject) {
	output := []string{}
	for i := 0; i < len(args); i++ {
		output = append(output, LispObject2String(args[i]))
	}
	fmt.Println(output)
}
