package main

import (
	"fmt"
	"os"
)

func LispError(err interface{}) {
	fmt.Println("***", err)
	os.Exit(1)
}
