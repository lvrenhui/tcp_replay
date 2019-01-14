package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/lvrenhui/tcp_replay/output"
)

// MultiOption allows to specify multiple flags with same name and collects all values into array
type MultiOption []string

func (h *MultiOption) String() string {
	return fmt.Sprint(*h)
}

// Set gets called multiple times for each flag with same name
func (h *MultiOption) Set(value string) error {
	*h = append(*h, value)
	return nil
}

// AppSettings is the struct of main configuration
type AppSettings struct {
	exitAfter time.Duration

	splitOutput  bool
	outputStdout bool
	outputNull   bool

	inputFile        MultiOption
	inputFileLoop    bool
	outputFile       MultiOption
	outputFileConfig output.FileOutputConfig

	inputTCP        MultiOption
	outputTCP       MultiOption
	tcpOutputConfig output.TCPOutputConfig
}

// Settings holds Goreplay configuration
var Settings AppSettings

func init() {
	flag.DurationVar(&Settings.exitAfter, "exit-after", 0, "exit after specified duration")

	flag.BoolVar(&Settings.splitOutput, "split-output", false, "By default each output gets same traffic. If set to `true` it splits traffic equally among all outputs")

	flag.Var(&Settings.inputTCP, "input-tcp", "Capture traffic from given port (use RAW sockets and require *sudo* access):\n\t# Capture traffic from 8080 port\n\ttcp_replay --input-raw :8080 --output-stdout")

	flag.Var(&Settings.inputFile, "input-file", "Read requests from file: \n\ttcp_replay --input-file ./requests.gor --output-stdout")
	flag.BoolVar(&Settings.inputFileLoop, "input-file-loop", false, "Loop input files, useful for performance testing")

	flag.Var(&Settings.outputFile, "output-file", "Write incoming requests to file: \n\ttcp_replay --input-tcp :80 --output-file ./requests.gor")
	flag.DurationVar(&Settings.outputFileConfig.FlushInterval, "output-file-flush-interval", time.Second, "Interval for forcing buffer flush to the file, default: 1s")
	flag.BoolVar(&Settings.outputFileConfig.Append, "output-file-append", false, "The flushed chunk is appended to existence file or not")

	flag.BoolVar(&Settings.outputStdout, "output-stdout", false, "Used for testing inputs. Just prints to console data coming from inputs")

	flag.Var(&Settings.outputTCP, "output-tcp", "Used for out put to tcp address like:\n\t tcp_replay --input-file pcap.out --output-tcp 127.0.0.1:4000")
	flag.BoolVar(&Settings.tcpOutputConfig.Secure, "output-tcp-secure", false, "Use TLS secure connection. --input-file on another end should have TLS turned on as well.")
	flag.BoolVar(&Settings.tcpOutputConfig.Stats, "output-tcp-stats", false, "Report TCP output queue stats to console every 5 seconds.")
	flag.IntVar(&Settings.tcpOutputConfig.Repeat, "output-tcp-repeat", 1, "Reapt times for each request for perf testing to .")

	// Set default
	Settings.outputFileConfig.SizeLimit.Set("32mb")
	flag.Var(&Settings.outputFileConfig.SizeLimit, "output-file-size-limit", "Size of each chunk. Default: 32mb")
	flag.IntVar(&Settings.outputFileConfig.QueueLimit, "output-file-queue-limit", 25600, "The length of the chunk queue. Default: 25600")

}
