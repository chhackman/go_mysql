package main

//T类型参数，名字叫做T,约束是any, 等于没有约束
type List[T any] interface {
	Add(idx int, t T)
	Append(t T)
}

func main() {
	println(Sum[int](1, 2, 3))

}

type MyMarshal struct {
}

func (m *MyMarshal) MarshalJSON() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func UseList() {
	//1.Append(12)
	var lany List[any]
	lany.Append(12.3)
	lany.Append(123)
	lk := Lin
}

//type parameter
type LinkedList[Daming any] struct {
	head *node[Daming]
	t    Daming
	tp   *Daming
	tp2  ***Daming
}

type node[T any] struct {
	val T
}
