package memoryStorage

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/grahovsky/system-stats-daemon/internal/config"
	"github.com/grahovsky/system-stats-daemon/internal/logger"
)

type element struct {
	timestamp time.Time
	data      interface{}
}

type MemoryStorage struct {
	rwm  sync.RWMutex
	list *list.List
	size int64
}

func New() *MemoryStorage {
	return &MemoryStorage{rwm: sync.RWMutex{}, list: list.New(), size: config.Settings.Stats.Limit}
}

func (ms *MemoryStorage) SetSize(owner string, newsize int64) {
	ms.rwm.Lock()
	defer ms.rwm.Unlock()

	if newsize > ms.size {
		ms.size = newsize
		logger.Info(fmt.Sprintf("[%s] changed size of storage. New size: %d", owner, newsize))
	}
}

func (ms *MemoryStorage) Push(s interface{}, t time.Time) {
	ms.rwm.Lock()
	defer ms.rwm.Unlock()

	if ms.size == 0 {
		return
	}
	if ms.list.Len() == int(ms.size) {
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

func (ms *MemoryStorage) Show() {
	ms.rwm.RLock()
	defer ms.rwm.RUnlock()

	for e := ms.list.Front(); e != nil; e = e.Next() {
		fmt.Printf("%s: %+v\n", e.Value.(element).timestamp, e.Value.(element).data)
	}
}
