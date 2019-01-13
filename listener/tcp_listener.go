package listener

import (
	"log"
	"strconv"

	"github.com/lvrenhui/tcp_replay/proto"
)

type TCPListener struct {
	// IP to listen
	addr string
	// Port to listen
	port uint16

	messagesChan chan *proto.TCPMessage

	underlying *IPListener
}

func NewTCPListener(addr string, port string, trackResponse bool) (l *TCPListener) {
	l = &TCPListener{}
	l.messagesChan = make(chan *proto.TCPMessage, 10000)
	l.addr = addr
	intPort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Invaild Port: %s, %v\n", port, err)
	}
	l.port = uint16(intPort)

	l.underlying = NewIPListener(addr, l.port, trackResponse)

	if l.underlying.IsReady() {
		go l.recv()
	} else {
		log.Fatalln("IP Listener is not ready after 5 seconds")
	}

	return
}

func (l *TCPListener) parseTCPPacket(packet *ipPacket) (message *proto.TCPMessage) {
	data := packet.payload
	message = proto.NewTCPMessage(data, false)
	if message.DstPort == l.port {
		message.IsIncoming = true
	}
	message.Start = packet.timestamp
	return
}

func (l *TCPListener) recv() {
	for {
		ipPacketsChan := l.underlying.Receiver()
		select {
		case packet := <-ipPacketsChan:
			message := l.parseTCPPacket(packet)
			l.messagesChan <- message
		}
	}
}

func (l *TCPListener) Receiver() chan *proto.TCPMessage {
	return l.messagesChan
}
