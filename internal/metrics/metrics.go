package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Product struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Published     bool    `json:"published"`
	Price         float64 `json:"price"`
	SalesCount    int     `json:"sales_count"`
	SalesUSDCents float64 `json:"sales_usd_cents"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
}

type Report struct {
	TS          string         `json:"ts"`
	Gumroad     GumroadMetrics `json:"gumroad"`
	GoatCounter GoatMetrics    `json:"goatcounter"`
	Notes       []string       `json:"notes,omitempty"`
}

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

func runGumroad(args ...string) ([]byte, error) {
	cmd := exec.Command("gumroad", append(args, "--json")...)
	// gumroad CLI uses its own stored auth; an MCP token in env can confuse it
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "GUMROAD_ACCESS_TOKEN=") {
			continue
		}
		env = append(env, e)
	}
	cmd.Env = env
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return out, fmt.Errorf("%w: %s", err, string(ee.Stderr))
		}
		return out, err
	}
	return out, nil
}

func Collect() Report {
	rep := Report{TS: time.Now().UTC().Format(time.RFC3339)}
	notes := []string{}

	gm := GumroadMetrics{}
	prodOut, err := runGumroad("products", "list", "--all")
	if err != nil {
		notes = append(notes, "gumroad products list failed: "+err.Error())
	} else {
		var pr ProductsResponse
		if err := json.Unmarshal(prodOut, &pr); err != nil {
			notes = append(notes, "gumroad products parse failed: "+err.Error())
		} else {
			gm.ProductsTotal = len(pr.Products)
			for _, p := range pr.Products {
				if p.Published {
					gm.ProductsPublished++
				} else {
					gm.ProductsDraft++
				}
				gm.SalesCount += p.SalesCount
				gm.RevenueUSDCents += p.SalesUSDCents
			}
			gm.RevenueUSD = gm.RevenueUSDCents / 100.0
		}
	}

	salesOut, err := runGumroad("sales", "list", "--all")
	utm := map[string]*UTMEntry{}
	if err != nil {
		notes = append(notes, "gumroad sales list failed: "+err.Error())
	} else {
		addSale := func(m map[string]any) {
			src := utmSource(m)
			e := utm[src]
			if e == nil {
				e = &UTMEntry{Source: src}
				utm[src] = e
			}
			e.Count++
			e.RevenueUSDCents += saleRevenueCents(m)
		}
		var arr []map[string]any
		if err := json.Unmarshal(salesOut, &arr); err == nil {
			gm.SalesListCount = len(arr)
			for _, m := range arr {
				addSale(m)
			}
		} else {
			var obj map[string]any
			if err := json.Unmarshal(salesOut, &obj); err == nil {
				if sl, ok := obj["sales"].([]any); ok {
					gm.SalesListCount = len(sl)
					for _, item := range sl {
						if m, ok := item.(map[string]any); ok {
							addSale(m)
						}
					}
				}
			} else {
				notes = append(notes, "gumroad sales parse failed: "+err.Error())
			}
		}
	}
	for _, e := range utm {
		gm.UTMAttribution = append(gm.UTMAttribution, *e)
	}

	token := os.Getenv("GUMROAD_ACCESS_TOKEN")
	mcpURL := os.Getenv("GUMROAD_MCP_URL")
	if mcpURL == "" {
		mcpURL = "https://gumroad.com/claude/v1/mcp"
	}
	if token != "" {
		views, vsrc, verr := fetchViewsMCP(mcpURL, token)
		if verr == nil {
			gm.Views = views
			gm.ViewsSource = vsrc
		} else {
			notes = append(notes, "gumroad mcp views failed: "+verr.Error())
		}
	} else {
		notes = append(notes, "GUMROAD_ACCESS_TOKEN not set; views unavailable")
	}

	rep.Gumroad = gm

	gc := GoatMetrics{}
	gcBase := os.Getenv("GOATCOUNTER_URL")
	gcToken := os.Getenv("GOATCOUNTER_TOKEN")
	if gcBase != "" {
		gc.Configured = true
		gc.BaseURL = gcBase
		v, vis, err := fetchGoatCounter(gcBase, gcToken)
		if err != nil {
			gc.Error = err.Error()
			notes = append(notes, "goatcounter fetch failed: "+err.Error())
		} else {
			gc.Verified = true
			gc.Views = v
			gc.Visitors = vis
		}
	} else {
		notes = append(notes, "GOATCOUNTER_URL not set; traffic unavailable")
	}
	rep.GoatCounter = gc
	rep.Notes = notes
	return rep
}

func fetchViewsMCP(url, token string) (int, string, error) {
	body := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]any{
			"name":      "get_analytics",
			"arguments": map[string]any{},
		},
	}
	b, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(b)))
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, "", fmt.Errorf("http %d", resp.StatusCode)
	}
	var raw map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return 0, "", err
	}
	text := ""
	if res, ok := raw["result"].(map[string]any); ok {
		if content, ok := res["content"].([]any); ok && len(content) > 0 {
			if c0, ok := content[0].(map[string]any); ok {
				if t, ok := c0["text"].(string); ok {
					text = t
				}
			}
		}
	}
	if text == "" {
		return 0, "mcp", fmt.Errorf("no text in mcp response")
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		return 0, "mcp", fmt.Errorf("unparseable analytics: %s", truncate(text, 200))
	}
	views := 0
	if v, ok := parsed["views"]; ok {
		views = toInt(v)
	}
	if rows, ok := parsed["data"].([]any); ok {
		for _, row := range rows {
			if m, ok := row.(map[string]any); ok {
				views += toInt(m["views"])
			}
		}
	}
	return views, "mcp:get_analytics", nil
}

func VerifyGoatCounter() error {
	base := os.Getenv("GOATCOUNTER_URL")
	token := os.Getenv("GOATCOUNTER_TOKEN")
	if base == "" {
		return fmt.Errorf("GOATCOUNTER_URL not set")
	}
	_, _, err := fetchGoatCounter(base, token)
	return err
}

func fetchGoatCounter(base, token string) (int, int, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	url := strings.TrimRight(base, "/") + "/api/v0/stats/total"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, 0, fmt.Errorf("http %d", resp.StatusCode)
	}
	var parsed map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return 0, 0, err
	}
	views := 0
	visitors := 0
	if v, ok := parsed["total"]; ok {
		views = toInt(v)
	}
	if v, ok := parsed["views"]; ok {
		views = toInt(v)
	}
	if v, ok := parsed["visitors"]; ok {
		visitors = toInt(v)
	}
	return views, visitors, nil
}

func utmSource(m map[string]any) string {
	for _, key := range []string{"utm_source", "utmSource", "source", "referrer"} {
		if v, ok := m[key]; ok {
			if str := toString(v); str != "" {
				return str
			}
		}
	}
	return "direct"
}

func saleRevenueCents(m map[string]any) float64 {
	for _, key := range []string{"price_cents", "amount_cents", "price"} {
		if v, ok := m[key]; ok {
			switch t := v.(type) {
			case float64:
				return t
			case int:
				return float64(t)
			case json.Number:
				f, _ := t.Float64()
				return f
			case string:
				var n float64
				fmt.Sscanf(t, "%f", &n)
				return n
			}
		}
	}
	return 0
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return fmt.Sprintf("%.0f", t)
	}
	return ""
}

func toInt(v any) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case json.Number:
		i, _ := t.Int64()
		return int(i)
	case string:
		var n int
		fmt.Sscanf(t, "%d", &n)
		return n
	}
	return 0
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
