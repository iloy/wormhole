package topic

import (
	"errors"
	"sync"
)

var (
	// ErrInvalidPublisher :
	ErrInvalidPublisher = errors.New("invalid publisher")
	// ErrPublisherAlreadyExists :
	ErrPublisherAlreadyExists = errors.New("publisher already exists")
	// ErrPublisherNotExist :
	ErrPublisherNotExist = errors.New("publisher not exist")

	// ErrInvalidSubscriber :
	ErrInvalidSubscriber = errors.New("invalid subscriber")
	// ErrSubscriberAlreadyExists :
	ErrSubscriberAlreadyExists = errors.New("subscriber already exists")
	// ErrSubscriberNotExist :
	ErrSubscriberNotExist = errors.New("subscriber not exist")

	// ErrInvalidTopic :
	ErrInvalidTopic = errors.New("invalid topic")
	// ErrTopicAlreadyExists :
	ErrTopicAlreadyExists = errors.New("topic already exists")
	// ErrTopicNotExist :
	ErrTopicNotExist = errors.New("topic not exist")

	// ErrNotYourTopic :
	ErrNotYourTopic = errors.New("not your topic")
)

// Topic :
type Topic struct {
	name    string
	creator string

	muPublisher sync.Mutex
	publisher   Publisher

	muSubscriber sync.Mutex
	subscriber   Subscriber
}

// Publisher :
type Publisher interface {
	Name() string
	Send(no uint64, start bool, end bool, data []byte) error
	Close(err error)
}

// Subscriber :
type Subscriber interface {
	Name() string
	Send(no uint64, start bool, end bool, data []byte) error
	Close(err error)
}

var (
	topicPool = sync.Pool{
		New: func() interface{} {
			return &Topic{}
		},
	}
)

// New :
func New(name string, creator string) *Topic {
	topic := topicPool.Get().(*Topic)

	topic.name = name
	topic.creator = creator

	return topic
}

// Delete :
func Delete(t *Topic) {
	t.muPublisher.Lock()
	t.muSubscriber.Lock()

	if t.publisher != nil {
		t.publisher.Close(nil)
		t.publisher = nil
	}

	if t.subscriber != nil {
		t.subscriber.Close(nil)
		t.subscriber = nil
	}

	t.muSubscriber.Unlock()
	t.muPublisher.Unlock()

	topicPool.Put(t)
}

// Creator :
func (t *Topic) Creator() string {
	return t.creator
}

// AddPublisher :
func (t *Topic) AddPublisher(publisher Publisher) error {
	t.muPublisher.Lock()

	if t.publisher != nil {
		t.muPublisher.Unlock()
		return ErrPublisherAlreadyExists
	}

	t.publisher = publisher

	t.muPublisher.Unlock()

	return nil
}

// RemovePublisher :
func (t *Topic) RemovePublisher(publisher Publisher) error {
	t.muPublisher.Lock()

	if t.publisher == nil {
		t.muPublisher.Unlock()
		return ErrPublisherNotExist
	}

	if t.publisher != publisher {
		t.muPublisher.Unlock()
		return ErrInvalidPublisher
	}

	t.publisher = nil

	t.muPublisher.Unlock()

	return nil
}

// AddSubscriber :
func (t *Topic) AddSubscriber(subscriber Subscriber) error {
	t.muSubscriber.Lock()

	if t.subscriber != nil {
		t.muSubscriber.Unlock()
		return ErrSubscriberAlreadyExists
	}

	t.subscriber = subscriber

	t.muSubscriber.Unlock()

	return nil
}

// RemoveSubscriber :
func (t *Topic) RemoveSubscriber(subscriber Subscriber) error {
	t.muSubscriber.Lock()

	if t.subscriber == nil {
		t.muSubscriber.Unlock()
		return ErrSubscriberNotExist
	}

	if t.subscriber != subscriber {
		t.muSubscriber.Unlock()
		return ErrInvalidSubscriber
	}

	t.subscriber = nil

	t.muSubscriber.Unlock()

	return nil
}

// SendToPublisher :
func (t *Topic) SendToPublisher(no uint64, start bool, end bool, data []byte) error {
	t.muPublisher.Lock()

	if t.publisher == nil {
		t.muPublisher.Unlock()
		return ErrPublisherNotExist
	}

	err := t.publisher.Send(no, start, end, data)
	if err != nil {
		name := t.publisher.Name()
		t.publisher.Close(err)
		t.publisher = nil

		err = errors.New("publisher " + name + " error: " + err.Error())

	}

	t.muPublisher.Unlock()

	return err
}

// SendToSubscriber :
func (t *Topic) SendToSubscriber(no uint64, start bool, end bool, data []byte) error {
	t.muSubscriber.Lock()

	if t.subscriber == nil {
		t.muSubscriber.Unlock()
		return ErrSubscriberNotExist
	}

	err := t.subscriber.Send(no, start, end, data)
	if err != nil {
		name := t.subscriber.Name()
		t.subscriber.Close(err)
		t.subscriber = nil

		err = errors.New("subscriber " + name + " error: " + err.Error())
	}

	t.muSubscriber.Unlock()

	return err
}
