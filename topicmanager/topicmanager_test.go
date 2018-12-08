package topicmanager_test

import (
	"testing"

	"github.com/iloy/wormhole/topicmanager"
)

func TestCreateTopic(t *testing.T) {
	manager := topicmanager.New()
	creator := "go test"

	{
		t1, err := manager.AddTopic("t1", creator)
		if err != nil {
			t.Fatal(err)
		}
		if t1 == nil {
			t.Fatal()
		}
	}

	{
		t2, err := manager.AddTopic("t2", creator)
		if err != nil {
			t.Fatal(err)
		}
		if t2 == nil {
			t.Fatal()
		}
	}

	{
		t1, err := manager.AddTopic("t1", creator)
		if err == nil {
			t.Fatal(err)
		}
		if t1 != nil {
			t.Fatal()
		}
	}
}
