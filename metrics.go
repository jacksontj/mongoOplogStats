package main

import (
	"fmt"
	"sync"
)

func NewMetrics() *Metrics {
	return &Metrics{
		Insert: &sync.Map{},
		Update: &sync.Map{},
		Delete: &sync.Map{},
	}
}

type Metrics struct {
	Insert *sync.Map
	Update *sync.Map
	Delete *sync.Map
}

func (m *Metrics) Print() {
	// print it out
	m.Insert.Range(func(k, v interface{}) bool {
		fmt.Println(k, "Insert", *(v.(*int64)))
		return true
	})

	m.Update.Range(func(k, v interface{}) bool {
		fmt.Println(k, "Update", *(v.(*int64)))
		return true
	})

	m.Delete.Range(func(k, v interface{}) bool {
		fmt.Println(k, "Delete", *(v.(*int64)))
		return true
	})
}
