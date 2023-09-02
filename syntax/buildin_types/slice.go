package main

import "fmt"

func Slice() {
	s1 := []int{9, 8, 7}
	fmt.Printf("s1:%v,len=%d,cap=%d\n", s1, len(s1), cap(s1))
	//s1:[9 8 7],len=3,cap=3

	s2 := make([]int, 3, 4)
	fmt.Printf("s2:%v, len=%d, cap=%d\n", s2, len(s2), cap(s2))
	//s2:[0 0 0], len=3, cap=4

	s3 := make([]int, 4) //如果没有说明容量的话，和长度保持一致
	fmt.Printf("s3:%v,len=%d,cap=%d\n", s3, len(s3), cap(s3))
	//s3:[0 0 0 0],len=4,cap=4

	s4 := make([]int, 0, 4)
	s4 = append(s4, 1)
	s4 = append(s4, 2)
	fmt.Printf("s4:%v,len=%d,cap=%d\n", s4, len(s4), cap(s4))
	//s4:[1 2],len=2,cap=4

	s4 = append(s4, 3)
	s4 = append(s4, 4)
	s4 = append(s4, 5)
	fmt.Printf("s4:%v,len=%d,cap=%d\n", s4, len(s4), cap(s4))
	//s4:[1 2 3 4 5],len=5,cap=8 超过cap会自动扩容
}

func SubSlice() {
	s1 := []int{2, 4, 6, 8, 10}
	s2 := s1[1:3]
	fmt.Printf("s2:%v,len=%d,cap=%d\n", s2, len(s2), cap(s2))
	//s2:[4 6],len=2,cap=4
	//容量就是start开始往后，包括原本s1的底层数组的元素个数

	s3 := s1[2:]
	fmt.Printf("s3=%v,len=%d,cap=%d\n", s3, len(s3), cap(s3))
	//s3=[6 8 10],len=3,cap=3

	s4 := s1[:3]
	fmt.Printf("s4=%v,len=%d,cap=%d\n", s4, len(s4), cap(s4))
	//s4=[2 4 6],len=3,cap=5
}

func ShareSlice() {
	s1 := []int{1, 2, 3, 4}
	s2 := s1[2:]
	s2 = append(s2, 111)
	fmt.Printf("s2:%v,len=%d,cap=%d\n", s2, len(s2), cap(s2))
	//fmt.Printf("s1地址:%p,s1:%v,len=%d,cap=%d\n", &s1, s1, len(s1), cap(s1))
	//s2:[3 4],len=2,cap=2

	s2[0] = 99
	fmt.Printf("s2:%v,len=%d,cap=%d\n", s2, len(s2), cap(s2))
	fmt.Printf("s1:%v,len=%d,cap=%d\n", s1, len(s1), cap(s1))
	//s2:[99 4],len=2,cap=2

	////s2 = append(s2, 199)
	////s2[0] = 188
	//fmt.Printf("s2:%v,len=%d,cap=%d\n", s2, len(s2), cap(s2))
	//fmt.Printf("s1:%v,len=%d,cap=%d\n", s1, len(s2), cap(s2))

}
