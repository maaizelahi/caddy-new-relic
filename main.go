package newrelic

import (
	"fmt"
	"net/http"
	"os"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// app, err := newrelic.NewApplication(
//     newrelic.ConfigAppName("Your Application Name"),
//     newrelic.ConfigLicense("YOUR_NEW_RELIC_LICENSE_KEY")
// )
// http.HandleFunc(newrelic.WrapHandleFunc(app, "/users", usersHandler))

// app, err := newrelic.NewApplication(
//     newrelic.ConfigFromEnvironment()
// )

func init() {
	caddy.RegisterModule(Newrelic{})

	 // Create an Application:
	 app, err := newrelic.NewApplication(
        // Name your application
        // // Fill in your New Relic license key
        // Add logging:
		newrelic.ConfigFromEnvironment(),
		newrelic.ConfigDebugLogger(os.Stdout),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigCodeLevelMetricsEnabled(true),
        // Optional: add additional changes to your configuration via a config function:
        // func(cfg *newrelic.Config) {
        //     cfg.CustomInsightsEvents.Enabled = false
        // },
    )
    // If an application could not be created then err will reveal why.
    if err != nil {
        fmt.Println("unable to create New Relic Application", err)
    }
    // Now use the app to instrument everything!
}

type Newrelic struct {
	// SpanName is a span name. It should follow the naming guidelines here:
	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/api.md#span
	// SpanName string `json:"span"`

	// otel implements opentelemetry related logic.
	// otel openTelemetryWrapper

	// logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (Newrelic) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.newrelic",
		New: func() caddy.Module { return new(Newrelic) },
	}
}

func (m *Newrelic) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	return next.ServeHTTP(w, r)
}

// Interface guards
var (
	// _ caddy.Provisioner           = (*Newrelic)(nil)
	_ caddyhttp.MiddlewareHandler = (*Newrelic)(nil)
	_ caddy.Module      		  = (*Newrelic)(nil)
)