package hostmatcher

import (
	"testing"
)

func Test_pathPatternMatcher_match(t *testing.T) {
	type fields struct {
		pattern string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			fields: fields{
				pattern: "*",
			},
			args: args{
				"a",
			},
			want: true,
		},
		{
			fields: fields{
				pattern: "*/a",
			},
			args: args{
				"b/a",
			},
			want: true,
		},
		{
			fields: fields{
				pattern: "b/*",
			},
			args: args{
				"b/a",
			},
			want: true,
		},
		{
			fields: fields{
				pattern: "*/a",
			},
			args: args{
				"a",
			},
			want: false,
		},
		{
			fields: fields{
				pattern: "a/*",
			},
			args: args{
				"a",
			},
			want: false,
		},

		{
			fields: fields{
				pattern: "**",
			},
			args: args{
				"a",
			},
			want: true,
		},
		{
			fields: fields{
				pattern: "**/a",
			},
			args: args{
				"b/a",
			},
			want: true,
		},
		{
			fields: fields{
				pattern: "b/**",
			},
			args: args{
				"b/a",
			},
			want: true,
		},
		{
			fields: fields{
				pattern: "**/a",
			},
			args: args{
				"a",
			},
			want: false,
		},
		{
			fields: fields{
				pattern: "a/**",
			},
			args: args{
				"a",
			},
			want: false,
		},

		{
			fields: fields{
				pattern: "**",
			},
			args: args{
				"a/b",
			},
			want: true,
		},
		{
			fields: fields{
				pattern: "**/a",
			},
			args: args{
				"c/b/a",
			},
			want: true,
		},
		{
			fields: fields{
				pattern: "b/**",
			},
			args: args{
				"b/a/c",
			},
			want: true,
		},
		{
			fields: fields{
				pattern: "**/a",
			},
			args: args{
				"a",
			},
			want: false,
		},
		{
			fields: fields{
				pattern: "a/**",
			},
			args: args{
				"a",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pathPatternMatcher{
				pattern: tt.fields.pattern,
			}
			if got := p.match(tt.args.path); got != tt.want {
				t.Errorf("match() = %v, want %v", got, tt.want)
			}
		})
	}
}
