package traefik_correlation

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
	"net/http"
)

const (
	DEFAULT_HEADER_NAME = "correlation-id"
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
		return nil, fmt.Errorf("[Correlation] headerName option cannot be empty")
	}

	return &Correlation{
		next:       next,
		name:       name,
		headerName: config.HeaderName,
	}, nil
}

func (c *Correlation) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Debug().Msg("[Correlation] Try to serve next")
	var id = uuid.NewV4().String()

	if request.Header.Get(c.headerName) == "" {
		log.Debug().Msgf("[Correlation] HeaderName by value %s is empty", c.headerName)
		request.Header.Add(c.headerName, id)
	}

	log.Debug().Msgf("[Correlation] All headers are incoming with correlationId: %s", id)
	c.next.ServeHTTP(writer, request)
}
