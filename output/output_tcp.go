package output

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/lvrenhui/tcp_replay/proto"
	"github.com/lvrenhui/tcp_replay/stats"
)

// TCPOutput used for sending raw tcp payloads
// Currently used for internal communication between listener and replay server
// Can be used for transfering binary payloads like protocol buffers
type TCPOutput struct {
	address  string
	limit    int
	buf      chan []byte
	bufStats *stats.GorStat
	config   *TCPOutputConfig
}

type TCPOutputConfig struct {
	Secure bool
	Stats  bool
	Repeat int
}

// NewTCPOutput constructor for TCPOutput
// Initialize 10 workers which hold keep-alive connection
func NewTCPOutput(address string, config *TCPOutputConfig) io.Writer {
	o := new(TCPOutput)

	o.address = address
	o.config = config

	o.buf = make(chan []byte, 1000)
	if o.config.Stats {
		o.bufStats = stats.NewGorStat("output_tcp")
	}

	for i := 0; i < 10; i++ {
		go o.worker()
	}

	return o
}

func (o *TCPOutput) worker() {
	retries := 0
	conn, err := o.connect(o.address)
	// log.Println(o.address)
	for {
		if err == nil || retries > 10 {
			break
		}

		log.Println("Can't connect to aggregator instance, reconnecting in 1 second. Retries:", retries)
		time.Sleep(1 * time.Second)

		conn, err = o.connect(o.address)
		retries++
	}

	if retries > 0 {
		log.Println("Connected to aggregator instance after ", retries, " retries")
	}

	defer conn.Close()

	for {
		data := <-o.buf

		// get body data ,ignore header
		body := proto.PayloadBody(data)

		num := o.config.Repeat
		for i := 0; i < num; i++ {

			conn.Write(body)
			_, err := conn.Write([]byte(proto.PayloadSeparator))

			if err != nil {
				log.Println("ERROR: TCP output connection closed!")
				// not retry
				// o.buf <- data
				// go o.worker()
				// break
			}

		}
	}
}

func (o *TCPOutput) Write(data []byte) (n int, err error) {
	if !proto.IsRequestPayload(data) {
		return len(data), nil
	}

	// We have to copy, because sending data in multiple threads
	newBuf := make([]byte, len(data))
	copy(newBuf, data)

	o.buf <- newBuf

	if o.config.Stats {
		o.bufStats.Write(len(o.buf))
	}

	return len(data), nil
}

func (o *TCPOutput) connect(address string) (conn net.Conn, err error) {
	if o.config.Secure {
		conn, err = tls.Dial("tcp", address, &tls.Config{})
	} else {
		conn, err = net.Dial("tcp", address)
	}

	return
}

func (o *TCPOutput) String() string {
	return fmt.Sprintf("TCP output %s, limit: %d", o.address, o.limit)
}
