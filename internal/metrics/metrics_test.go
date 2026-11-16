package metrics

import "testing"

func TestToInt(t *testing.T) {
	cases := []struct {
		in   any
		want int
	}{
		{float64(42), 42},
		{int(7), 7},
		{"9", 9},
		{nil, 0},
	}
	for _, c := range cases {
		if got := toInt(c.in); got != c.want {
			t.Errorf("toInt(%v) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestToString(t *testing.T) {
	if toString("abc") != "abc" {
		t.Fatal("string")
	}
	if toString(float64(5)) != "5" {
		t.Fatal("float")
	}
	if toString(nil) != "" {
		t.Fatal("nil")
	}
}

func TestUTMSource(t *testing.T) {
	m := map[string]any{"utm_source": "hn"}
	if utmSource(m) != "hn" {
		t.Fatalf("got %s", utmSource(m))
	}
	m2 := map[string]any{"referrer": "direct"}
	if utmSource(m2) != "direct" {
		t.Fatalf("got %s", utmSource(m2))
	}
	m3 := map[string]any{}
	if utmSource(m3) != "direct" {
		t.Fatalf("got %s", utmSource(m3))
	}
}

func TestSaleRevenueCents(t *testing.T) {
	m := map[string]any{"price_cents": float64(2900)}
	if saleRevenueCents(m) != 2900 {
		t.Fatalf("got %v", saleRevenueCents(m))
	}
	m2 := map[string]any{"price": "19.5"}
	if saleRevenueCents(m2) != 19.5 {
		t.Fatalf("got %v", saleRevenueCents(m2))
	}
	if saleRevenueCents(map[string]any{}) != 0 {
		t.Fatal("empty")
	}
}

func TestTruncate(t *testing.T) {
	if truncate("hello", 10) != "hello" {
		t.Fatal("short")
	}
	if truncate("hello world", 5) != "hello..." {
		t.Fatal("long")
	}
}