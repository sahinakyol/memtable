package main

import (
	"testing"
	"time"
)

func TestMemTablePutAndGet(t *testing.T) {
	memTable := NewMemTable()

	memTable.Put("name", NewMemTableValue("John"))
	memTable.Put("surname", NewMemTableValue("Doe"))

	if name, exists := memTable.Get("name"); !exists || name.(*MemTableValue).value != "John" {
		t.Errorf("Expected 'John', got '%v'", name.(*MemTableValue).value)
	}
	if surname, exists := memTable.Get("surname"); !exists || surname.(*MemTableValue).value != "Doe" {
		t.Errorf("Expected 'Doe', got '%v'", surname.(*MemTableValue).value)
	}
}

func TestMemTableDelete(t *testing.T) {
	memTable := NewMemTable()

	memTable.Put("name", NewMemTableValue("John"))
	memTable.Delete("name")

	if name, exists := memTable.Get("name"); exists != true || !name.(*MemTableValue).deleted {
		t.Errorf("Expected 'name' to be deleted, but got '%v'", name)
	}
}

func TestConcurrentAccess(t *testing.T) {
	memTable := NewMemTable()

	go memTable.Put("name", NewMemTableValue("Concurrent"))
	go memTable.Put("age", NewMemTableValue("30"))

	time.Sleep(100 * time.Millisecond) // Give some time for goroutines to complete

	if name, exists := memTable.Get("name"); exists && name.(*MemTableValue).value != "Concurrent" {
		t.Errorf("Expected 'Concurrent', got '%v'", name)
	}
	if age, exists := memTable.Get("age"); exists && age.(*MemTableValue).value != "30" {
		t.Errorf("Expected '30', got '%v'", age)
	}
}
