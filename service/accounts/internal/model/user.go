/**
    @author: dongjs
    @date: 2023/9/13
    @description:
**/

package model

// User users table of organization_schema(uuid) struct
type User struct {
	ID                 int
	Email              string
	FullName           string
	Phone              string
	Status             bool
	LanguagePreference string
	PolicyId           int
	Theme              string
}
