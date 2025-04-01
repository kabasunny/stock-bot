package main

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
	return "MethodD" + " + " + i.A.MethodA()
}

type ImplementationE struct {
	B InterfaceB
}

func (i *ImplementationE) MethodE() string {
	return "MethodE" + " + " + i.B.MethodB()
}

type ImplementationF struct {
	C InterfaceC
}

func (i *ImplementationF) MethodF() string {
	return "MethodF" + " + " + i.C.MethodC()
}

type ImplementationDEF struct {
	ImplementationD
	ImplementationE
	ImplementationF
}

func NewImplementationDEF(abc *DelegatingABC) *ImplementationDEF {
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
	// abc := NewImplementationABC()

	// ImplementationDEFのインスタンスを生成 (InterfaceAの実装を注入)
	// def := NewImplementationDEF(abc)

	// メソッドの呼び出し
	// fmt.Println(def.ImplementationD.MethodD())
}
