/**
    @author: dongjs
    @date: 2023/9/12
    @description:
**/

package repository

// UpdateProfilePhotoOfUsers update column profile_photo in table users
func (r *accountRepositoryHandler) UpdateProfilePhotoOfUsers(accountID int, profilePhoto string) error {
	if _, err := r.db.Exec(`UPDATE users SET profile_photo = ? WHERE id = ?`, profilePhoto, accountID); err != nil {
		return err
	}
	return nil
}
