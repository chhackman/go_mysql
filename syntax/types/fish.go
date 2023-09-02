package main

type Integer int

func UseInt() {
	i1 := 10
	i2 := Integer(i1)
	var i3 Integer = 11
	println(i2, i3)
}

type Fish struct {
	Name string
}

func (f Fish) swim() {
	println("fish在游")
}

type FakeFish Fish

func UseFish() {
	f1 := Fish{}
	f2 := FakeFish(f1)
	f2.Name = "Tom"
	println(f1.Name)
	//println(f2.Name)

	var y Yu
	y.Name = "yu"
	
}

type Yu = Fish
