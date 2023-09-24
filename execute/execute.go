package execute

import (
	"log"
	"reflect"
)

func Execute(testsuite interface{}) {
	suiteType := reflect.TypeOf(testsuite)
	suiteValue := reflect.ValueOf(testsuite)
	funcElements := []reflect.Value{suiteValue}
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
			for i := 0; i < suiteType.NumMethod(); i++ {
				if suiteType.Method(i).Name == "Execute" || suiteType.Method(i).Name == "Setup" || suiteType.Method(i).Name == "Teardown" {
					continue
				}
				log.Printf("FAIL: %s.%s", suiteType.Elem().Name(), suiteType.Method(i).Name)
			}
			return
		}
	}
	// testcase
	for i := 0; i < suiteType.NumMethod(); i++ {
		if suiteType.Method(i).Name == "Execute" || suiteType.Method(i).Name == "Setup" || suiteType.Method(i).Name == "Teardown" {
			continue
		}
		method := suiteType.Method(i)
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
