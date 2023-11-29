package repository

// UserFieldNames is a map of user field names to their respective database column names.
var UserFieldNames = map[string][]string{
	"id":                  {"id"},
	"email":               {"email"},
	"full_name":           {"full_name"},
	"phone":               {"phone"},
	"status":              {"status"},
	"language_preference": {"language_preference"},
	"theme":               {"theme"},
	"policies":            {"policies"},
}

// UserFieldNamesString is a map of user field names to their respective database column names.
var UserFieldNamesString = map[string][]string{
	"email":               {"email"},
	"full_name":           {"full_name"},
	"phone":               {"phone"},
	"language_preference": {"language_preference"},
	"theme":               {"theme"},
}
