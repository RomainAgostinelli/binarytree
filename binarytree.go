package binarytree

import (
	"errors"
)

// BinaryNode struct representing a node of the BinaryTree, can contain an Item
type BinaryNode struct {
	Left, Right, Parent *BinaryNode
	Item                Item
}

// BinaryTree the tree itself, containing only its root
type BinaryTree struct {
	root *BinaryNode
}

// Iterator struct used to navigate in a BinaryTree
type Iterator struct {
	whole     *BinaryTree
	down      *BinaryNode
	up        *BinaryNode
	isLeftArc bool
}

// Item interface representing what kind of Item can be stored
type Item interface {
}

// IsEmpty tells if the tree is empty of not
func (t *BinaryTree) IsEmpty() bool {
	return t.root == nil
}

// Root returns new iterator positioned on the root of this tree
func (t *BinaryTree) Root() *Iterator {
	return NewIterator(t)
}

// -----------------------------------------------------------------------------
// ITERATOR
// -----------------------------------------------------------------------------

// NewIterator returns a new Iterator positioned on the top of the BinaryTree t.
func NewIterator(t *BinaryTree) *Iterator {
	return &Iterator{
		whole:     t,
		down:      t.root,
		up:        nil,
		isLeftArc: false,
	}
}

// IsRoot tells if the Iterator i is on the root (above the root node)
func (i *Iterator) IsRoot() bool {
	return i.up == nil
}

// IsBottom tells if the Iterator i is on the bottom of the B-Tree (under a leaf node)
func (i *Iterator) IsBottom() bool {
	return i.down == nil
}

// HasLeft tells if the down node has a left branch
func (i *Iterator) HasLeft() bool {
	return !i.IsBottom() && i.down.Left != nil
}

// HasRight tells if the down node has a right branch
func (i *Iterator) HasRight() bool {
	return !i.IsBottom() && i.down.Right != nil
}

// IsLeaf tells if the iterator is on a leaf node
func (i *Iterator) IsLeaf() bool {
	return !i.IsBottom() && !i.HasRight() && !i.HasLeft()
}

// Consult returns the data of the node on which the iterator is (down node)
func (i *Iterator) Consult() Item {
	return i.down.Item
}

// Up moves the Iterator upward (creates a new one)
// PRE: !isRoot returns an error if isRoot
func (i *Iterator) Up() (*Iterator, error) {
	if i.IsRoot() {
		return nil, errors.New("cannot go up when root")
	}
	itr := NewIterator(i.whole)
	itr.down = i.up
	itr.up = i.up.Parent
	itr.isLeftArc = !itr.IsRoot() && itr.down == itr.up.Left
	return itr, nil
}

// Right moves the Iterator on the right branch of down node (creates a new one)
// PRE: !IsBottom
func (i *Iterator) Right() *Iterator {
	itr := NewIterator(i.whole)
	itr.down = i.down.Right
	itr.up = i.down
	itr.isLeftArc = false
	return itr
}

// Left moves the Iterator on the left branch of down node (creates a new one)
// PRE: !IsBottom
func (i *Iterator) Left() *Iterator {
	itr := NewIterator(i.whole)
	itr.down = i.down.Left
	itr.up = i.down
	itr.isLeftArc = true
	return itr
}

// RightMost moves the Iterator down as far as possible following the right branches (creates aa new one)
func (i *Iterator) RightMost() *Iterator {
	// Copy the current
	itr := i.Alias()

	for !itr.IsBottom() {
		// Do not use "Right" method as it creates a new one each time
		itr.up = itr.down
		itr.down = itr.down.Right
		itr.isLeftArc = false
	}
	return itr
}

// LeftMost moves the Iterator down as far as possible following the left branches (creates aa new one)
func (i *Iterator) LeftMost() *Iterator {
	// Copy the current
	itr := i.Alias()

	for !itr.IsBottom() {
		// Do not use "Left" method as it creates a new one each time
		itr.up = itr.down
		itr.down = itr.down.Left
		itr.isLeftArc = true
	}
	return itr
}

// Alias returns a new Iterator on the same position
func (i *Iterator) Alias() *Iterator {
	return &Iterator{
		whole:     i.whole,
		down:      i.down,
		up:        i.up,
		isLeftArc: i.isLeftArc,
	}
}

// IsInside tells if the iterator is in the tree given in parameter
func (i *Iterator) IsInside(t *BinaryTree) bool {
	return i.whole == t
}

// Update the data in the down node
func (i *Iterator) Update(item Item) error {
	if i.IsBottom() {
		return errors.New("cannot update when iterator pointing to bottom")
	}
	i.down.Item = item
	return nil
}

// Insert replaces the subtree with a single node containing the item
func (i *Iterator) Insert(item Item) {
	i.Cut()
	node := &BinaryNode{
		Left:   nil,
		Right:  nil,
		Parent: i.up,
		Item:   item,
	}
	if i.IsRoot() {
		i.whole.root = node
	} else if i.isLeftArc {
		i.up.Left = node
	} else {
		i.up.Right = node
	}
	i.down = node
}

// Cut the subtree and return it as new element
func (i *Iterator) Cut() *BinaryTree {
	tree := &BinaryTree{}
	if i.IsBottom() {
		return tree
	}
	tree.root = i.down
	tree.root.Parent = nil
	if i.IsRoot() {
		i.whole.root = nil
	} else if i.isLeftArc {
		i.up.Left = nil
	} else {
		i.up.Right = nil
	}
	i.down = nil
	return tree
}

// Paste replaces the current subtree with the one given in parameter
// PRE : !IsInside (error if so) POST: t is now empty
func (i *Iterator) Paste(t *BinaryTree) error {
	if i.IsInside(t) {
		return errors.New("cannot paste a tree in which the iterator is inside")
	}
	i.Cut()
	if t.IsEmpty() { // nothing to do in this case
		return nil
	}
	// Take the root as node
	n := t.root
	n.Parent = i.up
	// Three cases
	if i.IsRoot() {
		i.whole.root = n
	} else if i.isLeftArc {
		i.up.Left = n
	} else {
		i.up.Right = n
	}
	// update iterator
	i.down = n
	// empty the tree
	t.root = nil
	return nil
}

// RotateRight rotate the tree regarding this schematic:
//       \ /                    \ /
// i -->  |						 | <-- i
//     	 [y]					[x]
//       / \	RotateRight		/ \
//    [x]  /c\       --> 	  /a\  [y]
//    / \					  	   / \
//  /a\ /b\						 /b\ /c\
// Where "[]" is a node, "/\" is subtree and "\ /" is parent tree
// Reverse operation of RotateLeft
func (i *Iterator) RotateRight() {
	y := i.Cut().Root()
	b := y.Left().Right().Cut()
	x := y.Left().Cut().Root()
	_ = y.Left().Paste(b)
	_ = x.Right().Paste(y.whole)
	_ = i.Paste(x.whole)
}

// RotateLeft rotate the tree regarding this schematic:
//       \ /                         \ /
//        | <-- i              i -->  |
//       [x]                       	 [y]
//       / \      RotateLeft         / \
//     /a\  [y]      -->          [x]  /c\
//          / \                   / \
//        /b\ /c\               /a\ /b\
// Where "[]" is a node, "/\" is subtree and "\ /" is parent tree
// Reverse operation of RotateRight
func (i *Iterator) RotateLeft() {
	x := i.Cut().Root()
	b := x.Right().Left().Cut()
	y := x.Right().Cut().Root()
	_ = x.Right().Paste(b)
	_ = y.Left().Paste(x.whole)
	_ = i.Paste(y.whole)
}
