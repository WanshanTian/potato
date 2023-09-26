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

	"github.com/txy2023/potato/tpl"
)

func GetTestSuiteName(suite interface{}) (name string) {
	typeofsuite := reflect.TypeOf(suite)
	name = typeofsuite.Elem().Name()
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

func CommentWrite(m interface{}, dst string) {
	testsuiteCommnetFile, err := os.OpenFile(path.Join(dst, "comment", "testSuiteCommentFile.go"), os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer testsuiteCommnetFile.Close()
	testsuiteCommentTemplate := template.Must(template.New("testsuiteComment").Parse(string(tpl.TestSuiteCommentTemplate())))
	err = testsuiteCommentTemplate.Execute(testsuiteCommnetFile, m)
	if err != nil {
		log.Panic(err)
	}
}
