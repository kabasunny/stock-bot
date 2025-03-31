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

type DelegatingABC struct { // まとめてフィールドを持ち、InterfaceABC を実装する構造体
	abc InterfaceABC
} // 処理を別のオブジェクトに委譲する委譲パターン

func NewDelegatingABC(abc InterfaceABC) *DelegatingABC {
	return &DelegatingABC{abc: abc}
}

// InterfaceABC のメソッドを実装
func (i *DelegatingABC) MethodA() string {
	return i.abc.MethodA()
}

func (i *DelegatingABC) MethodB() string {
	return i.abc.MethodB()
}

func (i *DelegatingABC) MethodC() string {
	return i.abc.MethodC()
}
