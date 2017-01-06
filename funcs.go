package main

import (
	"fmt"
	"strings"
)

type LispGoFn func(*LispList) LispObject

func LispFn_Add(args *LispList) LispObject {
	sum := float64(0)
	for {
		if args == LISP_NIL {
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

func LispFn_StringConcat(args *LispList) LispObject {
	s := []string{}
	for {
		if args == LISP_NIL {
			break
		}
		switch arg := args.First().(type) {
		case string:
			s = append(s, fmt.Sprintf("%s", arg))
		default:
			s = append(s, LispObject2String(arg))
		}

		args = args.Rest()
	}
	return LispObject(strings.Join(s, ""))
}

func LispFn_Print(args *LispList) LispObject {
	for {
		if args == LISP_NIL {
			break
		}

		switch arg := args.First().(type) {
		case string:
			fmt.Printf("%s ", arg)
		default:
			fmt.Printf("%s ", LispObject2String(arg))
		}

		args = args.Rest()
	}
	fmt.Println()
	return LISP_NIL
}
