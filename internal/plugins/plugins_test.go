package plugins

import "testing"

func TestFind(t *testing.T) {
	m := Manifest{Plugins: []Plugin{{Name: "a"}, {Name: "b"}}}
	p, idx := Find(m, "b")
	if idx != 1 || p.Name != "b" {
		t.Fatalf("got %d %v", idx, p)
	}
	_, idx2 := Find(m, "zzz")
	if idx2 != -1 {
		t.Fatal("expected -1")
	}
}
