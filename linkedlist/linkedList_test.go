package linkedlist

import (
	"testing"
	"testing/quick"
)

func TestLinkedListPropertiesQuick(t *testing.T) {
	err := quick.Check(func(inputs []int) bool {
		l := NewLinkedList()

		for k, v := range inputs {
			ok := l.Insert(uint(k), v)
			if !ok {
				return false
			}
		}

		if l.length != uint(len(inputs)) {
			return false
		}

		for k, v := range inputs {
			out, ok := l.Get(uint(k))
			if !ok {
				return false
			}
			if out != v {
				return false
			}
		}

		for k, v := range inputs {
			index, found := l.Find(v)
			if !found {
				return false
			}
			if index != uint(k) {
				return false
			}
		}

		for v := uint(len(inputs)); v > 0; v-- {
			ok := l.Remove(v - 1)
			if !ok {
				return false
			}
		}

		if l.length != 0 {
			return false
		}

		for k := range inputs {
			_, ok := l.Get(uint(k))
			if ok {
				return false
			}
		}

		return true
	}, nil) // used default config

	if err != nil {
		t.Fatal(err)
	}
}

func TestLinkedListInsert(t *testing.T) {
	l := NewLinkedList()

	if !l.Insert(0, 10) {
		t.Error("Insert failed to insert the first element")
	}
	if l.head.Value != 10 || l.length != 1 {
		t.Errorf("Expected head value of 10 and length 1, got value %d and length %d", l.head.Value, l.length)
	}

	if !l.Insert(1, 20) {
		t.Error("Insert failed to append second element")
	}
	if l.head.Next.Value != 20 || l.length != 2 {
		t.Errorf("Expected second value of 20 and length 2, got value %d and length %d", l.head.Next.Value, l.length)
	}

	if !l.Insert(1, 15) {
		t.Error("Insert failed to insert element in the middle")
	}
	if l.head.Next.Value != 15 || l.head.Next.Next.Value != 20 || l.length != 3 {
		t.Errorf("Expected middle value of 15, got %d, expected third value of 20, got %d, expected length 3, got %d", l.head.Next.Value, l.head.Next.Next.Value, l.length)
	}

	if l.Insert(5, 30) {
		t.Error("Insert did not fail when trying to insert out of bounds")
	}
}

func TestLinkedListRemove(t *testing.T) {
	l := NewLinkedList()
	l.Insert(0, 10)
	l.Insert(1, 20)
	l.Insert(2, 30)

	if !l.Remove(0) {
		t.Error("Remove failed to remove the first element")
	}
	if l.head.Value != 20 || l.length != 2 {
		t.Errorf("Expected new head value of 20 and length 2, got value %d and length %d", l.head.Value, l.length)
	}

	if !l.Remove(1) {
		t.Error("Remove failed to remove the middle element")
	}
	if l.head.Next != nil || l.length != 1 {
		t.Errorf("Expected final element to be nil and length 1, got next value %v and length %d", l.head.Next, l.length)
	}

	if !l.Remove(0) {
		t.Error("Remove failed to remove the last element")
	}
	if l.head != nil || l.length != 0 {
		t.Errorf("Expected empty list, got head %v and length %d", l.head, l.length)
	}

	if l.Remove(0) {
		t.Error("Remove did not fail when trying to remove from an empty list")
	}
}

func TestLinkedListGet(t *testing.T) {
	l := NewLinkedList()
	l.Insert(0, 10)
	l.Insert(1, 20)
	l.Insert(2, 30)

	tests := []struct {
		index    uint
		expected int
		ok       bool
	}{
		{0, 10, true},
		{1, 20, true},
		{2, 30, true},
		{3, 0, false},
	}
	for _, tt := range tests {
		value, ok := l.Get(tt.index)
		if ok != tt.ok || (ok && value != tt.expected) {
			t.Errorf("Get(%d): expected %d, ok %t, got %d, ok %t", tt.index, tt.expected, tt.ok, value, ok)
		}
	}
}

func TestLinkedListFind(t *testing.T) {
	l := NewLinkedList()
	l.Insert(0, 10)
	l.Insert(1, 20)
	l.Insert(2, 30)

	tests := []struct {
		val      int
		expected uint
		found    bool
	}{
		{10, 0, true},
		{20, 1, true},
		{30, 2, true},
		{40, 0, false},
	}
	for _, tt := range tests {
		index, found := l.Find(tt.val)
		if found != tt.found || (found && index != tt.expected) {
			t.Errorf("Find(%d): expected index %d, found %t, got index %d, found %t", tt.val, tt.expected, tt.found, index, found)
		}
	}
}
