package initializers

import "testing"

func TestConnectToDB(t *testing.T) {
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
