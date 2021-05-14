package messaging

//go:generate mockgen -source=interface.go -destination=mocks/mock_interface.go -package=mock_gox_aws_messaging

type Event struct {
	Key   string
	Value interface{}
}

type Response struct {
	RawPayload interface{}
}

type Producer interface {
	Send(request *Event) (*Response, error)
}
