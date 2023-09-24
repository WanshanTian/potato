package tpl

func MainTemplate() []byte {
	return []byte(`package main

import (
	"{{ .ModName }}/execute"
)

func main() {
	execute.Execute()
}
`)
}

func ExecuteTemplate() []byte {
	return []byte(`package execute

import (
	"github.com/spf13/cobra"
	"github.com/txy2023/potato/execute"
	"github.com/txy2023/potato/register"
)

var rootCmd = &cobra.Command{
	Use: "go run main.go OR ./main(the compiled name)",
	Long: ` + "`" +
		`when flag testcase is specified, the testcase specified will be executed, so is flag testsuite
if flag testcase or flag testsuite is not specified, all testcases will be executed` + "`" + `,
	Run: func(cmd *cobra.Command, args []string) {
		if *TestCasesSpecified != "" {
			execute.TestCasesSpecified = TestCasesSpecified
			execute.IsExistTestcases()
			for _, testsuite := range execute.TestSuitesExecForTestCasesSpecified {
				testsuite.Execute()
			}
		}
		if *TestSuitesSpecified != "" {
			execute.TestSuitesSpecified = TestSuitesSpecified
			execute.IsExistTestSuites()
			for _, testsuite := range execute.TestSuitesExec {
				testsuite.Execute()
			}
		} 
		if *TestCasesSpecified == "" && *TestSuitesSpecified == "" {
			for _, testsuite := range register.Registered {
				testsuite.Execute()
			}
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

var (
	TestCasesSpecified  *string
	TestSuitesSpecified *string
)

func init() {
	TestCasesSpecified = rootCmd.Flags().StringP("testcase", "c", "", "specify the testcases to execute(separated by commas), such as(total:):")
	TestSuitesSpecified = rootCmd.Flags().StringP("testsuite", "s", "", "specify the testsuites to execute(separated by commas), such as(total:):")
}
`)
}

func AddSuiteTemplate() []byte {
	return []byte(`package {{ .PackageName }}

// Here you can define Setup of testsuite optionally
func (d *{{ .StructName }}) Setup() (err error) {
	// implementation
	return 
}	

// Here you can define Teardown of testsuite optionally
func (d *{{ .StructName }}) Teardown() (err error) {
	// implementation
	return
}	

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
	execute.Execute(d)
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
