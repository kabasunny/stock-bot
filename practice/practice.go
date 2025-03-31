package main

import "fmt"

//--------第1層---------
type InterfaceABC interface {
	InterfaceA // 埋め込みフィールド(匿名フィールド)
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

type ComposedABC struct { // 個別にフィールドを持ち、InterfaceABC を実装する構造体
	a InterfaceA
	b InterfaceB
	c InterfaceC
} // 複数のインターフェースを組み合わせて新しい機能を実現する合成パターン

func NewComposedABC(a InterfaceA, b InterfaceB, c InterfaceC) *ComposedABC {
	return &ComposedABC{a: a, b: b, c: c}
}

func (i *ComposedABC) MethodA() string {
	return i.a.MethodA()
}

func (i *ComposedABC) MethodB() string {
	return i.b.MethodB()
}

func (i *ComposedABC) MethodC() string {
	return i.c.MethodC()
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

//--------代2層---------
type InterfaceDEF interface {
	InterfaceD
	InterfaceE
	InterfaceF
}

type InterfaceD interface {
	MethodD() string
}
type InterfaceE interface {
	MethodE() string
}
type InterfaceF interface {
	MethodF() string
}
type ImplementationD struct {
	A InterfaceA
}

func (i *ImplementationD) MethodD() string {
	return "MethodD" + A.MethodA() // i.ImplementationA.MethodA() から i.MethodA() に変更
}

type ImplementationE struct {
	b InterfaceB
}

type ImplementationF struct {
	//InterfaceF への依存
}

type ImplementationDEF struct {
	ImplementationD // InterfaceDに依存 (フィールド名を削除)
	ImplementationE
	ImplementationF
}

func NewImplementationDEF(abc *ABCImplementations) *ImplementationDEF {
	// d := ImplementationD{InterfaceA: &abc.ImplementationA}
	d := ImplementationD{InterfaceA: a}
	e := ImplementationE{}
	f := ImplementationF{}
	return &ImplementationDEF{
		ImplementationD: d,
		ImplementationE: e,
		ImplementationF: f,
	}
}

// --- メイン関数 ---

func main() {
	// ImplementationABCのインスタンスを生成
	abc := NewImplementationABC()

	// ImplementationDEFのインスタンスを生成 (InterfaceAの実装を注入)
	def := NewImplementationDEF(abc)

	// メソッドの呼び出し
	fmt.Println(def.ImplementationD.MethodD())
}
