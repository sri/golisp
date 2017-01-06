package main

import (
	"bufio"
	"strconv"
)

func ReadList(reader *bufio.Reader) LispObject {
	result := []LispObject{}

	reader.Discard(1)
loop:
	for {
		b, err := reader.Peek(1)
		if err != nil {
			break
		}
		switch b[0] {
		case ' ', '\t', '\n':
			reader.Discard(1)
			continue
		case ')':
			reader.Discard(1)
			break loop
		default:
			lispObj, err := Read(reader)
			if err != nil {
				LispFatalError(err)
			}
			result = append(result, lispObj)
		}
	}

	list := LISP_NIL
	for i := len(result) - 1; i >= 0; i-- {
		list = Cons(result[i], list)
	}
	return list
}

func ReadString(reader *bufio.Reader) LispObject {
	delim, _ := reader.ReadByte()
	line, err := reader.ReadString(delim)
	if err != nil {
		LispFatalError(err)
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
			break
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

	if len(result) == 0 {
		LispError("Invalid read in ReadAtom")
	}

	s := string(result)
	if s == "nil" {
		return LISP_NIL
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
	lispObj, err := Read(reader)
	if err != nil {
		LispFatalError(err)
	}
	return List(SYMBOLS["quote"], lispObj)
}

func ReadBackquote(reader *bufio.Reader) LispObject {
	reader.Discard(1)
	lispObj, err := Read(reader)
	if err != nil {
		LispFatalError(err)
	}
	return List(SYMBOLS["backquote"], lispObj)
}

func ReadUnquote(reader *bufio.Reader) LispObject {
	reader.Discard(1)

	b, err := reader.Peek(1)
	if err != nil {
		LispFatalError(err)
	}

	head := SYMBOLS["unquote"]
	switch b[0] {
	case '@':
		reader.Discard(1)
		head = SYMBOLS["unquote-splice"]
	}

	lispObj, err := Read(reader)
	if err != nil {
		LispFatalError(err)
	}
	return List(head, lispObj)
}

func Read(reader *bufio.Reader) (LispObject, error) {
	for {
		b, err := reader.Peek(1)
		if err != nil {
			return LISP_NIL, err
		}

		switch b[0] {
		case ' ', '\n', '\t':
			// Skip whitespace.
			reader.Discard(1)
			continue
		case '(':
			return ReadList(reader), nil
		case '"':
			return ReadString(reader), nil
		case '\'':
			return ReadQuote(reader), nil
		case '`':
			return ReadBackquote(reader), nil
		case ',':
			return ReadUnquote(reader), nil
		default:
			return ReadAtom(reader), nil
		}
	}

	return LISP_NIL, nil
}
