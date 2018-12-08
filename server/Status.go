package server

import (
	"context"

	"github.com/iloy/wormhole/protocol"
)

// Status :
func (rs *rpcServer) Status(
	ctx context.Context, req *protocol.StatusRequest) (
	*protocol.StatusResponse, error) {

	return nil, nil
}
