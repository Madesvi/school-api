package linked

import (
	"fmt"
	"testing"
)

func TestLL(t *testing.T) {
	newLinkedList := LinkedList{0, nil}
	newLinkedList.AddAtHead(100)
	newLinkedList.Insert(1, 200)
	newLinkedList.AddAtTail(500)

	ourList := newLinkedList.DisplayAllList()
	fmt.Println("TLE LINKED LIST", ourList)

	newLinkedList.DeleteAtHead()
	ourList = newLinkedList.DisplayAllList()
	fmt.Println("TLE LINKED LIST after delete", ourList)
}
