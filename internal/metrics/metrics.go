package metrics

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/farishan-labs/dist/internal/plugins"
)

type UTMEntry struct {
	Source          string  `json:"source"`
	Count           int     `json:"count"`
	RevenueUSDCents float64 `json:"revenue_usd_cents"`
}

type GumroadMetrics struct {
	ProductsTotal     int        `json:"products_total"`
	ProductsPublished int        `json:"products_published"`
	ProductsDraft     int        `json:"products_draft"`
	SalesCount        int        `json:"sales_count"`
	RevenueUSDCents   float64    `json:"revenue_usd_cents"`
	RevenueUSD        float64    `json:"revenue_usd"`
	Views             int        `json:"views,omitempty"`
	ViewsSource       string     `json:"views_source,omitempty"`
	SalesListCount    int        `json:"sales_list_count"`
	UTMAttribution    []UTMEntry `json:"utm_attribution,omitempty"`
}

type GoatMetrics struct {
	Configured bool   `json:"configured"`
	Verified   bool   `json:"verified"`
	BaseURL    string `json:"base_url,omitempty"`
	Views      int    `json:"views"`
	Visitors   int    `json:"visitors"`
	Error      string `json:"error,omitempty"`
}

type Report struct {
	TS          string                      `json:"ts"`
	Gumroad     GumroadMetrics              `json:"gumroad"`
	GoatCounter GoatMetrics                 `json:"goatcounter"`
	Plugins     map[string]json.RawMessage `json:"plugins,omitempty"`
	Notes       []string                    `json:"notes,omitempty"`
}

func Collect() Report {
	rep := Report{TS: time.Now().UTC().Format(time.RFC3339)}
	notes := []string{}

	m, err := plugins.Load()
	if err != nil {
		notes = append(notes, "plugin manifest load failed: "+err.Error())
	}

	for _, p := range m.Plugins {
		if !p.Enabled {
			continue
		}
		out, err := plugins.Run(p)
		if err != nil {
			notes = append(notes, fmt.Sprintf("plugin %s failed: %v", p.Name, err))
			continue
		}
		var env plugins.Envelope
		if err := json.Unmarshal(out, &env); err != nil {
			notes = append(notes, fmt.Sprintf("plugin %s bad envelope: %v", p.Name, err))
			continue
		}
		notes = append(notes, env.Notes...)

		switch p.Name {
		case "gumroad":
			var gm GumroadMetrics
			if err := json.Unmarshal(env.Metrics, &gm); err != nil {
				notes = append(notes, fmt.Sprintf("plugin %s metrics parse failed: %v", p.Name, err))
			} else {
				rep.Gumroad = gm
			}
		case "goatcounter":
			var gc GoatMetrics
			if err := json.Unmarshal(env.Metrics, &gc); err != nil {
				notes = append(notes, fmt.Sprintf("plugin %s metrics parse failed: %v", p.Name, err))
			} else {
				rep.GoatCounter = gc
			}
		default:
			if rep.Plugins == nil {
				rep.Plugins = map[string]json.RawMessage{}
			}
			rep.Plugins[p.Name] = env.Metrics
		}
	}

	rep.Notes = notes
	return rep
}

func VerifyGoatCounter() error {
	m, err := plugins.Load()
	if err != nil {
		return err
	}
	p, idx := plugins.Find(m, "goatcounter")
	if idx < 0 {
		return fmt.Errorf("goatcounter plugin not configured")
	}
	out, err := plugins.Run(p)
	if err != nil {
		return err
	}
	var env plugins.Envelope
	if err := json.Unmarshal(out, &env); err != nil {
		return err
	}
	var gc GoatMetrics
	if err := json.Unmarshal(env.Metrics, &gc); err != nil {
		return err
	}
	if !gc.Verified {
		if gc.Error != "" {
			return fmt.Errorf("%s", gc.Error)
		}
		return fmt.Errorf("goatcounter not verified")
	}
	return nil
}
