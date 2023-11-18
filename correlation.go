package traefik_correlation

import (
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"os"
)

const (
	DEFAULT_HEADER_NAME = "correlation-id"
)

var (
	logDebug = os.Stdout.WriteString
	logError = os.Stderr.WriteString
)

// Config the plugin configuration.
type Config struct {
	HeaderName string `yaml:"headerName,omitempty" json:"header_name,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		HeaderName: "",
	}
}

type Correlation struct {
	next       http.Handler
	name       string
	headerName string
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.HeaderName == "" {
		logError("[Correlation] headerName option cannot be empty")
		return nil, fmt.Errorf("[Correlation] headerName option cannot be empty")
	}

	return &Correlation{
		next:       next,
		name:       name,
		headerName: config.HeaderName,
	}, nil
}

func (c *Correlation) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	logDebug("[Correlation][StdOut] Try to serve next")
	fmt.Println("[Correlation][Fmt] Try to serve next")
	log.Println("[Correlation][Log] Try to serve next")
	var id = uuid.NewV4().String()

	request.Header.Add(c.headerName, id)
	if request.Header.Get(c.headerName) != "" {
		logDebug(fmt.Sprintf("[Correlation] HeaderName by value %s is empty", c.headerName))
		request.Header.Add(c.headerName, request.Header.Get(c.headerName))
	}

	fmt.Printf("[Correlation][Fmt] All headers are incoming with correlationId: %s\n", id)
	log.Printf("[Correlation][Log] All headers are incoming with correlationId: %s\n", id)
	logDebug(fmt.Sprintf("[Correlation] All headers are incoming with correlationId: %s", id))
	c.next.ServeHTTP(writer, request)
}
