package config

import (
	"reflect"
	"testing"
)

func TestMustLoad(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			"config.MustLoad() test",
			args{"../../configs/main.yaml"},
			&Config{HTTPServer{"127.0.0.1:8080"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustLoad(tt.args.configPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustLoad() = %v, want %v", got, tt.want)
			}
		})
	}
}