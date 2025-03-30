package main

import "fmt"

//--------第1層---------
type InterfaceABC interface {
	InterfaceA // 埋め込みフィールド
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
	return "ImplementationA: MethodA"
}

type ImplementationB struct{}

func (i *ImplementationB) MethodB() string {
	return "ImplementationB: MethodB"
}

type ImplementationC struct{}

func (i *ImplementationC) MethodC() string {
	return "ImplementationC: MethodC"
}

type ImplementationABC struct { // フィールドに各実装を持つだけで、ABCInterfaceを実装したわけではない
	A ImplementationA
	B ImplementationB
	C ImplementationC
}

// NewImplementationABC は、ImplementationABCのファクトリメソッド
func NewImplementationABC() *ImplementationABC {
	return &ImplementationABC{
		A: ImplementationA{},
		B: ImplementationB{},
		C: ImplementationC{},
	}
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
	a InterfaceA
}

func (i *ImplementationD) MethodD() string {
	return "ImplementationD: MethodD, " + a.MethodA() // i.ImplementationA.MethodA() から i.MethodA() に変更
}

type ImplementationE struct {
	b InterfaceB
}

type ImplementationF struct {
	//InterfaceF への依存
}

// ImplementationDEF
type ImplementationDEF struct {
	ImplementationD // InterfaceDに依存 (フィールド名を削除)
	ImplementationE
	ImplementationF
}

// NewImplementationDEF は、ImplementationDEFのファクトリメソッド
func NewImplementationDEF(abc *ImplementationABC) *ImplementationDEF {
	// d := ImplementationD{InterfaceA: &abc.ImplementationA}
	d := ImplementationD{InterfaceA: abc}
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
