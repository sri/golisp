package main

import (
	"bufio"
	"io"
	"strconv"
)

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
			if err == io.EOF {
				break
			}
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
