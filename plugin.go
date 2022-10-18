package mixpanelproxy

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(MixpanelProxy{})
	httpcaddyfile.RegisterHandlerDirective("mixpanel_proxy", parseCaddyfile)
}

// MixpanelProxy implements an HTTP handler that dynamically adds the key to a mixpanel payload
type MixpanelProxy struct {
	MixpanelKey string `json:"mixpanel_key"`
}

// CaddyModule returns the Caddy module information.
func (MixpanelProxy) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.mixpanel_proxy",
		New: func() caddy.Module { return new(MixpanelProxy) },
	}
}

// Provision implements caddy.Provisioner.
func (m *MixpanelProxy) Provision(ctx caddy.Context) error {
	if m.Length < 1 {
		m.Length = 21
	}

	if m.Additional == nil {
		m.Additional = make(map[string]int)
	}

	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m MixpanelProxy) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile - this is a no-op
func (m *MixpanelProxy) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	arg1 := d.NextArg()

	// Parse standalone length
	if arg1 {
		val := d.Val()
		m.MixpanelKey = val

		if m.MixpanelKey == "" {
			return d.Err("empty mixpanel key")
		}
	} else {
		return d.Err("missing mixpanel key")
	}

	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	m := new(MixpanelProxy)
	err := m.UnmarshalCaddyfile(h.Dispenser)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Interface guards
var (
	_ caddy.Provisioner           = (*MixpanelProxy)(nil)
	_ caddyhttp.MiddlewareHandler = (*MixpanelProxy)(nil)
	_ caddyfile.Unmarshaler       = (*MixpanelProxy)(nil)
)
