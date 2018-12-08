package beancounter

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrAlreadyExists :
	ErrAlreadyExists = errors.New("already exists")
	// ErrNotExist :
	ErrNotExist = errors.New("not exist")
)

// Beancounter :
type Beancounter struct {
	sync.Mutex

	name string
	m    map[string]*time.Time
}

// New :
func New(name string) *Beancounter {
	ret := &Beancounter{
		name: name,
		m:    make(map[string]*time.Time),
	}

	return ret
}

// Add :
func (b *Beancounter) Add(name string) error {
	b.Lock()
	defer b.Unlock()

	if _, ok := b.m[name]; ok {
		return ErrAlreadyExists
	}

	now := time.Now()
	b.m[name] = &now

	return nil
}

// Remove :
func (b *Beancounter) Remove(name string) error {
	b.Lock()
	defer b.Unlock()

	if _, ok := b.m[name]; ok {
		delete(b.m, name)
		return nil
	}

	return ErrNotExist
}

// Rename :
func (b *Beancounter) Rename(fromname string, toname string) error {
	b.Lock()
	defer b.Unlock()

	if fromname == toname {
		return nil
	}

	if v, ok := b.m[fromname]; ok {
		delete(b.m, fromname)
		b.m[toname] = v
		return nil
	}

	return ErrNotExist
}
