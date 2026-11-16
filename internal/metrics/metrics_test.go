package metrics

import (
	"encoding/json"
	"testing"
)

func TestReportMarshal(t *testing.T) {
	r := Report{TS: "2026-01-01T00:00:00Z"}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("empty")
	}
}

func TestEnvelopeParse(t *testing.T) {
	raw := []byte(`{"metrics":{"views":5},"notes":["ok"]}`)
	var env struct {
		Metrics json.RawMessage `json:"metrics"`
		Notes   []string        `json:"notes"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatal(err)
	}
	if string(env.Metrics) != `{"views":5}` {
		t.Fatalf("got %s", string(env.Metrics))
	}
	if len(env.Notes) != 1 {
		t.Fatal("notes")
	}
}
