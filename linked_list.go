package main

import (
	"container/list"
	"sync"
)

type listMutex struct {
	*list.List
	sync.Mutex
}

func NewListMutex() *listMutex {
	return &listMutex{List: list.New()}
}
