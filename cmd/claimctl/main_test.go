package main

import "testing"

func TestCommandHelpers(t *testing.T) {
	if structRecord("b", "a", 2).Quantity != 2 {
		t.Fatal("record helper")
	}
}
