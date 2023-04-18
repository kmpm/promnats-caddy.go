package promnatscaddy

import (
	"log"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {

}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (m *PromNats) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	log.Println("---UnmarshalCaddyfile")
	for d.Next() {

		if !d.Args(&m.Subject) {
			//return d.ArgErr()
			log.Println("empty subject")
			m.Subject = "metrics"
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "context":
				if !d.Args(&m.ContextName) {
					return d.ArgErr()
				}
			case "server":
				if !d.Args(&m.ServerURL) {
					d.ArgErr()
				}
			}
		}

	}
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	log.Println("---parseCaddyfile")
	m := newModule()
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}
