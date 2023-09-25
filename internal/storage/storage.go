package storage

import "time"

type Storage interface {
	Push(interface{}, time.Time)
	GetElementsAt(time.Time) <-chan interface{}
	Show()
}
