package server

import (
	"context"

	"github.com/iloy/wormhole/protocol"
)

// Statistics :
func (rs *rpcServer) Statistics(
	ctx context.Context, req *protocol.StatisticsRequest) (
	*protocol.StatisticsResponse, error) {

	return nil, nil
}
