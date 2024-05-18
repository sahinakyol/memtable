package main

import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

const (
	LIMIT_SIZE = 64
)

// MEMTABLE_QUEUE level tables s0,s1,s3..Sn - It is simulating flush process
var MEMTABLE_QUEUE = make([]*MemTable, 0)

type MemTable struct {
	data      *SkipList
	limitSize int
	mu        sync.RWMutex
}

type MemTableValue struct {
	value     interface{}
	timestamp int64
	deleted   bool
}

func NewMemTableValue(value interface{}) *MemTableValue {
	return &MemTableValue{
		value:     value,
		timestamp: time.Now().UnixNano(),
		deleted:   false,
	}
}

func (mtv *MemTableValue) MemTableValueSize() int {
	valueSize := len(mtv.value.(string))
	timestampSize := int(unsafe.Sizeof(mtv.timestamp))
	deletedSize := int(unsafe.Sizeof(mtv.deleted))
	totalSize := valueSize + timestampSize + deletedSize
	return totalSize
}

func NewMemTable() *MemTable {
	return &MemTable{
		data:      NewSkipList(),
		limitSize: LIMIT_SIZE,
	}
}

func (mt *MemTable) IsMemTableFull(key []byte, valueSize int) bool {
	requestedSize := len(key) + valueSize
	isFit := mt.limitSize-requestedSize >= 0
	if isFit {
		mt.limitSize -= requestedSize
		return false
	}
	return true
}

func (mt *MemTable) appendMemTableToQueue() *MemTable {
	MEMTABLE_QUEUE = append(MEMTABLE_QUEUE, mt)
	return NewMemTable()
}

func (mt *MemTable) Put(key string, value *MemTableValue) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	hashKey := int(Fnv1aHash(key))
	data := mt.data.Search(hashKey)

	if data == nil {
		if mt.IsMemTableFull([]byte(key), value.MemTableValueSize()) {
			mt = mt.appendMemTableToQueue()
		}
		mt.data.Insert(hashKey, value)
		return
	}

	mt.data.Insert(hashKey, value)
}

func (mt *MemTable) Get(key string) (interface{}, bool) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	hashKey := int(Fnv1aHash(key))
	data := mt.data.Search(hashKey)
	if data != nil {
		return data.value, true
	}

	for i := len(MEMTABLE_QUEUE) - 1; i >= 0; i-- {
		data = MEMTABLE_QUEUE[i].data.Search(hashKey)
		if data != nil {
			return data.value, true
		}
	}
	return nil, false
}

func (mt *MemTable) Delete(key string) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	hashKey := int(Fnv1aHash(key))
	data := mt.data.Search(hashKey)
	if data != nil {
		mt.data.Insert(hashKey, &MemTableValue{
			value:     data.value.(*MemTableValue).value,
			timestamp: time.Now().UnixNano(),
			deleted:   true,
		})
	}
}

func main() {
	memTable := NewMemTable()
	memTable.Put("name", NewMemTableValue("John"))
	memTable.Put("surname", NewMemTableValue("Doe"))
	memTable.Put("salute", NewMemTableValue("HELLO"))

	if name, exists := memTable.Get("name"); exists {
		fmt.Println("Name:", name)
	}
	if surname, exists := memTable.Get("surname"); exists {
		fmt.Println("Surname:", surname)
	}

	memTable.Delete("surname")

	if surname, exists := memTable.Get("surname"); exists {
		fmt.Println("Surname:", surname)
	}
}
