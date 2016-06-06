package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type LispObject interface{}
type LispList struct {
	first LispObject
	rest  *LispList
}

var NIL *LispList = new(LispList)

type LispSymbol struct {
	name string
}

var SYMBOLS = make(map[string]LispSymbol)

func InitSymbols() {
	SYMBOLS["if"] = LispSymbol{"if"}
	SYMBOLS["quote"] = LispSymbol{"quote"}
}

func LispObject2String(obj LispObject) string {
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

func IsTrue(obj LispObject) bool {
	if obj == NIL {
		return false
	}
	return true
}

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

func Evalis(list *LispList) LispObject {
	if list == NIL {
		return NIL
	}

	switch obj := list.First().(type) {
	case LispSymbol:
		if obj == SYMBOLS["if"] {
			// (if <cond> <if-true> <if-false>)
			cond := list.Rest().First()
			body := list.Rest().Rest()
			if IsTrue(Eval(cond)) {
				return Eval(body.First())
			}
			return Eval(body.Rest().First())
		} else if obj == SYMBOLS["quote"] {
			// (quote (a b c)) => (a b c)
			// (quote z) => z
			return list.Rest().First()
		}
	default:
		panic("do not know how to evaluate " + LispObject2String(list))
		return NIL
	}

	panic("do not know how to evaluate " + LispObject2String(list))
	return NIL
}

func Eval(obj LispObject) LispObject {
	switch o := obj.(type) {
	case *LispList:
		return Evalis(o)
	default:
		return obj
	}
}

func ReadList(reader *bufio.Reader) LispObject {
	result := []LispObject{}

	reader.Discard(1)
loop:
	for {
		b, err := reader.Peek(1)
		if err != nil {
			panic(err)
		}
		switch b[0] {
		case ' ', '\t', '\n':
			reader.Discard(1)
			continue
		case ')':
			reader.Discard(1)
			break loop
		default:
			result = append(result, Read(reader))
		}
	}

	list := NIL
	for i := len(result) - 1; i >= 0; i-- {
		list = Push(result[i], list)
	}
	return list

}

func ReadString(reader *bufio.Reader) LispObject {
	delim, _ := reader.ReadByte()
	line, err := reader.ReadString(delim)
	if err != nil {
		panic(err)
	}
	return LispObject(line[:len(line)-1])
}

func IsValidIdentifier(b byte) bool {
	if ('0' <= b && b <= '9') ||
		('a' <= b && b <= 'z') ||
		('A' <= b && b <= 'Z') ||
		b == '_' || b == '-' ||
		b == '?' || b == '%' {
		return true
	}
	return false
}

func ReadAtom(reader *bufio.Reader) LispObject {
	result := []byte{}

	for {
		buf, err := reader.Peek(1)
		if err != nil {
			panic(err)
		}
		b := buf[0]
		if b == ' ' || b == '\n' || b == '\t' {
			reader.Discard(1)
			break
		} else if IsValidIdentifier(b) {
			reader.Discard(1)
			result = append(result, b)
		} else if b == '(' || b == ')' {
			break
		} else {
			panic("invalid char: " + string(b))
		}
	}

	s := string(result)
	if s == "nil" {
		return NIL
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return LispObject(f)
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return LispObject(i)
	}
	return LispSymbol{s}
}

func ReadQuote(reader *bufio.Reader) LispObject {
	reader.Discard(1)
	fmt.Println("in ReadQuote")
	return NewList(SYMBOLS["quote"], Read(reader))
}

func Read(reader *bufio.Reader) LispObject {
	for {
		b, err := reader.Peek(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		switch b[0] {
		case ' ', '\n', '\t':
			// Skip whitespace.
			reader.Discard(1)
			continue
		case '(':
			return ReadList(reader)
		case '"':
			return ReadString(reader)
		case '\'':
			return ReadQuote(reader)
		default:
			return ReadAtom(reader)
		}
	}

	return NIL
}

func Repl() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("LISP> ")
		lispObj := Read(reader)
		result := Eval(lispObj)
		fmt.Println(LispObject2String(result))
	}
}

func main() {
	InitSymbols()
	Repl()
}
