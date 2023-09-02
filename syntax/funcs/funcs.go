package main

import "strings"

//Func1没有任何参数
func Func1() {

}

//Func2有一个参数
func Func2(a int) {

}

//Func3多个参数
func Func3(a int, b string) {

}

//Func4多个参数，一个类型
func Func4(a, b string) {

}

func Func5(a string, b string) {

}

func Func6(string) string {
	//有返回值要保证一定返回
	return "hello,world"
}

//Func7 多个返回值
func Func7(a, b string) (string, string) {
	//有返回值要保证一定返回
	return "Deng", "Ming"
}

func Func8() (name string, age int) {
	return "Daming", 18
}

func Func9() (name string, age int) {
	name = "Daming"
	age = 18
	return
}
func Func10() (name string, age int) {
	//等价于"",0
	//对应类型的零值
	return
}

func Func11() (string, int) {
	var name string
	var age int
	return name, age
}

//要么都给名字，要么都没有名字
//func Func12()(name string,int)  {
//	return
//}
//

func Func13(abc string) (string, int) {
	segs := strings.Split(abc, " ")
	println(segs[0], len(segs))
	return segs[0], len(segs)
}

func Func14(abc string) (first string, length int) {
	//从这里开始生效
	segs := strings.Split(abc, " ")
	first = segs[0]
	length = len(segs)
	println(first, length)
	return
}
