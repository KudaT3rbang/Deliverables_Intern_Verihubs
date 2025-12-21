package repository

import (
	"context"
	"lendbook/internal/entity"
)

type BookRepository interface {
	Create(ctx context.Context, book *entity.Book) error
	GetByID(ctx context.Context, id int) (*entity.Book, error)
	Update(ctx context.Context, book *entity.Book) error
	ListAvailable(ctx context.Context) ([]entity.Book, error)
	CreateBorrowHistory(ctx context.Context, borrow *entity.BorrowHistory) error
	GetBorrowHistory(ctx context.Context, bookid int) ([]entity.BorrowHistory, error)
	IsBookBorrowed(ctx context.Context, bookid int) (bool, error)
	GetActiveBorrowHistory(ctx context.Context, userid, bookid int) (*entity.BorrowHistory, error)
	UpdateBorrowHistory(ctx context.Context, borrow *entity.BorrowHistory) error
}
