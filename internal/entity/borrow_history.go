package entity

import "time"

type BorrowHistory struct {
	ID            int        `json:"id"`
	BookID        int        `json:"book_id"`
	UserID        int        `json:"user_id"`
	BorrowedAt    time.Time  `json:"borrowed_at"`
	BorrowedUntil time.Time  `json:"borrowed_until"`
	ReturnedAt    *time.Time `json:"returned_at"`
}

type OverdueBorrow struct {
	BorrowHistory
	BookTitle string `json:"book_title"`
}
