package models

import "time"

type Comment struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	PostID    uint64 `gorm:"not null"`
	UserID    uint64 `gorm:"not null"`
	Content   string `gorm:"type:text;not null"`
	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;"`
}
