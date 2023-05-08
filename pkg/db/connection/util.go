package sql_connection

import (
	"encoding/json"
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

func JSON2DatabaseConnection(jsonStr string) string {
	sec := map[string]interface{}{}
	err := json.Unmarshal([]byte(jsonStr), &sec)
	if err != nil {
		panic(err)
	}
	// root:123456@tcp(localhost:3306)/bynar
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		sec["username"],
		sec["password"],
		sec["host"],
		sec["port"],
		sec["dbname"],
	)
}
