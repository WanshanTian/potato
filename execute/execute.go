package execute

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/txy2023/potato/register"
	"github.com/txy2023/potato/utils"
)

var (
	testCasesSpecified                  *string
	testSuitesSpecified                 *string
	TestSuitesExec                      []register.TestSuite
	TestSuitesExecForTestCasesSpecified []register.TestSuite
)

func IsExistTestcases() {
	if testCasesSpecified == nil {
		return
	}
	reg := regexp.MustCompile(`(\w+,?)+`)
	if !reg.MatchString(*testCasesSpecified) {
		log.Printf("testcases should be seperated by comma")
		os.Exit(1)
	}
	testcases := strings.Split(*testCasesSpecified, ",")
	allcases := make(map[string]register.TestSuite)
	for _, suite := range register.Registered {
		for _, s := range utils.GetMethodsImplementedByUser(suite) {
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
	if testSuitesSpecified == nil {
		return
	}
	reg := regexp.MustCompile(`(\w+,?)+`)
	if !reg.MatchString(*testSuitesSpecified) {
		log.Printf("testsuites should be seperated by comma")
		os.Exit(1)
	}
	testSuitesSpecifiedSlice := strings.Split(*testSuitesSpecified, ",")
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
	//global variable testCasesSpecified !=nil
	var testCasesExec = []reflect.Method{}
	if testCasesSpecified != nil {
		testcases := strings.Split(*testCasesSpecified, ",")
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
	log.Println(strings.Repeat("*", len(fmt.Sprintf("Setup of testsuite: %s is being executing", suiteType.Elem().Name()))))
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
		log.Printf("Setup of testsuite: %s is being executing", suiteType.Elem().Name())
		ret := method.Func.Call(funcElements)
		// if the return of Setup !=nil, the testcases will be assumed to be failed
		if ret[0].Interface() != nil {
			log.Printf("  FAIL, ErrMsg: %s", ret[0].Interface())
			for _, method := range utils.GetMethodsImplementedByUser(testsuite.(register.TestSuite)) {
				log.Printf("Testcase: %s.%s is skipped", suiteType.Elem().Name(), method.Name)
				log.Printf("  FAIL, ErrMsg: %s", "Setup FAIL")
				fmt.Println()
			}
			return
		} else {
			log.Printf("  DOWN")
			fmt.Println()
		}
	}
	// testcase
	if len(testCasesExec) == 0 {
		testCasesExec = utils.GetMethodsImplementedByUser(testsuite.(register.TestSuite))
	}
	for _, method := range testCasesExec {
		if method.Type.NumIn() > 1 {
			log.Printf("  FAIL, ErrMsg: the numIn of %s should be equal 1", method.Name)
			continue
		}
		if method.Type.NumOut() > 1 {
			log.Printf("  FAIL, ErrMsg: the numOut of %s should be equal 1)", method.Name)
			continue
		}
		if method.Type.Out(0).String() != "error" {
			log.Printf("  FAIL, ErrMsg: the typeOut of %s should be error)", method.Name)
			continue
		}
		log.Printf("%s.%s is being testing", suiteType.Elem().Name(), method.Name)
		timeStat := time.Now()
		ret := method.Func.Call(funcElements)
		timeElapse := time.Since(timeStat)
		if ret[0].Interface() == nil {
			log.Printf("  PASS(%s)", timeElapse)
		} else {
			log.Printf("  FAIL(%s), ErrMsg: %s", timeElapse, ret[0].Interface())
		}
		fmt.Println()
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
		log.Printf("Teardown of testsuite: %s is being executing", suiteType.Elem().Name())
		ret := method.Func.Call(funcElements)
		if ret[0].Interface() != nil {
			log.Printf("  FAIL, ErrMsg: %s", ret[0].Interface())
			fmt.Println()
			return
		} else {
			log.Printf("  DOWN")
			fmt.Println()
		}
	}
}

func Run(tc, ts *string) {
	if *tc != "" {
		testCasesSpecified = tc
		IsExistTestcases()
		for _, testsuite := range TestSuitesExecForTestCasesSpecified {
			testsuite.Execute()
		}
		testCasesSpecified = nil
	}
	if *ts != "" {
		testSuitesSpecified = ts
		IsExistTestSuites()
		for _, testsuite := range TestSuitesExec {
			testsuite.Execute()
		}
	}
	if *tc == "" && *ts == "" {
		for _, testsuite := range register.Registered {
			testsuite.Execute()
		}
	}
}
