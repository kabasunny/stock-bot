package main

import "fmt"

// InterfaceA
type InterfaceA interface {
	MethodA() string
}

// InterfaceB (例)
type InterfaceB interface {
	MethodB() string
}

// InterfaceC (例)
type InterfaceC interface {
	MethodC() string
}

// InterfaceABC は、A, B, C を集約するマーカーインターフェース
type InterfaceABC interface {
	InterfaceA
	InterfaceB
	InterfaceC
}

// InterfaceD
type InterfaceD interface {
	MethodD() string
}

// InterfaceE (例)
type InterfaceE interface {
	MethodE() string
}

// InterfaceF (例)
type InterfaceF interface {
	MethodF() string
}

// InterfaceDEF は、D, E, F を集約するマーカーインターフェース
type InterfaceDEF interface {
	InterfaceD
	InterfaceE
	InterfaceF
}

// ImplementationA
type ImplementationA struct{}

func (i *ImplementationA) MethodA() string {
	return "ImplementationA: MethodA"
}

// ImplementationB
type ImplementationB struct{}

func (i *ImplementationB) MethodB() string {
	return "ImplementationB: MethodB"
}

// ImplementationC
type ImplementationC struct{}

func (i *ImplementationC) MethodC() string {
	return "ImplementationC: MethodC"
}

// ImplementationABC は、A, B, C の実装を埋め込むだけの構造体
type ImplementationABC struct {
	ImplementationA // InterfaceAの実装を埋め込む (匿名フィールド)
	ImplementationB // InterfaceBの実装を埋め込む (匿名フィールド)
	ImplementationC // InterfaceCの実装を埋め込む (匿名フィールド)
}

// NewImplementationABC は、ImplementationABCのファクトリメソッド
func NewImplementationABC() *ImplementationABC {
	return &ImplementationABC{
		ImplementationA: ImplementationA{},
		ImplementationB: ImplementationB{},
		ImplementationC: ImplementationC{},
	}
}

// ImplementationD
type ImplementationD struct {
	InterfaceA // InterfaceAに依存 (フィールド名を削除)
}

func (i *ImplementationD) MethodD() string {
	return "ImplementationD: MethodD, " + i.MethodA() // i.ImplementationA.MethodA() から i.MethodA() に変更
}

// ImplementationE
type ImplementationE struct {
	//InterfaceE への依存
}

//ImplementationF
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
