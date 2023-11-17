package traefik_correlation

import (
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"os"
)

var (
	logDebug = os.Stdout.WriteString
	//logError = os.Stderr.WriteString
)

// Config the plugin configuration.
type Config struct {
	Headers []string `yaml:"headers,omitempty" json:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: []string{},
	}
}

type Correlation struct {
	next    http.Handler
	name    string
	headers map[string]bool
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	var headers = make(map[string]bool, len(config.Headers))

	if len(config.Headers) > 0 {
		logDebug(fmt.Sprintf("Finded %d count headers at request", len(config.Headers)))
		for _, name := range config.Headers {
			logDebug(fmt.Sprintf("Incoming header by name: %s", name))
			headers[name] = true
		}
	}

	return &Correlation{
		headers: headers,
		next:    next,
		name:    name,
	}, nil
}

func (c *Correlation) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var id = uuid.NewV4().String()
	logDebug(fmt.Sprintf("All headers are incoming with correlationId: %s", id))
	c.next.ServeHTTP(writer, request)
}
