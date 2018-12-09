package server

import (
	log "github.com/sirupsen/logrus"

	"github.com/iloy/wormhole/protocol"
	"github.com/iloy/wormhole/publisher"
	"github.com/iloy/wormhole/token"
)

// Publish :
func (rs *rpcServer) Publish(
	stream protocol.Wormhole_PublishServer) error {

	errorSendFunc := func(msg string) error {
		_ = stream.Send(&protocol.PublishResponse{
			OK:      false,
			Message: msg,
		})
		return nil
	}

	req, err := stream.Recv()
	if err != nil {
		return errorSendFunc(err.Error())
	}

	id, err := token.Decode(req.Token)
	if err != nil {
		return errorSendFunc(err.Error())
	}
	topicName := req.TopicName

	errCh := make(chan error, 1) // from sender
	pub := publisher.New(id, stream, errCh)

	type MsgAndErr struct {
		Msg *protocol.PublishRequest
		Err error
	}
	topic, err := topicManager.AddPublisher(topicName, pub)
	if err != nil {
		return errorSendFunc(err.Error())
	}

	exitInternalCh := make(chan struct{})
	defer func() { close(exitInternalCh) }()

	const BUFFERCOUNT = 32
	msgCh := make(chan *MsgAndErr, BUFFERCOUNT-2)
	var msgAndErrBuffer [BUFFERCOUNT]MsgAndErr

	go func() {
		i := 0
		for {
			select {
			case <-exitInternalCh:
				return
			default:
				m := &msgAndErrBuffer[i]
				m.Msg, m.Err = stream.Recv()
				select {
				case <-exitInternalCh:
					return
				case msgCh <- m:
				}
				i++
				if i == BUFFERCOUNT {
					i = 0
				}
			}
		}
	}()

	log.Infoln("Publish():", topicName, "by", id)

	for {
		select {
		case <-errCh:
			goto finalize
		case msgAndErr := <-msgCh:
			msg := msgAndErr.Msg
			err := msgAndErr.Err
			if err != nil {
				goto finalize
			}

			/*
				if msg.EOR {
					goto finalize
				}
			*/

			err = topic.SendToSubscriber(msg.UniqueSequenceNumber, msg.Start, msg.End, msg.Payload)
			if err != nil {
				_ = errorSendFunc(err.Error())
				goto finalize
			}

		}
	}

finalize:
	err = topicManager.RemovePublisher(topicName, pub)
	if err != nil {
		// TODO
	}
	pub.Close(err)
	return nil
}
