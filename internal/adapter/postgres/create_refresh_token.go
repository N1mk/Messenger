package postgres

import "context"

// СreateRefreshToken creates new element in the refresh_tokens table.
func (p *Postgres) CreateRefreshToken(ctx context.Context, tokenHash []byte, userID int) error {
	err := p.db.WithContext(ctx).Create(&RefreshToken{TokenHash: tokenHash, UserID: userID}).Error

	return err
}
