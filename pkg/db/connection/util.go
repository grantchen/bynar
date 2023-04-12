package sql_connection

import (
	"fmt"
	"strings"
)

func ChangeDatabaseConnectionSchema(connString, schema string) string {
	idx := strings.LastIndex(connString, "/")
	if idx == -1 {
		return fmt.Sprintf("%s/%s", connString, schema)
	}

	return fmt.Sprintf("%s/%s", connString[:idx], schema)
}
