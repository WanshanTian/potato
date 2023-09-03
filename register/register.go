package register

type TestSuite interface {
	Execute()
}

var Registered []TestSuite

func Registe(t TestSuite) []TestSuite {
	Registered = append(Registered, t)
	return Registered
}
