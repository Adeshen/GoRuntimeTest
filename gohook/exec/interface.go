package exec

var TestVar Test = DefaultTest{}

type Test interface {
	Test()
	AfterMain()
	BeforeMain()
	NewProc(uintptr)
}

type DefaultTest struct{}

func (DefaultTest) Test() {
	print("default test\n")
}

func (DefaultTest) AfterMain() {
	print("default after main\n")
}

func (DefaultTest) BeforeMain() {
	print("default before main\n")
}

func (DefaultTest) NewProc(uintptr) {
	print("default new proc\n")
}
