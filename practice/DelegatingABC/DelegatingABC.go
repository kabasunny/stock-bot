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
	abc InterfaceABC
}

func NewDelegatingABC(abc InterfaceABC) *DelegatingABC {
	return &DelegatingABC{abc: abc}
}

func (d *DelegatingABC) MethodA() string {
	return d.abc.MethodA()
}

func (d *DelegatingABC) MethodB() string {
	return d.abc.MethodB()
}

func (d *DelegatingABC) MethodC() string {
	return d.abc.MethodC()
}
