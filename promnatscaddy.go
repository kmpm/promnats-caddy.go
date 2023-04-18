package promnatscaddy

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func init() {
	log.Println("---init")
	caddy.RegisterModule(PromNats{})
	httpcaddyfile.RegisterHandlerDirective("promnats", parseCaddyfile)

}

// CaddyModule returns the Caddy module information.
func (PromNats) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.promnats",
		New: func() caddy.Module { return newModule() },
	}
}

func newModule() *PromNats {
	log.Println("---newModule")
	pn := &PromNats{
		Interval: time.Minute * 5,
		routes:   make(map[string][]byte),
	}
	return pn
}

func (m *PromNats) Provision(ctx caddy.Context) error {
	log.Println("---Provision")
	m.logger = ctx.Logger()

	var err error
	var nc *nats.Conn
	if m.ServerURL == "" {
		nc, err = natscontext.Connect(m.ContextName)
	} else {
		nc, err = nats.Connect(m.ServerURL)
	}
	if err != nil {
		return err
	}
	m.logger.Info("nats connected", zap.Strings("servers", nc.DiscoveredServers()))
	m.nc = nc
	go m.request()
	m.refresh()
	return nil
}

// Validate implements caddy.Validator.
func (m *PromNats) Validate() error {
	if m.ContextName == "" && m.ServerURL == "" {
		return fmt.Errorf("no context or server")
	}
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m *PromNats) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// m.w.Write([]byte(r.RemoteAddr))
	m.logger.Info("request", zap.String("path", r.URL.Path))

	if data, ok := m.routes[r.URL.Path]; ok {
		w.Write(data)
		return nil
	}
	return caddyhttp.Error(http.StatusNotFound, fmt.Errorf("no such path"))

	// return next.ServeHTTP(w, r)
}

// Interface guards
var (
	_ caddy.Provisioner           = (*PromNats)(nil)
	_ caddy.Validator             = (*PromNats)(nil)
	_ caddyhttp.MiddlewareHandler = (*PromNats)(nil)
	_ caddyfile.Unmarshaler       = (*PromNats)(nil)
)
