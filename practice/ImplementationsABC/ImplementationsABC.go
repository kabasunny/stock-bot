package ImplementationsABC

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

type ImplementationsABC struct { // フィールドに各実装を持つだけで、ABCInterfaceを実装したわけではない
	A InterfaceA
	B InterfaceB
	C InterfaceC
} // フィールドはインターフェース型で、柔軟で拡張性のあるコンポーネント。ポリモーフィズムの活用

func NewImplementationsABC(a InterfaceA, b InterfaceB, c InterfaceC) *ImplementationsABC {
	return &ImplementationsABC{
		A: a,
		B: b,
		C: c,
	}
}
