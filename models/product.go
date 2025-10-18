package models

import "time"

type Product struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Image     string    `json:"image"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug" gorm:"unique"`
	Content   string    `json:"content"`
	Owner     string    `json:"owner"`
	Price     int       `json:"price"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
