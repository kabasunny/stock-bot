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

type ConcreteABC struct { // フィールドに各実装を持つだけで、ABCInterfaceを実装したわけではない
	A ImplementationA
	B ImplementationB
	C ImplementationC
} // フィールドは具象型で、単なるデータ構造体。ポリモーフィズムが活用できない（蜜結合で、柔軟性や拡張性が不要な特化用途）

func NewConcreteABC() *ConcreteABC {
	return &ConcreteABC{
		A: ImplementationA{},
		B: ImplementationB{},
		C: ImplementationC{},
	}
}
