package main

import (
	"fmt"
	"strings"
)

type LispObj interface{}

type LispList struct {
	first LispObj
	rest  *LispList
}

var NIL *LispList = new(LispList)

type LispSymbol struct {
	name string
}

var SYMBOLS = make(map[string]LispSymbol)

func initSymbols() {
	SYMBOLS["if"] = LispSymbol{ "if" }
}

func LispObj2String(obj LispObj) string {
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

func IsTrue(obj LispObj) bool {
	if obj == NIL {
		return false
	}
	return true
}

// Methods on LispList:
func (list *LispList) First() LispObj {
	if list == NIL {
		return NIL
	}
	return list.first
}

func (list *LispList) Rest() *LispList {
	if list == NIL {
		return NIL
	}
	return list.rest
}

func (list *LispList) String() string {
	result := []string{}

	if list == NIL {
		return "nil"
	}

	result = append(result, "(")
	for {
		if list == NIL {
			break
		}
		result = append(result, LispObj2String(list.First()))
		list = list.Rest()
		if list != NIL {
			result = append(result, " ")
		}
	}
	result = append(result, ")")

	return strings.Join(result, "")
}

func Push(first LispObj, rest *LispList) *LispList {
	result := new(LispList)
	result.first = first
	result.rest = rest
	return result
}

func NewList(args ...LispObj) *LispList {
	result := NIL
	for i := len(args) - 1; i >= 0; i-- {
		result = Push(args[i], result)
	}
	return result
}

func Evalis(list *LispList) LispObj {
	if list == NIL {
		return NIL
	}

	switch obj := list.First().(type) {
	case LispSymbol:
		// Yuck -- specials forms are strings
		// for now.
		if obj == SYMBOLS["if"] {
			// ("if" <cond> <if-true> <if-false>)
			cond := list.Rest().First()
			body := list.Rest().Rest()
			if IsTrue(Eval(cond)) {
				return Eval(body.First())
			}
			return Eval(body.Rest().First())
		}
	default:
		return NIL
	}
	return NIL
}

func Eval(obj LispObj) LispObj {
	switch o := obj.(type) {
	case *LispList:
		return Evalis(o)
	default:
		return obj
	}
}

func main() {
	initSymbols()

	empty := NewList()
	result := Eval(empty)
	fmt.Printf("eval(empty list) => %s\n", LispObj2String(result))

	list := NewList(1, 2, NewList(3.1, 3.2, 3.3), 4, 5, "six")
	fmt.Printf("printing out list: %s\n", list.String())

	list = NewList(SYMBOLS["if"], 2, 3, 4)
	result = Eval(list)
	fmt.Printf("%s expr should eval to 3 => %s\n",
		list.String(),
		LispObj2String(result))

	list = NewList(SYMBOLS["if"], NIL, 3, 4)
	result = Eval(list)
	fmt.Printf("%s expr should eval to 4 => %s\n",
		list.String(),
		LispObj2String(result))
}
