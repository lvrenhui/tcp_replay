package proto

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type TCPMessage struct {
	IsIncoming bool
	Start      time.Time
	SrcPort    uint16
	DstPort    uint16
	length     uint16
	checksum   uint16
	data       []byte
}

func NewTCPMessage(data []byte, isIncoming bool) (m *TCPMessage) {
	m = &TCPMessage{}
	tcp := &layers.TCP{}
	err := tcp.DecodeFromBytes(data, gopacket.NilDecodeFeedback)
	if err != nil {
		log.Printf("Error decode tcp message, %v\n", err)
	}
	m.SrcPort = uint16(tcp.SrcPort)
	m.DstPort = uint16(tcp.DstPort)
	//todo : fix length
	// m.length = tcp.Length
	m.length = uint16(len(tcp.Payload))
	m.checksum = tcp.Checksum
	m.data = tcp.Payload
	m.IsIncoming = isIncoming

	return
}

func (m *TCPMessage) UUID() []byte {
	var key []byte

	key = strconv.AppendInt(key, m.Start.UnixNano(), 10)
	key = strconv.AppendUint(key, uint64(m.SrcPort), 10)
	key = strconv.AppendUint(key, uint64(m.DstPort), 10)
	key = strconv.AppendUint(key, uint64(m.length), 10)

	uuid := make([]byte, 40)
	sha := sha1.Sum(key)
	hex.Encode(uuid, sha[:20])

	return uuid
}

func (m *TCPMessage) Data() []byte {
	return m.data
}

func (m *TCPMessage) String() string {
	return fmt.Sprintf("SrcPort: %d | DstPort: %d | Length: %d | Checksum: %d | Data: %s",
		m.SrcPort, m.DstPort, m.length, m.checksum, string(m.data))
}
