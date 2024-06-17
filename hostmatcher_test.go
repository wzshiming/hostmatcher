package hostmatcher

import (
	"net"
	"reflect"
	"testing"
)

func Test_newHostMatcher(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		args    args
		want    matcher
		wantErr bool
	}{
		{
			args: args{
				host: "127.0.0.1",
			},
			want: ipMatch{
				ip: net.IPv4(127, 0, 0, 1),
			},
		},
		{
			args: args{
				host: "127.0.0.1:80",
			},
			want: ipMatch{
				ip:   net.IPv4(127, 0, 0, 1),
				port: "80",
			},
		},
		{
			args: args{
				host: "10.0.0.1/8",
			},
			want: cidrMatch{
				cidr: parseCIDR("10.0.0.1/8"),
			},
		},
		{
			args: args{
				host: "localhost",
			},
			want: domainMatch{
				host:      ".localhost",
				matchHost: true,
			},
		},
		{
			args: args{
				host: "localhost:80",
			},
			want: domainMatch{
				host:      ".localhost",
				port:      "80",
				matchHost: true,
			},
		},
		{
			args: args{
				host: "local.host",
			},
			want: domainMatch{
				host:      ".local.host",
				matchHost: true,
			},
		},
		{
			args: args{
				host: ".local.host",
			},
			want: domainMatch{
				host:      ".local.host",
				matchHost: false,
			},
		},
		{
			args: args{
				host: "*.local.host",
			},
			want: domainMatch{
				host:      ".local.host",
				matchHost: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newHostMatcher(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("newHostMatcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newHostMatcher() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func parseCIDR(host string) *net.IPNet {
	_, pnet, _ := net.ParseCIDR(host)
	return pnet
}
