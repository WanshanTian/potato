package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/txy2023/potato/tpl"
)

// Project contains name and paths to projects.
type Project struct {
	PkgName      string
	AbsolutePath string
	Viper        bool
	ProjectName  string
}

type Testsuite struct {
	SuiteName            string
	ToLowerSuiteBaseName string //path.Base(SuiteName)
	*Project
}

func (p *Project) Create() error {
	// check if AbsolutePath exists
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(p.AbsolutePath, 0754); err != nil {
			return err
		}
	}

	// create main.go
	mainFile, err := os.Create(fmt.Sprintf("%s/main.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer mainFile.Close()

	mainTemplate := template.Must(template.New("main").Parse(string(tpl.MainTemplate())))
	err = mainTemplate.Execute(mainFile, p)
	if err != nil {
		return err
	}

	// create execute/execute.go
	if _, err = os.Stat(fmt.Sprintf("%s/execute", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/execute", p.AbsolutePath), 0751))
	}
	executeFile, err := os.Create(fmt.Sprintf("%s/execute/execute.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer executeFile.Close()

	executeTemplate := template.Must(template.New("execute").Parse(string(tpl.ExecuteTemplate())))
	err = executeTemplate.Execute(executeFile, p)
	if err != nil {
		return err
	}
	// create testsuites/
	if _, err = os.Stat(fmt.Sprintf("%s/testsuites", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/testsuites", p.AbsolutePath), 0751))
	}

	return nil
}

func (c *Testsuite) Create() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	var suiteFile *os.File
	if strings.Contains(pwd, "testsuites") {
		dst := path.Join(pwd, c.SuiteName)
		if !strings.Contains(dst, "testsuites") {
			fmt.Printf("The absolute path is %s, please enter the correct relative path of testsuite", dst)
			os.Exit(1)
		}
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			cobra.CheckErr(os.MkdirAll(dst, 0751))
		}
		suiteFile, err = os.Create(path.Join(dst, c.ToLowerSuiteBaseName))
		if err != nil {
			return err
		}
	} else {
		if _, err := os.Stat(path.Join(c.AbsolutePath, "testsuites")); os.IsNotExist(err) {
			fmt.Println(err.Error() + "\nPlease add testsuites under the path of potato project")
			os.Exit(1)
		}
		dst := path.Join(c.AbsolutePath, "testsuites", c.SuiteName)
		if !strings.Contains(dst, "testsuites") {
			fmt.Printf("The absolute path is %s, please enter the correct relative path of testsuite", dst)
			os.Exit(1)
		}
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			cobra.CheckErr(os.MkdirAll(dst, 0751))
		}
		suiteFile, err = os.Create(path.Join(dst, fmt.Sprintf("%s.go", c.ToLowerSuiteBaseName)))
		if err != nil {
			return err
		}
	}
	defer suiteFile.Close()

	SuiteTemplate := template.Must(template.New("suite").Parse(string(tpl.AddSuiteTemplate())))
	err = SuiteTemplate.Execute(suiteFile, c)
	if err != nil {
		return err
	}
	return nil
}
