package main

import "time"

type LinkedList struct {
	head      *node
	tail      *node
	Len       int
	CreatTime time.Time
}

func (l LinkedList) Add(idx int, val any) error {
	//TODO implement me
	panic("implement me")
}

func (l LinkedList) Append(val any) {
	//TODO implement me
	panic("implement me")
}

func (l LinkedList) Delete(idx int) (any, error) {
	//TODO implement me
	panic("implement me")
}

type node struct {
	prev *node
	next *node
}
