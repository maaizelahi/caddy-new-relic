package newrelic

import (
	"fmt"
	"net/http"
	"os"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/newrelic/go-agent/v3/newrelic"
)

var app *newrelic.Application
var nrErr error

// app, err := newrelic.NewApplication(
//     newrelic.ConfigAppName("Your Application Name"),
//     newrelic.ConfigLicense("YOUR_NEW_RELIC_LICENSE_KEY")
// )
// http.HandleFunc(newrelic.WrapHandleFunc(app, "/users", usersHandler))

// app, err := newrelic.NewApplication(
//     newrelic.ConfigFromEnvironment()
// )

// const (
// 	nextCallCtxKey  caddy.CtxKey = "nextCall"
// )

// // nextCall store the next handler, and the error value return on calling it (if any)
// type nextCall struct {
// 	next caddyhttp.Handler
// 	err  error
// }

func init() {
	caddy.RegisterModule(Newrelic{})

	 // Create an Application:
	//  app, err := newrelic.NewApplication(
    //     // Name your application
    //     // // Fill in your New Relic license key
    //     // Add logging:
	// 	newrelic.ConfigFromEnvironment(),
	// 	newrelic.ConfigDebugLogger(os.Stdout),
	// 	newrelic.ConfigAppLogForwardingEnabled(true),
	// 	newrelic.ConfigCodeLevelMetricsEnabled(true),
    //     // Optional: add additional changes to your configuration via a config function:
    //     // func(cfg *newrelic.Config) {
    //     //     cfg.CustomInsightsEvents.Enabled = false
    //     // },
    // )
    // // If an application could not be created then err will reveal why.
    // if err != nil {
    //     fmt.Println("unable to create New Relic Application", err)
    // }
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

// Provision implements caddy.Provisioner.
func (nr *Newrelic) Provision(ctx caddy.Context) error {
	app, nrErr = newrelic.NewApplication(
		newrelic.ConfigAppName("maaiz-caddy-module-service"),
        // Fill in your New Relic license key
        newrelic.ConfigLicense("c6547b52947fb5d2adcd712deccda1339028NRAL"),
		// newrelic.ConfigFromEnvironment(),
		newrelic.ConfigDebugLogger(os.Stdout),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigCodeLevelMetricsEnabled(true),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigCodeLevelMetricsPathPrefix("go-agent/v3"),
	)

	   // If an application could not be created then err will reveal why.
	if nrErr != nil {
        fmt.Println("unable to create New Relic Application", nrErr)
    }
	if app != nil {
        fmt.Println("Maaiz app initialized", app)
    }
    // Now use the app to instrument everything!
	return nrErr
}

// serveHTTP injects a tracing context and call the next handler.
// func (nr *Newrelic) serveHTTP(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Maaiz in small serve http");
// 	next := r.Context().Value(nextCallCtxKey).(*nextCall)
// 	next.err = next.next.ServeHTTP(w, r)
// }

// func (nr *Newrelic) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
// 	fmt.Println("Maaiz in ServeHTTP")
// 	dummyHandler(w, r)
// 	// func WrapHandle(app Application, pattern string, handler http.Handler) (string, http.Handler)
// 	path, handler := newrelic.WrapHandle(app, "/", dummyHandler);
// 	fmt.Println(path);
// 	fmt.Println(handler);
// 	// http.HandleFunc("/", dummyHandler)
// 	fmt.Println("Maaiz after dummyHandler")
// 	return next.ServeHTTP(w, r)
// }

func (nr *Newrelic) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	fmt.Println("Maaiz in ServeHTTP")
	txn := app.StartTransaction(r.Method + " " + r.Host + r.URL.Path)
	// txn := app.StartTransaction(r.Method+" "+pattern, txnOptionList...)
	defer txn.End()

	w = txn.SetWebResponse(w)
	txn.SetWebRequestHTTP(r)

	r = newrelic.RequestWithTransactionContext(r, txn)
	
	return next.ServeHTTP(w, r)
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
// func (ot *Tracing) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
// 	return ot.otel.ServeHTTP(w, r, next)
// }

// Interface guards
var (
	_ caddy.Provisioner           = (*Newrelic)(nil)
	_ caddyhttp.MiddlewareHandler = (*Newrelic)(nil)
	_ caddy.Module      		  = (*Newrelic)(nil)
)