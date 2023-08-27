package tpl

func MainTemplate() []byte {
	return []byte(`package main

import (
	"{{ .PkgName }}/execute"
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
	return []byte(`package {{ .ToLowerSuiteBaseName }}



`)
}
