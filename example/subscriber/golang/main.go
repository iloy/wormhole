package main

import (
	"log"

	"github.com/iloy/wormhole/lib"
	"github.com/iloy/wormhole/protocol"
)

const (
	clientUniqueName = "example/subscriber/go"
)

var (
	token string
)

func dataRecvCallbackFunc(no uint64, start bool, end bool, data []byte) (*protocol.SubscribeRequest, error) {
	log.Println("no:", no, "start:", start, "end:", end, "data:", data)

	var ret *protocol.SubscribeRequest
	if end {
		ret = &protocol.SubscribeRequest{
			UniqueSequenceNumber: no,
			Start:                true,
			End:                  true,
			Payload:              []byte("no 9"),
		}
	}

	return ret, nil
}

func main() {
	client, conn := lib.NewClient("127.0.0.1:18554", 0)
	defer func() { _ = conn.Close() }()

	token, id, err := lib.Login(client, 4, clientUniqueName)
	if err != nil {
		panic(err)
	}
	log.Println("token:", token)
	log.Println("id:", id)

	topicName := "topic_1"

	if true {
		err := lib.DestroyTopic(client, 4, token, topicName)
		if err != nil {
			//panic(err)
		}
		log.Println("DestroyTopic():", topicName)
	}

	if true {
		err := lib.CreateTopic(client, 4, token, topicName)
		if err != nil {
			panic(err)
		}
		log.Println("CreateTopic():", topicName)
	}

	if true {
		log.Println("Subscribe():", topicName)
		err := lib.Subscribe(client, 0, token, topicName, dataRecvCallbackFunc)
		if err != nil {
			panic(err)
		}
	}
}
