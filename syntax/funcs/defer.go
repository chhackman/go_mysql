package main

import "fmt"

func Defer() {
	defer func() {
		println("第一个defer")
	}()

	defer func() {
		println("第二个defer")
	}()
}

func DeferClosure() {
	i := 0
	fmt.Printf("第一个 i:%p\n", &i)
	defer func() {
		fmt.Printf("defer内 i:%p\n", &i)
		println(i)
	}()
	i = 1
	fmt.Printf("第二个 i:%p\n", &i)
}

//1

func DeferClosureV1() {
	i := 0
	fmt.Printf("第一个 i:%p\n", &i)
	defer func(i int) {
		println(i)
		fmt.Printf("defer i:%p\n", &i)
	}(i)
	i = 1
	fmt.Printf("第二个 i:%p\n", &i)
}

//i:0xc0000180a8
//第二个 i:0xc0000180a8
//0
//defer i:0xc0000180e0

//返回0
func DeferReturn() int {
	a := 0
	fmt.Printf("a地址: %p\n", &a)
	defer func() {
		a = 1
		fmt.Printf("defer地址: %p\n", &a)
	}()
	return a

}

//返回1
func DeferReturnV1() (a int) {
	a = 0
	defer func() {
		a = 1
	}()
	return a
}
