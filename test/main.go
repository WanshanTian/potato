package main

import (
	"fmt"
	"os"
	"path"
)

func main() {
	path1, _ := os.Getwd()
	fmt.Println(path1)
	fmt.Println(path.Base(path1))
	fmt.Println(path.Join("/test/1/2", "../3"))
}
