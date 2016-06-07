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

type LispGoFn func(*LispList) LispObject

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

var NIL *LispList = new(LispList)

type LispSymbol struct {
	name string
}

var SYMBOLS = make(map[string]LispSymbol)

func InitSymbols() {
	SYMBOLS["if"] = LispSymbol{"if"}
	SYMBOLS["quote"] = LispSymbol{"quote"}
	SYMBOLS["lambda"] = LispSymbol{"lambda"}
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
			reader.Discard(1)
			result = append(result, b)
		}
	}

	s := string(result)
	if s == "nil" {
		return NIL
	} else if val, ok := SYMBOLS[s]; ok {
		return val
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

func GlobalEnv() *LispEnv {
	env := new(LispEnv)
	env.current = make(map[LispSymbol]LispObject)
	env.current[LispSymbol{"+"}] = LispGoFn(LispFn_Add)
	return env
}

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
			panic("Invalid arg to +: " + LispObject2String(n))
		}
		args = args.Rest()
	}

	return LispObject(sum)
}

func LispFn_Print(args ...LispObject) {
	output := []string{}
	for i := 0; i < len(args); i++ {
		output = append(output, LispObject2String(args[i]))
	}
	fmt.Println(output)
}

func Repl() {
	reader := bufio.NewReader(os.Stdin)
	env := GlobalEnv()

	for {
		fmt.Print("LISP> ")
		lispObj := Read(reader)
		result := Eval(lispObj, env)
		fmt.Println(LispObject2String(result))
	}
}

func main() {
	InitSymbols()
	Repl()
}
