package linked

import "fmt"

type Node struct {
	data int
	next *Node
}

type LinkedList struct {
	length int
	head   *Node
}

// Add data in head

func (l *LinkedList) AddAtHead(data int) {
	temp1 := &Node{data, nil}

	// {100, *Node2 | 200, nil}
	if l.head == nil {
		l.head = temp1
	} else {
		temp2 := l.head
		l.head = temp1
		temp1.next = temp2
	}
	l.length += 1
}

func (l *LinkedList) AddAtTail(data int) {
	temp1 := &Node{data, nil}

	if l.head == nil {
		l.head = temp1
	} else {
		temp2 := l.head
		for temp2.next != nil {
			temp2 = temp2.next
		}
		temp2.next = temp1

	}
	l.length += 1
}

// {[0]100, *node2 | [1]200, *node3 | [2]300, nil }
// if n == 1 ---> AddAtHead
// if n == l.length - 1 ---> AddAtTail
// if n == 2 ---> { 100, *node2 | 200, *temp1 | 250, *node3 | 300, nil}

func (l *LinkedList) Insert(n, data int) {
	switch {
	case n == 0:
		l.AddAtHead(data)
	case n == l.length:
		l.AddAtTail(data)
	default:
		temp1 := &Node{data, nil}
		temp2 := l.head // 200, *node3
		for i := 0; i < n-1; i++ {
			temp2 = temp2.next
		}
		temp1.next = temp2.next
		temp2.next = temp1
	}
	l.length += 1
}

func (l *LinkedList) DeleteAtHead() {
	temp := l.head
	l.head = temp.next

	l.length -= 1
}

func (l *LinkedList) DisplayAllList() []int {
	if l.head == nil {
		fmt.Println("List is empty")
		return nil
	}
	allList := []int{}
	temp := l.head // *Node1 -> *nil
	for temp != nil {
		allList = append(allList, temp.data)
		temp = temp.next
	}
	// fmt.Println(allList)
	return allList
}

// func LinkedListTest() {
// 	ll := LinkedList{0, nil}
// 	ll.DisplayAllList()
// 	ll.AddAtHead(200)
// 	ll.AddAtHead(100)
// 	ll.AddAtTail(300)
// 	ll.DisplayAllList()
// }
