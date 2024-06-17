package hostmatcher

import (
	"testing"
)

func Test_hostMatcher_Match(t *testing.T) {

	tests := []struct {
		name  string
		rules []string
		match string
		want  bool
	}{
		{
			rules: []string{
				"*/hello",
			},
			match: "local.host/hello",
			want:  true,
		},
		{
			rules: []string{
				"local.host/*",
			},
			match: "local.host/hello",
			want:  true,
		},
		{
			rules: []string{
				"local.host/*",
			},
			match: "local.host/hello/world",
			want:  false,
		},
		{
			rules: []string{
				"local.host/**",
			},
			match: "local.host/hello/world",
			want:  true,
		},
		{
			rules: []string{
				"*.host",
			},
			match: "local.host",
			want:  true,
		},
		{
			rules: []string{
				"*.host",
			},
			match: "local.nohost",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatcher(tt.rules)
			if got := m.Match(tt.match); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
