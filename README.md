#### A simple lisp interpreter in Go

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
