package main

import (
	"fmt"
)

// Convert
// (let (a 10 b 20 c 30) ...) =>
// ((lambda (a b c) ...) 10 20 30)
func Let2Lambda(let *LispList, env *LispEnv) *LispList {
	lambdaArgs := NIL
	actualArgs := NIL
	body := let.Third().(*LispList)

	letExpr := let.Second().(*LispList)
	for {
		if letExpr == NIL {
			break
		}

		lambdaArgs = Push(letExpr.First(), lambdaArgs)
		actualArgs = Push(letExpr.Second(), actualArgs)

		letExpr = letExpr.Rest().Rest()
	}

	lambdaArgs = ReverseList(lambdaArgs)
	actualArgs = ReverseList(actualArgs)

	result := Push(NewList(SYMBOLS["lambda"], lambdaArgs, body), actualArgs)
	return result
}

func Apply(obj LispObject, actualArgs *LispList, env *LispEnv) (LispObject, error) {
	switch fn := obj.(type) {
	case *LispList:
		head := fn.First()
		if head == SYMBOLS["lambda"] {
			switch fnArgs := fn.Second().(type) {
			case *LispList:
				fnBody := fn.Third()
				newEnv, err := MakeEnv(env, fnArgs, actualArgs)
				if err != nil {
					return NIL, err
				}
				return Eval(fnBody, newEnv)
			default:
				return NIL, LispError("function params need to be a list: " +
					LispObject2String(obj))
			}
		} else if head == SYMBOLS["macro"] {
			macroBody := fn.Third().(*LispList)
			expansion, err := Eval(macroBody, env)
			if err != nil {
				return NIL, err
			}
			fmt.Printf("MACRO EXPANSION: %s => %s\n",
				macroBody, expansion)
			return Eval(expansion, env)
		} else {
			return NIL, LispError("Unknown obj: " + LispObject2String(obj))
		}
	case LispGoFn:
		return fn(actualArgs), nil
	default:
		return NIL, LispError("currently only lambdas are supported")
	}

	return NIL, nil
}

func EvalList(list *LispList, env *LispEnv) (LispObject, error) {
	if list == NIL {
		return NIL, nil
	}

	switch obj := list.First().(type) {
	case LispSymbol:
		if obj == SYMBOLS["if"] {
			// (if <cond> <if-true> <if-false>)
			cond := list.Second()
			body := list.Rest().Rest()
			result, err := Eval(cond, env)
			if err != nil {
				return NIL, err
			}
			if IsTrue(result) {
				return Eval(body.First(), env)
			}
			return Eval(body.Second(), env)
		} else if obj == SYMBOLS["quote"] {
			// (quote (a b c)) => (a b c)
			// (quote z) => z
			return list.Second(), nil
		} else if obj == SYMBOLS["lambda"] {
			// (lambda (a b c) (+ a b c))
			// Lambdas evaluate to themselves
			return list, nil
		} else if obj == SYMBOLS["let"] {
			return Eval(Let2Lambda(list, env), env)
		} else if obj == SYMBOLS["macro"] {
			// simple macros:
			// (macro <name> (<macro-args>) <expansion> <body...>)
			// The macro gets transformed into a macro function.
			name := list.Second().(LispSymbol)
			macroArgs := list.Third().(*LispList)
			expansion := list.Nth(4).(*LispList)
			body := list.Nth(5).(*LispList)

			macroFn := NewList(SYMBOLS["macro"], macroArgs, expansion)
			newEnv, err := MakeEnv(env, NewList(name), NewList(NewList(SYMBOLS["quote"], macroFn)))
			if err != nil {
				return NIL, err
			}

			return Eval(body, newEnv)
		} else if obj == SYMBOLS["backquote"] {
			return list, nil
		}
	}

	fn, err := Eval(list.First(), env)
	if err != nil {
		return NIL, err
	}
	result := []LispObject{}
	for args := list.Rest(); args != NIL; args = args.Rest() {
		t, err := Eval(args.First(), env)
		if err != nil {
			return NIL, err
		}
		result = append(result, t)
	}
	args := NIL
	for i := len(result) - 1; i >= 0; i-- {
		args = Push(result[i], args)
	}
	return Apply(fn, args, env)
}

func EvalSymbol(sym LispSymbol, env *LispEnv) (LispObject, error) {
	for {
		if env == nil {
			break
		}
		if val, ok := env.current[sym]; ok {
			return val, nil
		}
		env = env.parent
	}

	return NIL, LispError("unidentified symbol: " + sym.name)
}

func Eval(obj LispObject, env *LispEnv) (LispObject, error) {
	switch o := obj.(type) {
	case *LispList:
		return EvalList(o, env)
	case LispSymbol:
		return EvalSymbol(o, env)
	default:
		return obj, nil
	}
}
