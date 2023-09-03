package main

import (
	"fmt"

	"golang.org/x/mod/modfile"
)

func main() {
	src := `
	module github.com/you/hello
	
	require rsc.io/quote v1.5.2
	`

	mod := modfile.ModulePath([]byte(src))
	fmt.Println(mod)
}
