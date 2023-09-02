package main

func main() {
	//判断变量是否为函数
	//isFunction := func(v any) bool {
	//	//获取变量的类型信息
	//	t := reflect.TypeOf(v)
	//	return t.Kind() == reflect.Func
	//}
	//Func13("wangjie is good man")
	//Func14("wangjie is good man")
	//name, age := Func9()
	//println(name, age)

	//Defer()
	//DeferClosure()
	//DeferClosureV1()
	//println(DeferReturn())
	//println(DeferReturnV1())

	//fn := Closure("大明")
	//if isFunction(fn) {
	//	fmt.Printf("fn是一个函数\n")
	//} else {
	//	fmt.Printf("fn不是一个函数\n")
	//}
	//fmt.Printf("fn的类型%T\n", fn)
	////fn其实已经从Closure里面返回了
	////但是我fn还要用到"大明"
	//println(fn())

	//getAge := Closure1()
	//println("age 是\n", getAge())
	//println("age 是\n", getAge())
	//println("age 是\n", getAge())
	//println("age 是\n", getAge())

	//DeferClosureLoopV1()

	//DeferClosureLoopV2()
	DeferClosureLoopV3()
}
