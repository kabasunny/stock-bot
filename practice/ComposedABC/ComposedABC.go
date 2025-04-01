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
	a InterfaceA
	b InterfaceB
	c InterfaceC
}

func NewComposedABC(a InterfaceA, b InterfaceB, c InterfaceC) *ComposedABC {
	return &ComposedABC{a: a, b: b, c: c}
}

func (c *ComposedABC) MethodA() string {
	return c.a.MethodA()
}

func (c *ComposedABC) MethodB() string {
	return c.b.MethodB()
}

func (c *ComposedABC) MethodC() string {
	return c.c.MethodC()
}
