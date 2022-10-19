package mixpanelproxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/tidwall/sjson"
)

var keysToClear = map[string]string{
	"$referrer":                 "https://example.com",
	"$referring_domain":         "example.com",
	"$current_url":              "https://example.com:1443/admin#/status/dashboard",
	"$initial_referrer":         "https://example.com:1443/admin",
	"$initial_referring_domain": "example.com",
}

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
	err := m.MassageRequestBody(r)

	if err == nil {
		return next.ServeHTTP(w, r)
	} else {
		fmt.Println("Error while massaging the request body", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}
}

func (m MixpanelProxy) MassageRequestBody(r *http.Request) error {
	defer r.Body.Close()

	bodyRaw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("mixpanel_proxy: Unable to read body: %s", err)
	}

	query, err := url.ParseQuery(string(bodyRaw))
	if err != nil {
		return fmt.Errorf("mixpanel_proxy: Unable to parse body: %s", err)
	}

	data := query.Get("data")
	if data == "" {
		return fmt.Errorf("mixpanel_proxy: Unable to read 'data' in body")
	}

	data, err = sjson.Set(data, "#.properties.token", m.MixpanelKey)
	if err != nil {
		return fmt.Errorf("mixpanel_proxy: Unable to set token value: %s", err)
	}

	for toClear, v := range keysToClear {
		data, err = sjson.Set(data, "#.properties."+toClear, v)
		if err != nil {
			return fmt.Errorf("mixpanel_proxy: Unable to set clear %s key: %s", toClear, err)
		}
	}

	newBody := "data=" + url.QueryEscape(data)
	newBody = strings.Replace(newBody, "+", "%20", -1)
	newBodyBytes := []byte(newBody)
	r.Header.Set("Content-Length", strconv.Itoa(len(newBodyBytes)))
	r.ContentLength = int64(len(newBodyBytes))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(newBodyBytes))

	return nil
}

// UnmarshalCaddyfile - this is a no-op
func (m *MixpanelProxy) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	arg1 := d.NextArg()
	arg2 := d.NextArg()

	// Parse standalone length
	if arg1 && arg2 {
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
