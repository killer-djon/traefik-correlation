package traefik_correlation

import (
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"os"
)

const (
	DEFAULT_HEADER_NAME = "correlation-id"
)

var (
	logDebug = os.Stdout.WriteString
	//logError = os.Stderr.WriteString
)

// Config the plugin configuration.
type Config struct {
	HeaderName string `yaml:"headerName,omitempty" json:"header_name,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		HeaderName: DEFAULT_HEADER_NAME,
	}
}

type Correlation struct {
	next       http.Handler
	name       string
	headerName string
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &Correlation{
		next:       next,
		name:       name,
		headerName: config.HeaderName,
	}, nil
}

func (c *Correlation) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var id = uuid.NewV4().String()
	writer.Header().Add(c.headerName, id)
	if request.Header.Get(c.headerName) != "" {
		writer.Header().Add(c.headerName, request.Header.Get(c.headerName))
	}

	logDebug(fmt.Sprintf("All headers are incoming with correlationId: %s", id))
	c.next.ServeHTTP(writer, request)
}
