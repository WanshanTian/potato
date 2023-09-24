package utils

import (
	"reflect"

	"github.com/txy2023/potato/register"
)

func GetTestSuiteName(suite register.TestSuite) (name string) {
	typeofsuite := reflect.TypeOf(suite)
	name = typeofsuite.Elem().Name()
	return
}
