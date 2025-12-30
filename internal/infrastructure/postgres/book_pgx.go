package postgres

import (
	"context"
	"lendbook/internal/db"
	"lendbook/internal/entity"
	"lendbook/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type bookRepository struct {
	queries *generated.Queries
}

func NewBookRepository(db *pgxpool.Pool) repository.BookRepository {
	return &bookRepository{
		queries: generated.New(db),
	}
}

func (b *bookRepository) Create(ctx context.Context, book *entity.Book) error {
	params := generated.CreateBookParams{
		Title:         book.Title,
		Author:        book.Author,
		PublishedDate: pgtype.Date{Time: book.PublishedDate, Valid: true},
		Language:      book.Language,
		AddedAt:       pgtype.Timestamp{Time: book.AddedAt, Valid: true},
		AddedBy:       int32(book.AddedBy),
	}

	id, err := b.queries.CreateBook(ctx, params)
	if err != nil {
		return err
	}

	book.ID = int(id)
	return nil
}

func (b *bookRepository) GetByID(ctx context.Context, id int) (*entity.Book, error) {
	row, err := b.queries.GetBookByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	book := &entity.Book{
		ID:            int(row.ID),
		Title:         row.Title,
		Author:        row.Author,
		PublishedDate: row.PublishedDate.Time,
		Language:      row.Language,
		AddedAt:       row.AddedAt.Time,
		AddedBy:       int(row.AddedBy),
	}

	if row.DeletedAt.Valid {
		t := row.DeletedAt.Time
		book.DeletedAt = &t
	}

	if row.DeletedBy != nil {
		val := int(*row.DeletedBy)
		book.DeletedBy = &val
	}

	return book, nil
}

func (b *bookRepository) Update(ctx context.Context, book *entity.Book) error {
	var deletedAt pgtype.Timestamp
	if book.DeletedAt != nil {
		deletedAt = pgtype.Timestamp{Time: *book.DeletedAt, Valid: true}
	} else {
		deletedAt = pgtype.Timestamp{Valid: false}
	}

	var deletedBy *int32
	if book.DeletedBy != nil {
		val := int32(*book.DeletedBy)
		deletedBy = &val
	}

	params := generated.UpdateBookParams{
		Title:         book.Title,
		Author:        book.Author,
		PublishedDate: pgtype.Date{Time: book.PublishedDate, Valid: true},
		Language:      book.Language,
		DeletedAt:     deletedAt,
		DeletedBy:     deletedBy,
		ID:            int32(book.ID),
	}

	return b.queries.UpdateBook(ctx, params)
}

func (b *bookRepository) ListAvailable(ctx context.Context) ([]entity.Book, error) {
	rows, err := b.queries.ListAvailableBooks(ctx)
	if err != nil {
		return nil, err
	}

	books := make([]entity.Book, 0, len(rows))
	for _, row := range rows {
		books = append(books, entity.Book{
			ID:            int(row.ID),
			Title:         row.Title,
			Author:        row.Author,
			PublishedDate: row.PublishedDate.Time,
			Language:      row.Language,
			AddedAt:       row.AddedAt.Time,
			AddedBy:       int(row.AddedBy),
		})
	}
	return books, nil
}

func (b *bookRepository) CreateBorrowHistory(ctx context.Context, borrow *entity.BorrowHistory) error {
	params := generated.CreateBorrowHistoryParams{
		BookID:     int32(borrow.BookID),
		UserID:     int32(borrow.UserID),
		BorrowedAt: pgtype.Timestamp{Time: borrow.BorrowedAt, Valid: true},
	}

	id, err := b.queries.CreateBorrowHistory(ctx, params)
	if err != nil {
		return err
	}

	borrow.ID = int(id)
	return nil
}

func (b *bookRepository) GetBorrowHistory(ctx context.Context, bookid int) ([]entity.BorrowHistory, error) {
	rows, err := b.queries.ListBorrowHistoryByBookID(ctx, int32(bookid))
	if err != nil {
		return nil, err
	}

	histories := make([]entity.BorrowHistory, 0, len(rows))
	for _, row := range rows {
		h := entity.BorrowHistory{
			ID:         int(row.ID),
			BookID:     int(row.BookID),
			UserID:     int(row.UserID),
			BorrowedAt: row.BorrowedAt.Time,
		}

		if row.ReturnedAt.Valid {
			t := row.ReturnedAt.Time
			h.ReturnedAt = &t
		}
		histories = append(histories, h)
	}
	return histories, nil
}

func (b *bookRepository) IsBookBorrowed(ctx context.Context, bookid int) (bool, error) {
	return b.queries.IsBookBorrowed(ctx, int32(bookid))
}

func (b *bookRepository) GetActiveBorrowHistory(ctx context.Context, userid, bookid int) (*entity.BorrowHistory, error) {
	params := generated.GetActiveBorrowHistoryParams{
		UserID: int32(userid),
		BookID: int32(bookid),
	}

	row, err := b.queries.GetActiveBorrowHistory(ctx, params)
	if err != nil {
		return nil, err
	}

	h := &entity.BorrowHistory{
		ID:         int(row.ID),
		BookID:     int(row.BookID),
		UserID:     int(row.UserID),
		BorrowedAt: row.BorrowedAt.Time,
	}

	if row.ReturnedAt.Valid {
		t := row.ReturnedAt.Time
		h.ReturnedAt = &t
	}

	return h, nil
}

func (b *bookRepository) UpdateBorrowHistory(ctx context.Context, borrow *entity.BorrowHistory) error {
	var returnedAt pgtype.Timestamp
	if borrow.ReturnedAt != nil {
		returnedAt = pgtype.Timestamp{Time: *borrow.ReturnedAt, Valid: true}
	} else {
		returnedAt = pgtype.Timestamp{Valid: false}
	}

	params := generated.UpdateBorrowHistoryReturnDateParams{
		ReturnedAt: returnedAt,
		ID:         int32(borrow.ID),
	}

	return b.queries.UpdateBorrowHistoryReturnDate(ctx, params)
}
