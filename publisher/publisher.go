package publisher

import (
	"sync"

	"github.com/iloy/wormhole/protocol"
)

// Publisher :
type Publisher struct {
	sync.Mutex

	name   string
	stream protocol.Wormhole_PublishServer
	errCh  chan<- error
}

// New :
func New(name string, stream protocol.Wormhole_PublishServer, errCh chan<- error) *Publisher {
	ret := &Publisher{
		name:   name,
		stream: stream,
		errCh:  errCh,
	}

	return ret
}

// Name :
func (p *Publisher) Name() string {
	return p.name
}

// Send :
func (p *Publisher) Send(no uint64, start bool, end bool, data []byte) error {
	err := p.stream.Send(&protocol.PublishResponse{
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
func (p *Publisher) Close(err error) {
	p.Lock()

	if p.stream != nil {
		if err != nil {
			err2 := p.stream.Send(&protocol.PublishResponse{
				OK:      false,
				Message: err.Error(),
			})
			if err2 != nil {
				// TODO?
			}
		}

		p.stream = nil
	}

	if p.errCh != nil {
		p.errCh <- err
		close(p.errCh)

		p.errCh = nil
	}

	p.Unlock()
}
