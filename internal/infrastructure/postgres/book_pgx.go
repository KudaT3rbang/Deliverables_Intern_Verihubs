package postgres

import (
	"context"
	"lendbook/internal/entity"
	"lendbook/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type bookRepository struct {
	db *pgxpool.Pool
}

func NewBookRepository(db *pgxpool.Pool) repository.BookRepository {
	return &bookRepository{
		db: db,
	}
}

func (b *bookRepository) Create(ctx context.Context, book *entity.Book) error {
	query := `INSERT INTO books (title, author, published_date, language, added_at, added_by)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`
	err := b.db.QueryRow(ctx, query,
		book.Title,
		book.Author,
		book.PublishedDate,
		book.Language,
		book.AddedAt,
		book.AddedBy).Scan(&book.ID)
	if err != nil {
		return err
	}
	return nil
}

func (b *bookRepository) GetByID(ctx context.Context, id int) (*entity.Book, error) {
	query := `SELECT id, title, author, published_date, language, added_at, added_by, deleted_at, deleted_by FROM books WHERE id = $1`
	var book entity.Book
	err := b.db.QueryRow(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.PublishedDate,
		&book.Language,
		&book.AddedAt,
		&book.AddedBy,
		&book.DeletedAt,
		&book.DeletedBy,
	)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (b *bookRepository) Update(ctx context.Context, book *entity.Book) error {
	query := `UPDATE books SET title=$1, author=$2, published_date=$3, language=$4, deleted_at=$5, deleted_by=$6 WHERE id=$7`
	_, err := b.db.Exec(ctx, query,
		book.Title,
		book.Author,
		book.PublishedDate,
		book.Language,
		book.DeletedAt,
		book.DeletedBy,
		book.ID)
	if err != nil {
		return err
	}
	return nil
}

func (b *bookRepository) ListAvailable(ctx context.Context) ([]entity.Book, error) {
	query := `SELECT id, title, author, published_date, language, added_at, added_by FROM books WHERE deleted_at IS NULL`
	rows, err := b.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var books []entity.Book
	for rows.Next() {
		var book entity.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedDate, &book.Language, &book.AddedAt, &book.AddedBy); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (b *bookRepository) CreateBorrowHistory(ctx context.Context, borrow *entity.BorrowHistory) error {
	query := `INSERT INTO borrow_history (book_id, user_id, borrowed_at) VALUES ($1, $2, $3) RETURNING id`
	return b.db.QueryRow(ctx, query, borrow.BookID, borrow.UserID, borrow.BorrowedAt).Scan(&borrow.ID)
}

func (b *bookRepository) GetBorrowHistory(ctx context.Context, bookid int) ([]entity.BorrowHistory, error) {
	query := `SELECT id, book_id, user_id, borrowed_at, returned_at FROM borrow_history WHERE book_id = $1 ORDER BY borrowed_at DESC`
	rows, err := b.db.Query(ctx, query, bookid)
	if err != nil {
		return nil, err
	}

	var histories []entity.BorrowHistory
	for rows.Next() {
		var h entity.BorrowHistory
		if err := rows.Scan(&h.ID, &h.BookID, &h.UserID, &h.BorrowedAt, &h.ReturnedAt); err != nil {
			return nil, err
		}
		histories = append(histories, h)
	}
	return histories, nil
}

func (b *bookRepository) IsBookBorrowed(ctx context.Context, bookid int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM borrow_history WHERE book_id = $1 AND returned_at IS NULL)`
	var exists bool
	err := b.db.QueryRow(ctx, query, bookid).Scan(&exists)
	return exists, err
}

func (b *bookRepository) GetActiveBorrowHistory(ctx context.Context, userid, bookid int) (*entity.BorrowHistory, error) {
	query := `SELECT id, book_id, user_id, borrowed_at, returned_at FROM borrow_history WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL`
	var borrowHistory entity.BorrowHistory
	err := b.db.QueryRow(ctx, query, userid, bookid).Scan(
		&borrowHistory.ID,
		&borrowHistory.BookID,
		&borrowHistory.UserID,
		&borrowHistory.BorrowedAt,
		&borrowHistory.ReturnedAt)
	if err != nil {
		return nil, err
	}
	return &borrowHistory, nil
}

func (b *bookRepository) UpdateBorrowHistory(ctx context.Context, borrow *entity.BorrowHistory) error {
	query := `UPDATE borrow_history SET returned_at = $1 WHERE id = $2`
	_, err := b.db.Exec(ctx, query, borrow.ReturnedAt, borrow.ID)
	return err
}
