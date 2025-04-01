package ConcreteABC

type InterfaceABC interface {
	InterfaceA
	InterfaceB
	InterfaceC
}
type InterfaceA interface {
	MethodA() string
}
type InterfaceB interface {
	MethodB() string
}
type InterfaceC interface {
	MethodC() string
}
type ImplementationA struct{}

func (i *ImplementationA) MethodA() string {
	return "MethodA"
}

type ImplementationB struct{}

func (i *ImplementationB) MethodB() string {
	return "MethodB"
}

type ImplementationC struct{}

func (i *ImplementationC) MethodC() string {
	return "MethodC"
}

type ConcreteABC struct {
	A ImplementationA
	B ImplementationB
	C ImplementationC
}

func NewConcreteABC() *ConcreteABC {
	return &ConcreteABC{
		A: ImplementationA{},
		B: ImplementationB{},
		C: ImplementationC{},
	}
}
