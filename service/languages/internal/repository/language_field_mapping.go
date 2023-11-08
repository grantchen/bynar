package repository

var LanguageFieldNames = map[string][]string{
	"id":            {"id"},
	"country":       {"country"},
	"language":      {"language"},
	"two_letters":   {"two_letters"},
	"three_letters": {"three_letters"},
	"number":        {"number"},
}

var LanguageFieldCountry = map[string][]string{
	"country": {"country"},
}
var LanguageFieldLanguage = map[string][]string{
	"language": {"language"},
}
var LanguageFieldLetters = map[string][]string{
	"two_letters":   {"two_letters"},
	"three_letters": {"three_letters"},
}
var LanguageFieldNamesFloat = map[string][]string{
	"number": {"number"},
}
