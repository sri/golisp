package main

import (
	"fmt"
)

type LispEnv struct {
	current map[LispSymbol]LispObject
	parent  *LispEnv
}

func (env *LispEnv) Print() {
	for k, v := range env.current {
		fmt.Printf(" %s => %s\n", LispObject2String(k), LispObject2String(v))
	}
	fmt.Println("===")
	if env.parent == nil {
		fmt.Println("(parent is NIL)")
	} else {
		env.parent.Print()
	}
}

func (env *LispEnv) Def(sym LispSymbol, result LispObject) {
	if val, ok := env.current[sym]; ok {
		fmt.Println(
			"(Warning: shadowing var: " +
				LispObject2String(sym) +
				", current val: " + LispObject2String(val) + ")")
	}
	env.current[sym] = result
}

func MakeEnv(env *LispEnv, args *LispList, vals *LispList) (*LispEnv, error) {
	newEnv := new(LispEnv)
	newEnv.current = make(map[LispSymbol]LispObject)

	if env != nil {
		newEnv.parent = env
	}

	for {
		if args == NIL && vals == NIL {
			break
		} else if args == NIL || vals == NIL {
			LispFatalError("excess params: args=" +
				LispObject2String(args) +
				", vals=" + LispObject2String(vals))
		}

		switch arg := args.First().(type) {
		case LispSymbol:
			result, err := Eval(vals.First(), env)
			if err != nil {
				return newEnv, err
			}
			newEnv.current[arg] = result
			break
		default:
			LispFatalError("Invalid obj, must be a symbol: " + LispObject2String(arg))

		}

		args = args.Rest()
		vals = vals.Rest()
	}

	return newEnv, nil
}

func GlobalEnv() *LispEnv {
	env := new(LispEnv)
	env.current = make(map[LispSymbol]LispObject)
	env.current[LispSymbol{"+"}] = LispGoFn(LispFn_Add)
	env.current[LispSymbol{"print"}] = LispGoFn(LispFn_Print)
	env.current[LispSymbol{"s+"}] = LispGoFn(LispFn_StringConcat)
	return env
}
