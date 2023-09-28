/**
    @author: dongjs
    @date: 2023/9/13
    @description:
**/

package organization_schema

// User users table of organization_schema(uuid) struct
type User struct {
	ID                 int    `db:"id"`
	Email              string `db:"email"`
	FullName           string `db:"full_name"`
	Phone              string `db:"phone"`
	Status             bool   `db:"status"`
	LanguagePreference string `db:"language_preference"`
	Theme              string `db:"theme"`
	ProfilePhoto       string `db:"profile_photo"`
	Policies           string `db:"policies"`
}
