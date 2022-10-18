package mixpanelproxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/tidwall/sjson"
)

var keysToClear = []string{"$referrer", "$referring_domain", "$current_url"}

const clearedValue = "CLEARED_BY_MIXPANEL_PROXY"

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
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m MixpanelProxy) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	return next.ServeHTTP(w, r)
}

func (m MixpanelProxy) MassageRequestBody(r *http.Request) error {
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("mixpanel_proxy: Unable to read body: %s", err)
	}

	data, err = sjson.SetBytes(data, "#.properties.token", m.MixpanelKey)
	if err != nil {
		return fmt.Errorf("mixpanel_proxy: Unable to set token value: %s", err)
	}

	for _, toClear := range keysToClear {
		data, err = sjson.SetBytes(data, "#.properties."+toClear, clearedValue)
		if err != nil {
			return fmt.Errorf("mixpanel_proxy: Unable to set clear %s key: %s", toClear, err)
		}
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	return nil
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
