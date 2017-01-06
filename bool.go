package main

func IsTrue(obj LispObject) bool {
	if obj == LISP_NIL {
		return false
	}
	return true
}
