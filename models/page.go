package models

import "time"

type Page struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug" gorm:"unique"`
	Content   string    `json:"content"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
