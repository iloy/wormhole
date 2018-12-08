package server

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/iloy/wormhole/protocol"
	"github.com/iloy/wormhole/token"
)

var (
	errMsgInvalidTopicName = "invalid topic name"
	errMsgInvalidToken     = "invalid token"
)

// CreateTopic :
func (rs *rpcServer) CreateTopic(
	ctx context.Context, req *protocol.CreateTopicRequest) (
	*protocol.CreateTopicResponse, error) {

	errorSendFunc := func(errMsg string) (*protocol.CreateTopicResponse, error) {
		return &protocol.CreateTopicResponse{
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

	// TODO?
	_, err = topicManager.AddTopic(req.TopicName, id)
	if err != nil {
		return errorSendFunc(err.Error())
	}
	log.Infoln("CreateTopic():", req.TopicName, "by", id)

	return &protocol.CreateTopicResponse{
		OK:      true,
		Message: "",
	}, nil
}
