package main

import (
	"strings"
)

type LispList struct {
	first LispObject
	rest  *LispList
}

var NIL *LispList = new(LispList)

// Methods on LispList:
func (list *LispList) First() LispObject {
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
		result = append(result, LispObject2String(list.First()))
		list = list.Rest()
		if list != NIL {
			result = append(result, " ")
		}
	}
	result = append(result, ")")

	return strings.Join(result, "")
}

func Push(first LispObject, rest *LispList) *LispList {
	result := new(LispList)
	result.first = first
	result.rest = rest
	return result
}

func NewList(args ...LispObject) *LispList {
	result := NIL
	for i := len(args) - 1; i >= 0; i-- {
		result = Push(args[i], result)
	}
	return result
}
