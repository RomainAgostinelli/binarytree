package binarytree

import "testing"

type UsedItem struct {
	val string
}

func TestBTree(t *testing.T) {
	tree := BinaryTree{}
	if tree.IsEmpty() != true {
		t.Fatal("A newly created BinaryTree must be empty")
	}

	itr := tree.Root()
	itr.Insert(&UsedItem{val: "test"})
	if itr.IsRoot() == itr.IsLeaf() == false {
		t.Fatal("An iterator on a tree with one element is root and leaf")
	}

	left, _ := itr.Left()
	right, _ := itr.Right()
	left.Insert(&UsedItem{"Coucou"})
	right.Insert(&UsedItem{"Tests"})
	if itr.HasRight() == itr.HasLeft() == false {
		t.Fatal("Itr must have right and left child")
	}
}
