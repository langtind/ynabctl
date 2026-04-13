package period

import "testing"

func TestParseSpecific(t *testing.T) {
	cases := []struct {
		kind, value, name, start, end string
	}{
		{"month", "2026-03", "2026-03", "2026-03-01", "2026-03-31"},
		{"quarter", "2026-Q1", "2026-Q1", "2026-01-01", "2026-03-31"},
		{"quarter", "2026-Q4", "2026-Q4", "2026-10-01", "2026-12-31"},
		{"year", "2026", "2026", "2026-01-01", "2026-12-31"},
		{"week", "2026-W15", "2026-W15", "2026-04-06", "2026-04-12"},
	}
	for _, c := range cases {
		got, err := Compute(c.kind, c.value)
		if err != nil {
			t.Errorf("%s %s: %v", c.kind, c.value, err)
			continue
		}
		if got.Name != c.name || got.StartDate != c.start || got.EndDate != c.end {
			t.Errorf("%s %s: got %+v, want name=%s start=%s end=%s", c.kind, c.value, got, c.name, c.start, c.end)
		}
	}
}

func TestInvalid(t *testing.T) {
	if _, err := Compute("decade", ""); err == nil {
		t.Error("expected error for unknown kind")
	}
	if _, err := Compute("month", "2026"); err == nil {
		t.Error("expected error for bad month format")
	}
}
