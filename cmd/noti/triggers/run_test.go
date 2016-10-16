package triggers

import "testing"

func TestKeyValue(t *testing.T) {
	cases := []struct {
		in      string
		wantKey string
		wantVal string
	}{
		{"exit", "exit", ""},
		{"contains=hello world", "contains", "hello world"},
		{"interval=10s", "interval", "10s"},
	}

	for i, c := range cases {
		k, v := keyValue(c.in)

		if k != c.wantKey {
			t.Error("Unexpected key")
			t.Errorf("%d: got: %q; want: %q", i, k, c.wantKey)
		}

		if v != c.wantVal {
			t.Error("Unexpected value")
			t.Errorf("%d: got: %q; want: %q", i, v, c.wantVal)
		}
	}
}
