package main

import "container/list"

type Node struct {
	value    string
	children *list.List
}

func NewNode(nodeValue string) Node {
	return Node{nodeValue, list.New()}
}
