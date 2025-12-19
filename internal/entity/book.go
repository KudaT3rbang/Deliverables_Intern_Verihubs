package entity

import "time"

type Book struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Author        string     `json:"author"`
	PublishedDate time.Time  `json:"published_date"`
	Language      string     `json:"language"`
	AddedAt       time.Time  `json:"added_at"`
	AddedBy       int        `json:"added_by"`
	DeletedAt     *time.Time `json:"deleted_at"`
	DeletedBy     *int       `json:"deleted_by"`
}
