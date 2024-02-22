package config

import (
	"reflect"
	"testing"

	"github.com/lunn06/video-hosting/internal/models"
)

func TestMustLoad(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name string
		args args
		want Config
	}{
		{
			"config.MustLoad() test",
			args{"../../configs/example_main.yaml"},
			Config{
				HTTPServer{"127.0.0.1", "8080"},
				Database{
					"pgsql.com:5432",
					"db_user",
					"db_name",
					"db_password",
					"disable",
				},
			},
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

func TestMustLoadDatabaseDefaults(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name string
		args args
		want DatabaseDefaults
	}{
		{
			"MustLoadDatabaseDefaults() Test",
			args{"../../configs/database_defaults.yaml"},
			DatabaseDefaults{
				[]models.Role{
					{1, "user", false, false},
					{2, "moderator", false, true},
					{3, "admin", true, true},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustLoadDatabaseDefaults(tt.args.configPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustLoadDatabaseDefaults() = %v, want %v", got, tt.want)
			}
		})
	}
}
