package main

import (
	"strings"
)

type LispList struct {
	first LispObject
	rest  *LispList
}

var LISP_NIL *LispList = new(LispList)

// Methods on LispList:
func (list *LispList) First() LispObject {
	if list == LISP_NIL {
		return LISP_NIL
	}
	return list.first
}

func (list *LispList) Second() LispObject {
	return list.Rest().First()
}

func (list *LispList) Third() LispObject {
	return list.Rest().Rest().First()
}

// Note: Nth(1) == First() & Nth(2) == Second()
func (list *LispList) Nth(n int) LispObject {
	l := list
	for {
		if n <= 0 {
			return LISP_NIL
		} else if n == 1 {
			return l.First()
		}
		l = l.Rest()
		n--
	}
	return LISP_NIL
}

func (list *LispList) Rest() *LispList {
	if list == LISP_NIL {
		return LISP_NIL
	}
	return list.rest
}

func (list *LispList) String() string {
	result := []string{}

	if list == LISP_NIL {
		return "nil"
	}

	result = append(result, "(")
	for {
		if list == LISP_NIL {
			break
		}
		result = append(result, LispObject2String(list.First()))
		list = list.Rest()
		if list != LISP_NIL {
			result = append(result, " ")
		}
	}
	result = append(result, ")")

	return strings.Join(result, "")
}

func Cons(first LispObject, rest *LispList) *LispList {
	result := new(LispList)
	result.first = first
	result.rest = rest
	return result
}

func ReverseList(list *LispList) *LispList {
	ary := []LispObject{}
	for {
		if list == LISP_NIL {
			break
		}
		ary = append(ary, list.First())
		list = list.Rest()
	}
	result := LISP_NIL
	for i := 0; i < len(ary); i++ {
		result = Cons(ary[i], result)
	}
	return result
}

func List(args ...LispObject) *LispList {
	result := LISP_NIL
	for i := len(args) - 1; i >= 0; i-- {
		result = Cons(args[i], result)
	}
	return result
}
