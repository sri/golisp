package main

func IsTrue(obj LispObject) bool {
	if obj == NIL {
		return false
	}
	return true
}
