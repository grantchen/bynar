package scope

import (
	"reflect"
	"testing"
)

func TestResolveFromToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		want    RequestScope
		wantErr bool
	}{
		{
			name:  "Token with string values",
			token: "Bearer xxxxxx.eyJjdXN0b206YWNjb3VudF9pZCI6ICIzOSIsICJjdXN0b206b3JnYW5pemF0aW9uX2lkIjogIjQ1In0.xxxx",
			want: RequestScope{
				AccountID:      39,
				OrganizationID: 45,
			},
		},
		{
			name:  "Token with int values", // this case shouldn't happen if AWS Cognito properties are properly configured.
			token: "Bearer xxxxxx.eyJjdXN0b206YWNjb3VudF9pZCI6IDM5LCAiY3VzdG9tOm9yZ2FuaXphdGlvbl9pZCI6IDQ1fQ.xxxx",
			want: RequestScope{
				AccountID:      39,
				OrganizationID: 45,
			},
		},
		{
			name:  "Token with string values",
			token: "xxxxxx.eyJjdXN0b206YWNjb3VudF9pZCI6ICIzOSIsICJjdXN0b206b3JnYW5pemF0aW9uX2lkIjogIjQ1In0.xxxx",
			want: RequestScope{
				AccountID:      39,
				OrganizationID: 45,
			},
		},
		{
			name:  "Token with int values", // this case shouldn't happen if AWS Cognito properties are properly configured.
			token: "xxxxxx.eyJjdXN0b206YWNjb3VudF9pZCI6IDM5LCAiY3VzdG9tOm9yZ2FuaXphdGlvbl9pZCI6IDQ1fQ.xxxx",
			want: RequestScope{
				AccountID:      39,
				OrganizationID: 45,
			},
		},
		{
			name:    "Token with invalid account id",
			token:   "xxxxxx.eyJjdXN0b206YWNjb3VudF9pZCI6ICJpbnZhbGlkIiwgImN1c3RvbTpvcmdhbml6YXRpb25faWQiOiA0NX0.xxxx",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveFromToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveFromToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResolveFromToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
