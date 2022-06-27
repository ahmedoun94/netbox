package netbox

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

// init registers this plugin.
func init() { plugin.Register("netbox", setup) }

//We use the *caddy.Controller to receive tokens from the Corefile and act upon them. Here we only check if there is nothing specified after the token netbox
func setup(c *caddy.Controller) error {
	n, err := parseNetbox(c)
	if err != nil {
		return plugin.Error("netbox", err)
	}
	c.Next()
	if c.NextArg() {
		return plugin.Error("netbox", c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return n // Set the Next field, so the plugin chaining works.
	})
	//all OK, return a nil error.
	return nil
}

func newNetbox() *Netbox {
	return &Netbox{}
}

func parseNetbox(c *caddy.Controller) (*Netbox, error) {
	n := newNetbox()
	i := 0
	for c.Next() {
		// ensure plugin is only included once in each block
		if i > 0 {
			return nil, plugin.ErrOnce
		}
		i++

		// parse inside block
		for c.NextBlock() {
			switch c.Val() {

			case "url":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				n.Url = c.Val()

			case "token":
				if !c.NextArg() {
					return n, c.ArgErr()
				}
				n.Token = c.Val()

			default:
				return nil, c.Errf("unknown property '%s'", c.Val())
			}
		}
	}

	// fail if url, token or localCacheDuration are not set
	if n.Url == "" || n.Token == "" {
		return nil, c.Err("Invalid config")
	}

	return n, nil
}
