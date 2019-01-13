package main

import (
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/lvrenhui/tcp_replay/input"
	"github.com/lvrenhui/tcp_replay/output"
)

// InOutPlugins struct for holding references to plugins
type InOutPlugins struct {
	Inputs  []io.Reader
	Outputs []io.Writer
	All     []interface{}
}

var pluginMu sync.Mutex

// Plugins holds all the plugin objects
var Plugins = new(InOutPlugins)

// extractLimitOptions detects if plugin get called with limiter support
// Returns address and limit
func extractLimitOptions(options string) (string, string) {
	split := strings.Split(options, "|")

	if len(split) > 1 {
		return split[0], split[1]
	}

	return split[0], ""
}

// Automatically detects type of plugin and initialize it
//
// See this article if curious about reflect stuff below: http://blog.burntsushi.net/type-parametric-functions-golang
func registerPlugin(constructor interface{}, options ...interface{}) {
	var path, limit string
	vc := reflect.ValueOf(constructor)

	// Pre-processing options to make it work with reflect
	vo := []reflect.Value{}
	for _, oi := range options {
		vo = append(vo, reflect.ValueOf(oi))
	}

	if len(vo) > 0 {
		// Removing limit options from path
		path, limit = extractLimitOptions(vo[0].String())

		// Writing value back without limiter "|" options
		vo[0] = reflect.ValueOf(path)
	}

	// Calling our constructor with list of given options
	plugin := vc.Call(vo)[0].Interface()

	if limit != "" {
		plugin = NewLimiter(plugin, limit)
	}

	_, isR := plugin.(io.Reader)
	_, isW := plugin.(io.Writer)

	// Some of the output can be Readers as well because return responses
	if isR && !isW {
		Plugins.Inputs = append(Plugins.Inputs, plugin.(io.Reader))
	}

	if isW {
		Plugins.Outputs = append(Plugins.Outputs, plugin.(io.Writer))
	}

	Plugins.All = append(Plugins.All, plugin)
}

// InitPlugins specify and initialize all available plugins
func InitPlugins() {
	pluginMu.Lock()
	defer pluginMu.Unlock()

	for _, options := range Settings.inputTCP {
		registerPlugin(input.NewTCPInput, options, false)
	}

	for _, options := range Settings.inputFile {
		registerPlugin(input.NewFileInput, options, Settings.inputFileLoop)
	}
	if Settings.outputStdout {
		registerPlugin(output.NewStdOutput)
	}

	for _, options := range Settings.outputFile {
		registerPlugin(output.NewFileOutput, options, &Settings.outputFileConfig)
	}

}
