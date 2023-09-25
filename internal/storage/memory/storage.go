package memoryStorage

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/grahovsky/system-stats-daemon/internal/logger"
)

type element struct {
	timestamp time.Time
	data      interface{}
}

type MemoryStorage struct {
	rwm  sync.RWMutex
	list *list.List
	size int
}

func New(size int) *MemoryStorage {
	return &MemoryStorage{rwm: sync.RWMutex{}, list: list.New(), size: size}
}

func (ms *MemoryStorage) SetSize(owner string, newsize int32) {
	ms.rwm.Lock()
	defer ms.rwm.Unlock()

	if int(newsize) > ms.size {
		ms.size = int(newsize)
		logger.Info(fmt.Sprintf("[%s] changed size of storage. New size: %d", owner, newsize))
	}
}

func (ms *MemoryStorage) Push(s interface{}, t time.Time) {
	ms.rwm.Lock()
	defer ms.rwm.Unlock()

	if ms.size == 0 {
		return
	}
	if ms.list.Len() == ms.size {
		ms.list.Remove(ms.list.Back())
	}
	ms.list.PushFront(element{timestamp: t.Truncate(time.Second), data: s})
}

func (ms *MemoryStorage) GetElementsAt(t time.Time) <-chan interface{} {
	elemCh := make(chan interface{})
	go func() {
		ms.rwm.RLock()
		defer close(elemCh)
		defer ms.rwm.RUnlock()
		for last := ms.list.Front(); last != nil; last = last.Next() {
			elem := last.Value.(element)
			if t.After(elem.timestamp) {
				return
			}
			elemCh <- elem.data
		}
	}()

	return elemCh
}

func (l *MemoryStorage) Show() {
	l.rwm.RLock()
	defer l.rwm.RUnlock()

	for e := l.list.Front(); e != nil; e = e.Next() {
		fmt.Printf("%s: %+v\n", e.Value.(element).timestamp, e.Value.(element).data)
	}
}
