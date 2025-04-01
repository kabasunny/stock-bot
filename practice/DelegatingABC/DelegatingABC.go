package DelegatingABC

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

type DelegatingABC struct {
	Abc InterfaceABC
}

func NewDelegatingABC(abc InterfaceABC) *DelegatingABC {
	return &DelegatingABC{Abc: abc}
}

func (d *DelegatingABC) MethodA() string {
	return d.Abc.MethodA()
}

func (d *DelegatingABC) MethodB() string {
	return d.Abc.MethodB()
}

func (d *DelegatingABC) MethodC() string {
	return d.Abc.MethodC()
}
