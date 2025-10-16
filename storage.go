package storage

import "sync"

type Orders struct {
	OrderStorage map[string]string
	Mu           sync.Mutex
}