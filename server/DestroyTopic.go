package server

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/iloy/wormhole/protocol"
	"github.com/iloy/wormhole/token"
)

// DestroyTopic :
func (rs *rpcServer) DestroyTopic(
	ctx context.Context, req *protocol.DestroyTopicRequest) (
	*protocol.DestroyTopicResponse, error) {

	errorSendFunc := func(errMsg string) (*protocol.DestroyTopicResponse, error) {
		return &protocol.DestroyTopicResponse{
			OK:      false,
			Message: errMsg,
		}, nil
	}

	id, err := token.Decode(req.Token)
	if err != nil {
		return errorSendFunc(errMsgInvalidToken)
	}

	if req.TopicName == "" {
		return errorSendFunc(errMsgInvalidTopicName)
	}

	err = topicManager.RemoveTopic(req.TopicName, id)
	if err != nil {
		return errorSendFunc(err.Error())
	}
	log.Infoln("DestroyTopic():", req.TopicName)

	return &protocol.DestroyTopicResponse{
		OK:      true,
		Message: "",
	}, nil
}
