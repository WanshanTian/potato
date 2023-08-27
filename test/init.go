package test

import (
	"fmt"

	"github.com/txy2023/potato/register"
)

type Dvs struct{}

func (d *Dvs) Execute() {
	fmt.Println("1")
}

func init() {
	register.Registe(new(Dvs))
}

func (*Dvs) VM() {
	fmt.Println(1)
}
