package main

import (
	"fmt"

	"github.com/txy2023/potato/execute"
)

// hellotest
type Hello2Suite struct{}

func (d *Hello2Suite) Execute() {
	execute.Execute(d)
}
func (d *Hello2Suite) Walk() error {
	// fmt.Println("hello walk")
	return nil
}
func (d *Hello2Suite) Walk2() error {
	// fmt.Println("hello walk")
	return nil
}
func (d *Hello2Suite) Setup() error {
	// fmt.Println("begin")
	return nil
}
func (d *Hello2Suite) Teardown() error {
	// fmt.Println("end")
	return fmt.Errorf("12")
}
func main() {
	// fmt.Println(utils.GetTestSuiteName(&Hello2Suite{}))
	// utils.GetAllTestCasesComment("../")
	// fmt.Println(utils.GetMethodsImplementedByUser(reflect.TypeOf(&Hello2Suite{})))
	a := Hello2Suite{}
	a.Execute()
}
