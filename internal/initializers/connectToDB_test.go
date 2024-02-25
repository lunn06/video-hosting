package initializers

import (
	"testing"

	"github.com/lunn06/video-hosting/internal/config"
)

func TestConnectToDB(t *testing.T) {
	config.CFG = config.MustLoad("../../configs/main.yaml")
	tests := []struct {
		name string
	}{
		{name: "ConnectToDB Test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ConnectToDB()
		})
	}
}
