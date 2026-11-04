package repository

import "gorm.io/gorm"

const defaultPagesize = 20

func Paginate(lastID, pageSize int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageSize <= 0 {
			pageSize = defaultPagesize
		}

		if lastID < 0 {
			lastID = 0
		}

		return db.Where("id < ?", lastID).Order("id DESC").Limit(pageSize)
	}
}
