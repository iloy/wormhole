package topicmanager

import (
	"sync"

	"github.com/iloy/wormhole/topic"
)

// TopicManager :
type TopicManager struct {
	sync.Mutex

	topics map[string]*topic.Topic
}

// New :
func New() *TopicManager {
	manager := &TopicManager{
		topics: make(map[string]*topic.Topic),
	}

	return manager
}

// AddTopic :
func (m *TopicManager) AddTopic(name string, creator string) (*topic.Topic, error) {
	m.Lock()

	if _, ok := m.topics[name]; ok {
		m.Unlock()
		return nil, topic.ErrTopicAlreadyExists
	}

	t := topic.New(name, creator)
	m.topics[name] = t

	m.Unlock()

	return t, nil
}

// RemoveTopic :
func (m *TopicManager) RemoveTopic(name string, creator string) error {
	m.Lock()

	t, ok := m.topics[name]

	if !ok {
		m.Unlock()
		return topic.ErrTopicNotExist
	}

	if t.Creator() != creator {
		m.Unlock()
		return topic.ErrNotYourTopic
	}

	delete(m.topics, name)
	topic.Delete(t)

	m.Unlock()

	return nil
}

// AddPublisher :
func (m *TopicManager) AddPublisher(topicName string, publisher topic.Publisher) (*topic.Topic, error) {
	if topicName == "" {
		return nil, topic.ErrInvalidTopic
	}

	if publisher == nil {
		return nil, topic.ErrInvalidPublisher
	}

	m.Lock()

	t, ok := m.topics[topicName]

	if !ok {
		m.Unlock()
		return nil, topic.ErrTopicNotExist
	}

	m.Unlock()

	err := t.AddPublisher(publisher)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// RemovePublisher :
func (m *TopicManager) RemovePublisher(topicName string, publisher topic.Publisher) error {
	if topicName == "" {
		return topic.ErrInvalidTopic
	}

	if publisher == nil {
		return topic.ErrInvalidPublisher
	}

	m.Lock()

	t, ok := m.topics[topicName]

	if !ok {
		m.Unlock()
		return topic.ErrTopicNotExist
	}

	m.Unlock()

	return t.RemovePublisher(publisher)
}

// AddSubscriber :
func (m *TopicManager) AddSubscriber(topicName string, subscriber topic.Subscriber) (*topic.Topic, error) {
	if topicName == "" {
		return nil, topic.ErrInvalidTopic
	}

	if subscriber == nil {
		return nil, topic.ErrInvalidSubscriber
	}

	m.Lock()

	t, ok := m.topics[topicName]

	if !ok {
		m.Unlock()
		return nil, topic.ErrTopicNotExist
	}

	m.Unlock()

	err := t.AddSubscriber(subscriber)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// RemoveSubscriber :
func (m *TopicManager) RemoveSubscriber(topicName string, subscriber topic.Subscriber) error {
	if topicName == "" {
		return topic.ErrInvalidTopic
	}

	if subscriber == nil {
		return topic.ErrInvalidSubscriber
	}

	m.Lock()

	t, ok := m.topics[topicName]

	if !ok {
		m.Unlock()
		return topic.ErrTopicNotExist
	}

	m.Unlock()

	return t.RemoveSubscriber(subscriber)
}
