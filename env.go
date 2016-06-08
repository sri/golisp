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

func MakeEnv(env *LispEnv, args *LispList, vals *LispList) *LispEnv {
	newEnv := new(LispEnv)
	newEnv.current = make(map[LispSymbol]LispObject)

	if env != nil {
		newEnv.parent = env
	}

	for {
		if args == NIL && vals == NIL {
			break
		} else if args == NIL || vals == NIL {
			panic("excess params: args=" +
				LispObject2String(args) +
				", vals=" + LispObject2String(vals))
		}

		switch arg := args.First().(type) {
		case LispSymbol:
			newEnv.current[arg] = Eval(vals.First(), env)
			break
		default:
			panic("Invalid obj, must be a symbol: " + LispObject2String(arg))

		}

		args = args.Rest()
		vals = vals.Rest()
	}

	return newEnv
}

func GlobalEnv() *LispEnv {
	env := new(LispEnv)
	env.current = make(map[LispSymbol]LispObject)
	env.current[LispSymbol{"+"}] = LispGoFn(LispFn_Add)
	return env
}
