package seqno

import "sync/atomic"

// Seqno :
type Seqno struct {
	number uint64
}

// Get :
func (s *Seqno) Get() uint64 {
	ret := atomic.AddUint64(&s.number, 1)

	return ret - 1
}
