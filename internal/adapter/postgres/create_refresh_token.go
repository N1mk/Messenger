package postgres

// СreateRefreshToken creates new element in the refresh_tokens table.
func (p *Postgres) CreateRefreshToken(tokenHash []byte, userID int) error {
	err := p.db.Create(&RefreshToken{TokenHash: tokenHash, UserID: userID}).Error

	return err
}
