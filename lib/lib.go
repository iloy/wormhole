package lib

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/iloy/wormhole/protocol"
)

// NewClient :
func NewClient(address string, pingTime uint64) (
	protocol.WormholeClient, *grpc.ClientConn) {

	kp := keepalive.ClientParameters{}
	kp.Time = time.Second * time.Duration(pingTime)
	kp.Timeout = time.Second * 20 // gRPC-go default value
	kp.PermitWithoutStream = true

	conn, err := grpc.Dial(address,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(kp),
	)
	if err != nil {
		panic(err)
	}

	rpcClient := protocol.NewWormholeClient(conn)

	return rpcClient, conn
}

// Login returns token, client ID and error
func Login(
	client protocol.WormholeClient, timeout uint64, uniqueClientName string) (
	string, string, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Second*time.Duration(timeout))
	defer func() { cancel() }()

	// TODO!
	req := &protocol.LoginRequest{
		ProtocolVersion: protocol.LoginRequest_VERSION_CURRENT,
		Id:              uniqueClientName,
		Password:        "",
		Role:            protocol.LoginRequest_ROLE_PUBLISHER,
	}

	res, err := client.Login(ctx, req)
	if err != nil {
		return "", "", err
	}
	if !res.OK {
		return "", "", errors.New(res.Message)
	}

	return res.Token, res.PublisherSubscriberMonitorId, nil
}

// CreateTopic :
func CreateTopic(
	client protocol.WormholeClient, timeout uint64, token string, topicName string) error {
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Second*time.Duration(timeout))
	defer func() { cancel() }()

	req := &protocol.CreateTopicRequest{
		Token:     token,
		TopicName: topicName,
	}

	res, err := client.CreateTopic(ctx, req)
	if err != nil {
		return err
	}
	if !res.OK {
		return errors.New(res.Message)
	}

	return nil
}

// DestroyTopic :
func DestroyTopic(
	client protocol.WormholeClient, timeout uint64, token string, topicName string) error {
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Second*time.Duration(timeout))
	defer func() { cancel() }()

	req := &protocol.DestroyTopicRequest{
		Token:     token,
		TopicName: topicName,
	}

	res, err := client.DestroyTopic(ctx, req)
	if err != nil {
		return err
	}
	if !res.OK {
		return errors.New(res.Message)
	}

	return nil
}

// Publish :
// TODO?
// remove the dependency to protocol.PublishResponse
func Publish(
	client protocol.WormholeClient,
	timeout uint64, token string, topicName string,
	dataSendFunc func() (uint64, bool, bool, []byte, error),
	dataRecvFunc func(*protocol.PublishResponse) error) error {
	var ctx context.Context
	var cancel context.CancelFunc
	if timeout == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	}
	defer func() { cancel() }()

	stream, err := client.Publish(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(&protocol.PublishRequest{
		Token:     token,
		TopicName: topicName,
	})

	var wg sync.WaitGroup
	defer wg.Wait()

	exitSender := make(chan struct{})
	defer close(exitSender)
	exitReceiver := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(exitReceiver)

		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}
			if msg == nil {
				return
			}
			if !msg.OK {
				exitReceiver <- errors.New(msg.Message)
				return
			}

			err = dataRecvFunc(msg)
			if err != nil {
				return
			}
		}
	}()

	// FIXME
	// parallelize between stream.Send() and dataSendFunc()
	for {
		select {
		case err := <-exitReceiver:
			return err
		default:
		}

		no, start, end, data, err := dataSendFunc()
		if err != nil && err != io.EOF {
			return err
		}

		if len(data) > 0 {
			req := &protocol.PublishRequest{
				UniqueSequenceNumber: no,
				Start:                start,
				End:                  end,
				Payload:              data,
				EOR:                  err == io.EOF,
			}

			err2 := stream.Send(req)
			if err2 != nil {
				return err2
			}
		}

		if err == io.EOF {
			return <-exitReceiver
		}
	}
}

// Subscribe :
// TODO?
// remove the dependency to protocol.SubscribeResponse
func Subscribe(
	client protocol.WormholeClient,
	timeout uint64, token string, topicName string,
	dataRecvFunc func(uint64, bool, bool, []byte) (*protocol.SubscribeRequest, error)) error {
	var ctx context.Context
	var cancel context.CancelFunc
	if timeout == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	}
	defer func() { cancel() }()

	stream, err := client.Subscribe(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(&protocol.SubscribeRequest{
		Token:     token,
		TopicName: topicName,
	})

	// FIXME
	// parallelize between stream.Recv() and dataReceiveFunc()
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		if !msg.OK {
			return errors.New(msg.Message)
		}

		req, err := dataRecvFunc(msg.UniqueSequenceNumber, msg.Start, msg.End, msg.Payload)
		if err != nil {
			return err
		}
		if req != nil {
			err := stream.Send(req)
			if err != nil {
				return err
			}
		}
	}
}
