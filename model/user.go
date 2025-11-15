package models

import "time"

type User struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	Username     string `gorm:"uniqueIndex;size:50;not null"`
	Email        string `gorm:"uniqueIndex;size:100;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Posts        []Post
	Comments     []Comment
}
