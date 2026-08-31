package devpanel

import "testing"

func TestValidateServiceName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "api"},
		{name: "url-shortener"},
		{name: "worker_2"},
		{name: "../api", wantErr: true},
		{name: "services/api", wantErr: true},
		{name: "api.log", wantErr: true},
		{name: "api\nforged", wantErr: true},
		{name: "", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := validateServiceName(test.name)
			if test.wantErr && err == nil {
				t.Fatal("expected service name to be rejected")
			}
			if !test.wantErr && err != nil {
				t.Fatalf("expected service name to be accepted: %v", err)
			}
		})
	}
}
