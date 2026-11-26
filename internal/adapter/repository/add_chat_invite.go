package repository

import "messenger/internal/domain"

// AddChatInvite adds new chat invite (domain.ChatInvite) into the storage.
func (pr *PostgresRepo) AddChatInvite(invite domain.ChatInvite) error {
	return pr.db.Create(&invite).Error
}
