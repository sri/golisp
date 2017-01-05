#### A simple lisp interpreter in Go

### How to run it

```
git clone https://github.com/sri/golisp.git
cd golisp
export GOPATH=`pwd`
go build && ./golisp
```

### Constants

```lisp
LISP> nil
nil
LISP> ()
nil
LISP> 1
1
LISP> 1.99
1.99
LISP> 'hello-world
hello-world
LISP> '(1 2 3)
(1 2 3)
LISP> (let (a 1) (+ a b c))
unidentified symbol: b
```

### let Binding

```lisp
LISP> (let (a 1 b 2 c 3) (+ a b c))
=> 6
```

### String Concatenation

```lisp
LISP> (s+ "Hello, " "world!")
"Hello, world!"
```

### Lambdas

```lisp
LISP> ((lambda (a b) (+ a b)) 1 2)
3
LISP> (lambda () 'a)
(lambda nil (quote a))
LISP> ((lambda () 'a))
a
```

### If expressions

```lisp
LISP> (if 1 2 3)
2
LISP> (if nil 2 3)
3
```