package main

import (
	"fmt"

	"github.com/txy2023/potato/register"
)

func main() {
	a := register.Registe(&register.TestSuite{})
	b := register.Registe(&register.TestSuite{})
	c := register.Registe(&register.TestSuite{})
	fmt.Println(a, b, c)
}
