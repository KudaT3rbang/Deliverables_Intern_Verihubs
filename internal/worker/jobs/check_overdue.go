package jobs

import (
	"context"
	"fmt"
	"lendbook/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

type CheckOverdueArgs struct{}

func (CheckOverdueArgs) Kind() string { return "check_overdue" }

type CheckOverdueWorker struct {
	river.WorkerDefaults[CheckOverdueArgs]
	bookRepo    repository.BookRepository
	riverClient *river.Client[pgx.Tx]
}

func NewCheckOverdueWorker(
	bookRepo repository.BookRepository,
	riverClient *river.Client[pgx.Tx],
) *CheckOverdueWorker {
	return &CheckOverdueWorker{
		bookRepo:    bookRepo,
		riverClient: riverClient,
	}
}

func (w *CheckOverdueWorker) Work(ctx context.Context, job *river.Job[CheckOverdueArgs]) error {
	overdueBorrows, err := w.bookRepo.GetOverdueBorrows(ctx)
	if err != nil {
		return fmt.Errorf("failed to get overdue borrows: %w", err)
	}

	for _, borrow := range overdueBorrows {
		message := fmt.Sprintf(
			"Your borrowed book '%s' was due on %s. Please return it as soon as possible.",
			borrow.BookTitle,
			borrow.BorrowedUntil.Format("2006-01-02"),
		)

		_, err := w.riverClient.Insert(ctx, SendNotificationArgs{
			UserID:  borrow.UserID,
			Message: message,
		}, nil)
		if err != nil {
			return fmt.Errorf("failed to insert notification job for user %d: %w", borrow.UserID, err)
		}
	}

	return nil
}
