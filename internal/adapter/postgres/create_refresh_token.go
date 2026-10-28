package postgres

// СreateRefreshToken creates new element in the refresh_tokens table.
func (p *Postgres) CreateRefreshToken(token string, userID int) error {
	err := p.db.Create(&RefreshToken{Token: token, UserID: userID}).Error

	return err
}
