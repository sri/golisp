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
			LispError("Invalid arg to +: " + LispObject2String(n))
		}
		args = args.Rest()
	}

	return LispObject(sum)
}

// TODO: fix bug with string: prints out with quotes
func LispFn_Print(args *LispList) LispObject {
	for {
		if args == NIL {
			break
		}
		fmt.Printf("%s ", LispObject2String(args.First()))
		args = args.Rest()
	}
	fmt.Println()
	return NIL
}
