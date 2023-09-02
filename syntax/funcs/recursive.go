package main

//递归使用不当，就有可能stack overflow

func Recursive(n int) {
	if n > 10 {
		return
	}
	Recursive(n + 1)
}
