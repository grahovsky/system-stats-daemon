package storage

import "time"

type Storage interface {
	Push(interface{}, time.Time)
	GetElementsAt(time.Time) <-chan interface{}
	GetElements(int64) <-chan interface{}
	StoreAt() <-chan interface{}
	Show()
}
