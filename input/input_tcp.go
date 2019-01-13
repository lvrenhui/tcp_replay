package input

import (
	"log"
	"net"

	"github.com/lvrenhui/tcp_replay/listener"
	"github.com/lvrenhui/tcp_replay/proto"
)

type TCPInput struct {
	data          chan *proto.TCPMessage
	address       string
	quit          chan bool
	listener      *listener.TCPListener
	trackResponse bool
}

func NewTCPInput(address string, trackResponse bool) (i *TCPInput) {
	i = new(TCPInput)
	i.data = make(chan *proto.TCPMessage)
	i.address = address
	i.quit = make(chan bool)
	i.trackResponse = trackResponse
	i.listen(address)
	return
}

func (i *TCPInput) Read(data []byte) (int, error) {
	msg := <-i.data
	buf := msg.Data()

	var header []byte

	if msg.IsIncoming {
		header = proto.PayloadHeader(proto.RequestPayload, msg.UUID(), msg.Start.UnixNano())
	} else {
		header = proto.PayloadHeader(proto.ResponsePayload, msg.UUID(), msg.Start.UnixNano())
	}

	copy(data[0:len(header)], header)

	copy(data[len(header):], buf)

	// copy(data, buf)

	return len(buf) + len(header), nil
}

func (i *TCPInput) listen(address string) {
	log.Println("Listening for traffic on: " + address)

	host, port, err := net.SplitHostPort(address)

	if err != nil {
		log.Fatal("input-raw: error while parsing address", err)
	}

	i.listener = listener.NewTCPListener(host, port, i.trackResponse)

	ch := i.listener.Receiver()

	go func() {
		for {
			select {
			case <-i.quit:
				return
			default:
			}
			// Receiving TCPMessage
			m := <-ch
			i.data <- m
		}
	}()
}
