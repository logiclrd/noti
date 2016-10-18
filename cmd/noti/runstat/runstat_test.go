package runstat

import "testing"

func TestParseExpansion(t *testing.T) {
	cases := []struct {
		s     string
		alias string
		want  string
	}{
		{"\x1b]7;file://foo.local/Users/bar/Developer/go/src/localhost/tmp\agss: aliased to git status --short",
			"gss",
			"git status --short"},
		{"gss: aliased to git status --short",
			"gss",
			"git status --short"},
		{"",
			"gss",
			""},
	}

	for i, c := range cases {
		got := parseExpansion(c.s, c.alias)
		if got != c.want {
			t.Error("Unexpected expansion")
			t.Errorf("%d: got: %q; want: %q", i, got, c.want)
		}
	}
}
