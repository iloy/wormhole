package main

import (
	"errors"
	"io"
	"log"
	"time"

	"github.com/iloy/wormhole/lib"
	"github.com/iloy/wormhole/protocol"
)

const (
	clientUniqueName = "example/publisher/go"
)

var (
	i = uint64(0)
	j = uint64(0)
)

func callbackPublisher() (uint64, bool, bool, []byte, error) {
	var no uint64
	var start bool
	var end bool
	var data []byte
	var err error

	no = i

	if j == 0 {
		start = true
	}

	if j == 9 {
		end = true
	}

	data = []byte{byte('a') + byte(j)}

	if i == 9 && j == 9 {
		err = io.EOF
	}

	j = j + 1

	if j == 9+1 {
		i = i + 1
		j = 0
	}

	time.Sleep(time.Millisecond * 100)

	return no, start, end, data, err
}

func callbackMsgHandle(msg *protocol.PublishResponse) error {
	log.Println("callbackMsgHandle():", "no:", msg.UniqueSequenceNumber, "start:", msg.Start, "end:", msg.End, "data:", msg.Payload)

	if !msg.OK {
		return errors.New(msg.Message)
	}

	if msg.UniqueSequenceNumber == 9 {
		return io.EOF
	}

	return nil
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

	if false {
		err = lib.CreateTopic(client, 4, token, topicName)
		if err != nil {
			panic(err)
		}
		log.Println("CreateTopic():", topicName, "by", id)
	}

	if false {
		err = lib.DestroyTopic(client, 4, token, topicName)
		if err != nil {
			panic(err)
		}
		log.Println("DestroyTopic():", topicName, "by", id)
	}

	if true {
		log.Println("Publish() start:", topicName, "by", id)
		err = lib.Publish(client, 0, token, topicName, callbackPublisher, callbackMsgHandle)
		if err != nil {
			panic(err)
		}
		log.Println("Publish() end:", topicName, "by", id)
	}
}
