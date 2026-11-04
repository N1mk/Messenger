package postgres

// СreateRefreshToken creates new element in the refresh_tokens table.
func (r *PostgresRepo) CreateRefreshToken(tokenHash []byte, userID int) error {
	err := r.db.Create(&RefreshToken{TokenHash: tokenHash, UserID: userID}).Error

	return err
}
