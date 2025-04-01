package ComposedABC

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

type ComposedABC struct {
	A InterfaceA
	B InterfaceB
	C InterfaceC
}

func NewComposedABC(a InterfaceA, b InterfaceB, c InterfaceC) *ComposedABC {
	return &ComposedABC{A: a, B: b, C: c}
}

func (c *ComposedABC) MethodA() string {
	return c.A.MethodA()
}

func (c *ComposedABC) MethodB() string {
	return c.B.MethodB()
}

func (c *ComposedABC) MethodC() string {
	return c.C.MethodC()
}
