package models

import "time"

type Post struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Image      string    `json:"image"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug" gorm:"unique"`
	Content    string    `json:"content"`
	CategoryID uint      `json:"category_id"`
	Category   Category    `json:"category" gorm:"foreignKey:CategoryID"`
	UserID     uint      `json:"user_id"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
