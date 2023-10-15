package utils

import (
	"bytes"
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
)

func GetTestSuiteName(suite interface{}) (name string) {
	typeofsuite := reflect.TypeOf(suite)
	name = typeofsuite.Elem().Name()
	return
}

func GetMethodsImplementedByUser(suite interface{}) (ret []reflect.Method) {
	var suiteType reflect.Type
	if v, ok := suite.(reflect.Type); ok {
		suiteType = v
	} else {
		suiteType = reflect.TypeOf(suite)
	}
	for i := 0; i < suiteType.NumMethod(); i++ {
		if suiteType.Method(i).Name == "Execute" || suiteType.Method(i).Name == "Setup" || suiteType.Method(i).Name == "Teardown" {
			continue
		}
		method := suiteType.Method(i)
		ret = append(ret, method)
	}
	return
}

func getAlldirs(dirPath string) (ret []string) {
	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil
	}
	if len(fs) == 0 {
		return
	}
	for _, f := range fs {
		if f.IsDir() {
			path := path.Join(dirPath, f.Name())
			ret = append(ret, path)
			tmp := getAlldirs(path)
			ret = append(ret, tmp...)
		}
	}
	return
}

// return the absolute dir of testsuite of the Potato Automated Testing Project
func GetTestSuiteAbsoluteRootDir(testsuitedirname string) (dst string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	if strings.Contains(wd, testsuitedirname) {
		dst = path.Join(strings.Split(wd, testsuitedirname)[0], testsuitedirname)
	} else {
		if _, err = os.Stat(path.Join(wd, testsuitedirname)); os.IsNotExist(err) {
			err = fmt.Errorf(err.Error() + "\nPlease execute comment under the rootpath of potato project or the path of testsuites")
			return
		}
		dst = path.Join(wd, testsuitedirname)
	}
	return
}

func getAllTestSuitesPath(testSuiteAbsoluteRootDir string) (ret []string) {
	return getAlldirs(testSuiteAbsoluteRootDir)
}

// format: map[Hello2Suite:description ...]
func GetAllTestSuitesComment(testSuiteAbsoluteRootDir string) (ret map[string]string, err error) {
	dirs := getAllTestSuitesPath(testSuiteAbsoluteRootDir)
	fset := token.NewFileSet()
	ret = make(map[string]string)
	for _, dir := range dirs {
		pkgs, err := parser.ParseDir(fset, dir, nil, 4)
		if err != nil {
			return nil, err
		}
		for k, v := range pkgs {
			targetPkg := doc.New(v, k, 4)
			for _, t := range targetPkg.Types {
				ret[t.Name] = t.Doc
			}
		}
	}
	return
}

func GetTestSuitesNum(m map[string]string) int {
	return len(m)
}

// format: map[Hello2Suite:map[Walk:description] ...]
func GetAllTestCasesComment(testSuiteAbsoluteRootDir string) (ret map[string]map[string]string, err error) {
	dirs := getAllTestSuitesPath(testSuiteAbsoluteRootDir)
	fset := token.NewFileSet()
	ret = make(map[string]map[string]string)
	for _, dir := range dirs {
		pkgs, err := parser.ParseDir(fset, dir, nil, 4)
		if err != nil {
			return nil, err
		}
		for k, v := range pkgs {
			targetPkg := doc.New(v, k, 4)
			for _, t := range targetPkg.Types {
				tmp := make(map[string]string, 0)
				for _, m := range t.Methods {
					if m.Name == "Execute" || m.Name == "Setup" || m.Name == "Teardown" {
						continue
					}
					tmp[m.Name] = m.Doc
				}
				ret[t.Name] = tmp
			}
		}
	}
	return
}

func GetTestCasesNum(m map[string]map[string]string) (ret int) {
	for _, v := range m {
		for range v {
			ret++
		}
	}
	return
}

// formatted output
func PrettySuiteComment(m map[string]string) string {
	var (
		ret   string
		max   int
		count int
	)
	for k := range m {
		if len(k) > max {
			max = len(k)
		}
	}
	for k, v := range m {
		space := string(bytes.Repeat([]byte(" "), max-len(k)))
		count++
		line := fmt.Sprintf("%d. %s%s: %s\n", count, k, space, v)
		ret += line
	}
	return ret
}

// formatted output
func PrettyCaseComment(m map[string]map[string]string) string {
	var (
		ret string
		max int
	)
	for _, v := range m {
		for k := range v {
			if len(k) > max {
				max = len(k)
			}
		}
	}
	for foo, bar := range m {
		ret += fmt.Sprintf("%s:\n", foo)
		count := 1
		for k, v := range bar {
			pre := fmt.Sprintf("%d. %s", count, k)
			space := string(bytes.Repeat([]byte(" "), max+5-len(pre)))
			line := fmt.Sprintf(" %d. %s%s: %s\n", count, k, space, v)
			count++
			ret += line
		}
	}
	return ret
}

// persistence
func CommentWrite(m interface{}, dst string, testsuiteText string, testcaseText string) {
	// testsuites
	testsuiteCommnetFile, err := os.OpenFile(path.Join(dst, "comment", "testSuiteCommentFile.go"), os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer testsuiteCommnetFile.Close()
	testsuiteCommentTemplate := template.Must(template.New("testsuiteComment").Parse(testsuiteText))
	err = testsuiteCommentTemplate.Execute(testsuiteCommnetFile, m)
	if err != nil {
		log.Panic(err)
	}
	// testcases
	testcaseCommnetFile, err := os.OpenFile(path.Join(dst, "comment", "testCaseCommentFile.go"), os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer testcaseCommnetFile.Close()
	testcaseCommentTemplate := template.Must(template.New("testcaseComment").Parse(testcaseText))
	err = testcaseCommentTemplate.Execute(testcaseCommnetFile, m)
	if err != nil {
		log.Panic(err)
	}
}
