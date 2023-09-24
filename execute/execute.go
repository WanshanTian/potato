package execute

import (
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/txy2023/potato/register"
)

var (
	TestCasesSpecified                  *string
	TestSuitesSpecified                 *string
	TestSuitesExec                      []register.TestSuite
	TestSuitesExecForTestCasesSpecified []register.TestSuite
)

func GetMethodsImplementedByUser(suite register.TestSuite) (ret []reflect.Method) {
	suiteType := reflect.TypeOf(suite)
	for i := 0; i < suiteType.NumMethod(); i++ {
		if suiteType.Method(i).Name == "Execute" || suiteType.Method(i).Name == "Setup" || suiteType.Method(i).Name == "Teardown" {
			continue
		}
		method := suiteType.Method(i)
		ret = append(ret, method)
	}
	return
}

func IsExistTestcases() {
	if *TestCasesSpecified == "" {
		return
	}
	reg := regexp.MustCompile(`(\w+,?)+`)
	if !reg.MatchString(*TestCasesSpecified) {
		log.Printf("testcases should be seperated by comma")
		os.Exit(1)
	}
	testcases := strings.Split(*TestCasesSpecified, ",")
	allcases := make(map[string]register.TestSuite)
	for _, suite := range register.Registered {
		for _, s := range GetMethodsImplementedByUser(suite) {
			allcases[strings.ToLower(s.Name)] = suite
		}
	}
	for _, c := range testcases {
		v, ok := allcases[c]
		if !ok {
			log.Printf("testcase: %s is not implemented\n", c)
			os.Exit(1)
		} else {
			flag := 0
			for _, s := range TestSuitesExecForTestCasesSpecified {
				if v == s {
					flag = 1
				}
			}
			if flag == 0 {
				TestSuitesExecForTestCasesSpecified = append(TestSuitesExecForTestCasesSpecified, v)
			}
		}
	}
}

func IsExistTestSuites() {
	if *TestSuitesSpecified == "" {
		return
	}
	reg := regexp.MustCompile(`(\w+,?)+`)
	if !reg.MatchString(*TestSuitesSpecified) {
		log.Printf("testsuites should be seperated by comma")
		os.Exit(1)
	}
	testSuitesSpecifiedSlice := strings.Split(*TestSuitesSpecified, ",")
	tmp := make(map[string]register.TestSuite)
	for _, testSuiteRegister := range register.Registered {
		tmp[strings.ToLower(reflect.TypeOf(testSuiteRegister).Elem().Name())] = testSuiteRegister
	}
	for _, testSuiteSlice := range testSuitesSpecifiedSlice {
		if _, ok := tmp[strings.ToLower(testSuiteSlice)]; ok {
			TestSuitesExec = append(TestSuitesExec, tmp[strings.ToLower(testSuiteSlice)])
		} else {
			log.Printf("testSuite: %s is not registered\n", testSuiteSlice)
			os.Exit(1)
		}
	}
}

func Execute(testsuite interface{}) {
	suiteType := reflect.TypeOf(testsuite)
	suiteValue := reflect.ValueOf(testsuite)
	funcElements := []reflect.Value{suiteValue}
	//global variable TestCasesSpecified !=nil
	var testCasesExec = []reflect.Method{}
	if *TestCasesSpecified == "" {
		testcases := strings.Split(*TestCasesSpecified, ",")
		tmp := make(map[string]reflect.Method)
		for i := 0; i < suiteType.NumMethod(); i++ {
			tmp[strings.ToLower(suiteType.Method(i).Name)] = suiteType.Method(i)
		}
		for _, testcase := range testcases {
			if _, ok := tmp[strings.ToLower(testcase)]; ok {
				testCasesExec = append(testCasesExec, tmp[strings.ToLower(strings.ToLower(testcase))])
			}
		}
		if len(testCasesExec) == 0 {
			return
		}
	}
	// setup
	if method, ok := suiteType.MethodByName("Setup"); ok {
		if method.Type.NumIn() > 1 {
			log.Printf("testsuite %s Setup fail (the numIn of %s should be equal 1)", suiteType.Elem().Name(), method.Name)
			return
		}
		if method.Type.NumOut() > 1 {
			log.Printf("testsuite %s Setup fail (the numOut of %s should be equal 1)", suiteType.Elem().Name(), method.Name)
			return
		}
		if method.Type.Out(0).String() != "error" {
			log.Printf("testsuite %s Setup fail (the typeOut of %s should be error)", suiteType.Elem().Name(), method.Name)
			return
		}
		ret := method.Func.Call(funcElements)
		// if the return of Setup !=nil, the testcases will be assumed to be failed
		if ret[0].Interface() != nil {
			log.Printf("testsuite %s Setup fail", suiteType.Elem().Name())
			for _, method := range GetMethodsImplementedByUser(testsuite.(register.TestSuite)) {
				log.Printf("FAIL: %s.%s", suiteType.Elem().Name(), method.Name)
			}
			return
		}
	}
	// testcase
	if len(testCasesExec) == 0 {
		testCasesExec = GetMethodsImplementedByUser(testsuite.(register.TestSuite))
	}
	for _, method := range testCasesExec {
		if method.Type.NumIn() > 1 {
			log.Printf("FAIL: %s.%s(the numIn of %s should be equal 1)", suiteType.Elem().Name(), method.Name, method.Name)
			continue
		}
		if method.Type.NumOut() > 1 {
			log.Printf("FAIL: %s.%s(the numOut of %s should be equal 1)", suiteType.Elem().Name(), method.Name, method.Name)
			continue
		}
		if method.Type.Out(0).String() != "error" {
			log.Printf("FAIL: %s.%s(the typeOut of %s should be error)", suiteType.Elem().Name(), method.Name, method.Name)
			continue
		}
		ret := method.Func.Call(funcElements)
		if ret[0].Interface() == nil {
			log.Printf("PASS: %s.%s", suiteType.Elem().Name(), method.Name)
		} else {
			log.Printf("FAIL: %s.%s", suiteType.Elem().Name(), method.Name)
		}
	}
	// teardown
	if method, ok := suiteType.MethodByName("Teardown"); ok {
		if method.Type.NumIn() > 1 {
			log.Printf("testsuite %s Teardown fail (the numIn of %s should be equal 1)", suiteType.Elem().Name(), method.Name)
			return
		}
		if method.Type.NumOut() > 1 {
			log.Printf("testsuite %s Teardown fail (the numOut of %s should be equal 1)", suiteType.Elem().Name(), method.Name)
			return
		}
		if method.Type.Out(0).String() != "error" {
			log.Printf("testsuite %s Teardown fail (the typeOut of %s should be error)", suiteType.Elem().Name(), method.Name)
			return
		}
		ret := method.Func.Call(funcElements)
		if ret[0].Interface() != nil {
			log.Printf("testsuite %s Teardown fail", suiteType.Elem().Name())
			return
		}
	}
}
