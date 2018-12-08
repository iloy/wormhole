package server

import (
	"context"
	"strconv"

	"github.com/iloy/wormhole/protocol"
	"github.com/iloy/wormhole/seqno"
	"github.com/iloy/wormhole/token"
)

var (
	errMsgInvalidProtocolVersion = "invalid protocol version"
	errMsgInvalidRole            = "invalid role"
	errMsgInternalError          = "internal error"
)

var (
	loginSeqno = seqno.Seqno{}
)

// Login :
func (ss *rpcServer) Login(
	ctx context.Context, req *protocol.LoginRequest) (
	*protocol.LoginResponse, error) {

	errorSendFunc := func(errMsg string) (*protocol.LoginResponse, error) {
		return &protocol.LoginResponse{
			OK:      false,
			Message: errMsg,
		}, nil
	}

	if req.ProtocolVersion != protocol.LoginRequest_VERSION_CURRENT {
		return errorSendFunc(errMsgInvalidProtocolVersion)
	}

	if req.Role != protocol.LoginRequest_ROLE_PUBLISHER &&
		req.Role != protocol.LoginRequest_ROLE_SUBSCRIBER &&
		req.Role != protocol.LoginRequest_ROLE_MONITOR {
		return errorSendFunc(errMsgInvalidRole)
	}

	// TODO?
	token, err := token.Encode(req.Id)
	if err != nil {
		return errorSendFunc(errMsgInternalError)
	}

	id := req.Role.String() + "_" + strconv.FormatUint(loginSeqno.Get(), 10)

	return &protocol.LoginResponse{
		OK:      true,
		Message: "welcome to wormhole",
		Token:   token,
		PublisherSubscriberMonitorId: id,
	}, nil
}
