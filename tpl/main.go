package tpl

func MainTemplate() []byte {
	return []byte(`package main

import (
	"{{ .ModName }}/execute"
	"flag"
	"os"
)

var (
	testcase  = flag.String("case", "", "specify the testcases to execute(separated by commas), such as(total:):")
	testsuite = flag.String("suite", "", "specify the testsuites to execute(separated by commas), such as(total:):")
)

func init() {
	flag.CommandLine.SetOutput(os.Stdout)
	flag.Parse()
}

func main() {
	execute.Execute(testcase, testsuite)
}
`)
}

func ExecuteTemplate() []byte {
	return []byte(`package execute

import (
	"fmt"
)

func Execute(testcase, testsuite *string) {
	fmt.Println("executing")
}
`)
}

func AddSuiteTemplate() []byte {
	return []byte(`package {{ .PackageName }}

// Here you will define your specific testcase
func (d *{{ .StructName }}) TestCase1() (err error) {
	return 
}

`)
}

func SuiteInitTemplate() []byte {
	return []byte(`package {{ .PackageName }}

import (
	"github.com/txy2023/potato/execute"
)

type {{ .StructName }} struct{}

func (d *{{ .StructName }}) Execute() {
	execute.Execute()
}

`)
}

func SuiteRegisteTemplate() []byte {
	return []byte(`package execute

import (
	"{{ .ImportedPath }}"

	"github.com/txy2023/potato/register"
)

func init() {
	register.Registe(new({{ .PackageName }}.{{ .StructName }}))
}
`)
}
