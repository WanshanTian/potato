package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/txy2023/potato/tpl"
	"golang.org/x/mod/modfile"
)

// Project contains name and paths to projects.
type Project struct {
	ModName      string
	AbsolutePath string
	Viper        bool
	// ProjectName  string
}

type Testsuite struct {
	SuiteName     string //相对路径名
	SuiteBaseName string //path.Base(SuiteName)
	Dst           string
	PackageName   string //the packageName of testsuite
	StructName    string
	ImportedPath  string //execute包中导入该testsuite时的import路径

	*Project
}

var (
	testsuites = "testsuites"
)

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
	var (
		suiteFile         *os.File
		suiteInitFile     *os.File
		suiteRegisterFile *os.File
		dst               string
	)
	// acquire dst
	if strings.Contains(c.AbsolutePath, "testsuites") {
		dst = path.Join(c.AbsolutePath, c.SuiteName)
	} else {
		if _, err := os.Stat(path.Join(c.AbsolutePath, "testsuites")); os.IsNotExist(err) {
			fmt.Println(err.Error() + "\nPlease add a testsuite under the rootpath of potato project or the path of testsuites")
			os.Exit(1)
		}
		dst = path.Join(c.AbsolutePath, "testsuites", c.SuiteName)
	}
	if !strings.Contains(dst, "testsuites") {
		fmt.Printf("The absolute path is %s, please enter the correct relative path of testsuite", dst)
		os.Exit(1)
	}
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		cobra.CheckErr(os.MkdirAll(dst, 0751))
	}
	c.Dst = dst
	// acquire the module name
	modFile := path.Join(strings.Split(dst, testsuites)[0], "go.mod")
	if _, err := os.Stat(modFile); os.IsNotExist(err) {
		fmt.Println("please go mod init [mod-name] manually")
		os.Exit(1)
	}
	goModBytes, err := ioutil.ReadFile(path.Join(modFile))
	if err != nil {
		return err
	}
	c.ModName = modfile.ModulePath(goModBytes)
	c.ImportedPath = path.Join(c.ModName, testsuites, strings.Split(dst, testsuites)[1])
	// create suiteFile where you can add new testcase
	suiteFilePath := path.Join(dst, fmt.Sprintf("%s.go", c.SuiteBaseName))
	suiteFile, err = os.Create(suiteFilePath)
	if err != nil {
		return err
	}
	defer suiteFile.Close()

	SuiteTemplate := template.Must(template.New("suite").Parse(string(tpl.AddSuiteTemplate())))
	err = SuiteTemplate.Execute(suiteFile, c)
	if err != nil {
		return err
	}
	// create initFile which completes initialization work
	suiteInitFile, err = os.Create(path.Join(dst, "init.go"))
	if err != nil {
		return err
	}
	defer suiteInitFile.Close()
	SuiteInitTemplate := template.Must(template.New("init").Parse(string(tpl.SuiteInitTemplate())))
	err = SuiteInitTemplate.Execute(suiteInitFile, c)
	if err != nil {
		return err
	}
	// create suiteRegisterFile
	dir := path.Join(strings.Split(dst, "testsuites")[0], "execute")
	suiteRegisterFile, err = os.Create(path.Join(dir, fmt.Sprintf("%sSuiteRegiste.go", c.SuiteBaseName)))
	if err != nil {
		return err
	}
	defer suiteRegisterFile.Close()
	suiteRegisterFileTemplate := template.Must(template.New("registe").Parse(string(tpl.SuiteRegisteTemplate())))
	err = suiteRegisterFileTemplate.Execute(suiteRegisterFile, c)
	return err
}
