package main

import "fmt"

func YourName(name string, aliases ...string) {
	if len(aliases) > 0 {
		fmt.Printf("aliases : %+v", aliases)
	}
	println(name)
}

func main() {
	YourName("大明", "小明", "瘦明")
}
