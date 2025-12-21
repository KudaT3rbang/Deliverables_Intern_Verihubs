package usecase

import (
	"context"
	"errors"
	"lendbook/internal/entity"
	"lendbook/internal/repository"
	"time"
)

type BookUsecase interface {
	AddBook(ctx context.Context, userid int, params entity.AddBookParams) error
	DeleteBook(ctx context.Context, userid int, bookid int) error
	ListBooks(ctx context.Context) ([]entity.Book, error)
	GetBookDetails(ctx context.Context, bookid int) (*entity.BookDetailResponse, error)
	BorrowBook(ctx context.Context, userid int, bookid int) error
	ReturnBook(ctx context.Context, userid int, bookid int) error
}

type bookUsecase struct {
	repo repository.BookRepository
}

func NewBookUsecase(repo repository.BookRepository) BookUsecase {
	return &bookUsecase{
		repo: repo,
	}
}

func (b *bookUsecase) AddBook(ctx context.Context, userid int, params entity.AddBookParams) error {
	newBook := entity.Book{
		Title:         params.Title,
		Author:        params.Author,
		PublishedDate: params.PublishedDate,
		Language:      params.Language,
		AddedBy:       userid,
		AddedAt:       time.Now(),
	}

	return b.repo.Create(ctx, &newBook)
}

func (b *bookUsecase) DeleteBook(ctx context.Context, userid int, bookid int) error {
	deleteBook, err := b.repo.GetByID(ctx, bookid)
	if err != nil {
		return errors.New("book not found")
	}
	if deleteBook.AddedBy != userid {
		return errors.New("user not authorized to delete book")
	}

	timeDeleted := time.Now()
	deleteBook.DeletedBy = &userid
	deleteBook.DeletedAt = &timeDeleted
	return b.repo.Update(ctx, deleteBook)
}

func (b *bookUsecase) ListBooks(ctx context.Context) ([]entity.Book, error) {
	return b.repo.ListAvailable(ctx)
}

func (b *bookUsecase) GetBookDetails(ctx context.Context, bookid int) (*entity.BookDetailResponse, error) {
	book, err := b.repo.GetByID(ctx, bookid)
	if err != nil || book.DeletedAt != nil {
		return nil, errors.New("book not available")
	}

	bookRecord, err := b.repo.GetBorrowHistory(ctx, bookid)
	if err != nil {
		return nil, errors.New("book not available")
	}

	bookDetailResponse := entity.BookDetailResponse{
		Book:          *book,
		BorrowHistory: bookRecord,
	}

	return &bookDetailResponse, nil
}

func (b *bookUsecase) BorrowBook(ctx context.Context, userid int, bookid int) error {
	borrowBook, err := b.repo.GetByID(ctx, bookid)
	if err != nil || borrowBook.DeletedAt != nil {
		return errors.New("book not available")
	}
	isBorrowed, err := b.repo.IsBookBorrowed(ctx, bookid)
	if err != nil {
		return errors.New("book not available")
	}
	if isBorrowed {
		return errors.New("book is borrowed by others")
	}
	borrow := entity.BorrowHistory{
		BookID:     bookid,
		UserID:     userid,
		BorrowedAt: time.Now(),
		ReturnedAt: nil,
	}

	return b.repo.CreateBorrowHistory(ctx, &borrow)
}

func (b *bookUsecase) ReturnBook(ctx context.Context, userid int, bookid int) error {
	borrowedBook, err := b.repo.GetActiveBorrowHistory(ctx, userid, bookid)
	if err != nil {
		return errors.New("you do not have this book borrowed")
	}

	now := time.Now()
	borrowedBook.ReturnedAt = &now

	return b.repo.UpdateBorrowHistory(ctx, borrowedBook)
}
