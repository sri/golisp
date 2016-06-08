package main

import (
	"fmt"
	"os"
)

func Apply(obj LispObject, actualArgs *LispList, env *LispEnv) LispObject {
	switch fn := obj.(type) {
	case *LispList:
		if fn.First() != SYMBOLS["lambda"] {
			panic("lambdas only for now")
		}
		switch fnArgs := fn.Rest().First().(type) {
		case *LispList:
			fnBody := fn.Rest().Rest().First()
			newEnv := MakeEnv(env, fnArgs, actualArgs)
			return Eval(fnBody, newEnv)
		default:
			panic("function params need to be a list: " +
				LispObject2String(obj))

		}
	case LispGoFn:
		return fn(actualArgs)
	default:
		panic("currently only lambdas are supported")

	}

	return NIL
}

func EvalList(list *LispList, env *LispEnv) LispObject {
	if list == NIL {
		return NIL
	}

	switch obj := list.First().(type) {
	case LispSymbol:
		if obj == SYMBOLS["if"] {
			// (if <cond> <if-true> <if-false>)
			cond := list.Rest().First()
			body := list.Rest().Rest()
			if IsTrue(Eval(cond, env)) {
				return Eval(body.First(), env)
			}
			return Eval(body.Rest().First(), env)
		} else if obj == SYMBOLS["quote"] {
			// (quote (a b c)) => (a b c)
			// (quote z) => z
			return list.Rest().First()
		} else if obj == SYMBOLS["lambda"] {
			// (lambda (a b c) (+ a b c))
			// Lambdas evaluate to themselves
			return list
		}
	}

	fn := Eval(list.First(), env)
	result := []LispObject{}
	for args := list.Rest(); args != NIL; args = args.Rest() {
		result = append(result, Eval(args.First(), env))
	}
	args := NIL
	for i := len(result) - 1; i >= 0; i-- {
		args = Push(result[i], args)
	}
	return Apply(fn, args, env)
}

func EvalSymbol(sym LispSymbol, env *LispEnv) LispObject {
	for {
		if env == nil {
			break
		}
		if val, ok := env.current[sym]; ok {
			return val
		}
		env = env.parent
	}

	fmt.Println("unidentified symbol", sym.name)
	env.Print()
	os.Exit(1)
	return NIL
}

func Eval(obj LispObject, env *LispEnv) LispObject {
	switch o := obj.(type) {
	case *LispList:
		return EvalList(o, env)
	case LispSymbol:
		return EvalSymbol(o, env)
	default:
		return obj
	}
}
