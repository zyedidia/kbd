package kbd

import "testing"

func TestExpand(t *testing.T) {
	tests := []struct {
		template string
		args     []string
		expect   string
	}{
		{"$0 $1", []string{"hello", "world"}, "hello world"},
		{"$$1", []string{}, "$1"},
		{"$foo $0", []string{"bar"}, "$foo bar"},
		{"$-2", []string{}, "$-2"},
	}

	for _, tt := range tests {
		t.Run(tt.template, func(t *testing.T) {
			n := nargs(tt.template)
			if n != len(tt.args) {
				t.Fatalf("nargs: %d, expected: %d", n, len(tt.args))
			}
			expanded := expand(tt.template, tt.args)
			if expanded != tt.expect {
				t.Fatalf("expand: %s, expected: %s", expanded, tt.expect)
			}
		})
	}
}
