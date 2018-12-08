package subscriber

import (
	"sync"

	"github.com/iloy/wormhole/protocol"
)

// Subscriber :
type Subscriber struct {
	sync.Mutex

	name   string
	stream protocol.Wormhole_SubscribeServer
	errCh  chan<- error
}

// New :
func New(name string, stream protocol.Wormhole_SubscribeServer, errCh chan<- error) *Subscriber {
	ret := &Subscriber{
		name:   name,
		stream: stream,
		errCh:  errCh,
	}

	return ret
}

// Name :
func (s *Subscriber) Name() string {
	return s.name
}

// Send :
func (s *Subscriber) Send(no uint64, start bool, end bool, data []byte) error {
	err := s.stream.Send(&protocol.SubscribeResponse{
		OK:                   true,
		Message:              "",
		UniqueSequenceNumber: no,
		Start:                start,
		End:                  end,
		Payload:              data,
	})

	return err
}

// Close :
func (s *Subscriber) Close(err error) {
	s.Lock()

	if s.stream != nil {
		if err != nil {
			err2 := s.stream.Send(&protocol.SubscribeResponse{
				OK:      false,
				Message: err.Error(),
			})
			if err2 != nil {
				// TODO?
			}
		}

		s.stream = nil
	}

	if s.errCh != nil {
		s.errCh <- err
		close(s.errCh)

		s.errCh = nil
	}

	s.Unlock()
}
