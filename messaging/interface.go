package messaging

import "context"

//go:generate mockgen -source=interface.go -destination=mocks/mock_interface.go -package=mock_gox_aws_messaging

type Event struct {
	Key      string
	Value    interface{}
	RawEvent interface{}
}

type Response struct {
	RawPayload interface{}
}

type Producer interface {
	Send(request *Event) (*Response, error)
}

type ConsumerFunc func(messageChannel chan Event)

type Consumer interface {
	Process(ctx context.Context, messagePressedAckChannel chan Event) (chan Event, error)
}
