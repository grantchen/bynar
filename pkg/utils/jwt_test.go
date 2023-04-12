package utils

import (
	"reflect"
	"testing"
)

func TestExtractClaimsFromJWT(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		wantClaims map[string]interface{}
		wantErr    bool
	}{
		{
			name:  "valid JWT payload",
			token: "xxxxx.eyJrZXkiOiAidmFsdWUiLCAiYm9vbCI6IHRydWV9.xxxxx",
			wantClaims: map[string]interface{}{
				"key":  "value",
				"bool": true,
			},
		},
		{
			name:  "valid Bearer JWT payload",
			token: "Bearer xxxxx.eyJrZXkiOiAidmFsdWUiLCAiYm9vbCI6IHRydWV9.xxxxx",
			wantClaims: map[string]interface{}{
				"key":  "value",
				"bool": true,
			},
		},
		{
			name:       "invalid token",
			token:      "invalid_token",
			wantClaims: nil,
			wantErr:    true,
		},
		{
			name:       "invalid payload",
			token:      "xxxxx.xxx.xxxxx",
			wantClaims: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClaims, err := ExtractClaimsFromJWT(tt.token)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ExtractClaimsFromJWT() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotClaims, tt.wantClaims) {
				t.Errorf("ExtractClaimsFromJWT() gotClaims = %v, want %v", gotClaims, tt.wantClaims)
			}
		})
	}
}
