package sql_connection

import "testing"

func TestChangeDatabaseConnectionSchema(t *testing.T) {
	type args struct {
		connString string
		schema     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "with suffix",
			args: args{
				connString: "connection_string/default_schema",
				schema:     "another_schema",
			},
			want: "connection_string/another_schema",
		},
		{
			name: "with slash only",
			args: args{
				connString: "connection_string/",
				schema:     "another_schema",
			},
			want: "connection_string/another_schema",
		},
		{
			name: "without suffix",
			args: args{
				connString: "connection_string",
				schema:     "another_schema",
			},
			want: "connection_string/another_schema",
		},
		{
			name: "empty string",
			args: args{
				connString: "",
				schema:     "another_schema",
			},
			want: "/another_schema",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ChangeDatabaseConnectionSchema(tt.args.connString, tt.args.schema); got != tt.want {
				t.Errorf("ChangeDatabaseConnectionSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}
