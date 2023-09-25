package utils

import (
	"fmt"
	"io/ioutil"
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

func GetAlldirs(dirPath string) (ret []string) {
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
			tmp := GetAlldirs(path)
			ret = append(ret, tmp...)
		}
	}
	return
}

func GetTestSuiteAbsoluteDir(testsuitedirname string) (dst string, err error) {
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

func GetAllTestSuitesPath(testSuiteAbsoluteDir string) (ret []string) {
	return GetAlldirs(testSuiteAbsoluteDir)
}
