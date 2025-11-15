package models

import "time"

type Post struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	UserID    uint64 `gorm:"not null"`
	Title     string `gorm:"size:200;not null"`
	Content   string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	User     User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Comments []Comment
}
