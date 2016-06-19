package main

import (
	"fmt"
	"os"
	"errors"
)

func LispError(s string) error {
	return errors.New(s)
}

func LispFatalError(err interface{}) {
	fmt.Println("***", err)
	os.Exit(1)
}
