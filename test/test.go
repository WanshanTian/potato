package main

import (
	"fmt"

	"github.com/txy2023/potato/execute"
)

type Hello2Suite struct{}

func (d *Hello2Suite) Execute() {
	execute.Execute(d)
}
func (d *Hello2Suite) Walk() error {
	fmt.Println("hello walk")
	return nil
}
func (d *Hello2Suite) Setup() error {
	fmt.Println("begin")
	return fmt.Errorf("errr")
}
func (d *Hello2Suite) Teardown(a int) error {
	fmt.Println("end")
	return nil
}
func main() {
	a := new(Hello2Suite)
	a.Execute()
}
