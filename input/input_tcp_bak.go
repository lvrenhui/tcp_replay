package input

// import (
// 	"fmt"
// 	"log"
// 	"net"
// 	"strings"
// 	"time"

// 	"github.com/google/gopacket"
// 	"github.com/google/gopacket/pcap"
// )

// // TCPInput used for internal communication
// type TCPInput struct {
// 	data     chan []byte
// 	listener net.Listener
// 	address  string
// 	// config   *TCPInputConfig
// }

// // type TCPInputConfig struct {
// // 	secure          bool
// // 	certificatePath string
// // 	keyPath         string
// // }

// // NewTCPInput constructor for TCPInput, accepts address with port
// func NewTCPInput(address string) (i *TCPInput) {
// 	i = new(TCPInput)
// 	i.data = make(chan []byte, 1000)
// 	i.address = address
// 	// i.config = config

// 	i.listen(address)

// 	return
// }

// func (i *TCPInput) Read(data []byte) (int, error) {
// 	buf := <-i.data
// 	copy(data, buf)

// 	return len(buf), nil
// }

// func (i *TCPInput) listen(address string) {

// 	//todo : auto detect device,like lo0 for 127.0.0.1
// 	var (
// 		device            = "lo0"
// 		snapshotLen int32 = 1024
// 		promiscuous       = false
// 		err         error
// 		timeout     time.Duration = 3 * time.Second
// 		handle      *pcap.Handle
// 	)

// 	// Open device
// 	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer handle.Close()

// 	// Set filter
// 	var filter = "tcp and dst port " + strings.Split(address, ":")[1]
// 	err = handle.SetBPFFilter(filter)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
// 	go func() {

// 		for packet := range packetSource.Packets() {
// 			// Do something with a packet here.
// 			// fmt.Println(packet)

// 			applicationLayer := packet.ApplicationLayer()
// 			if applicationLayer != nil {
// 				// fmt.Println("Application layer/Payload found.")
// 				fmt.Printf("%s\n", applicationLayer.Payload())
// 				i.data <- applicationLayer.Payload()
// 			}
// 		}
// 	}()

// }

// func (i *TCPInput) String() string {
// 	return "TCP input: " + i.address
// }
