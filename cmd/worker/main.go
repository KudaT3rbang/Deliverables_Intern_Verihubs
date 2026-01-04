package main

import (
	"context"
	"log"
	"time"

	"lendbook/internal/infrastructure/postgres"
	"lendbook/internal/worker/jobs"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

func main() {
	if err := godotenv.Load(". env"); err != nil {
		log.Println("Error loading . env file")
	}

	ctx := context.Background()

	dbPool, err := postgres.InitDb()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer dbPool.Close()

	log.Println("Worker connecting to database...")

	bookRepo := postgres.NewBookRepository(dbPool)
	notificationRepo := postgres.NewNotificationRepository(dbPool)

	workers := river.NewWorkers()
	river.AddWorker(workers, jobs.NewSendNotificationWorker(notificationRepo))

	driver := riverpgxv5.New(dbPool)
	riverClient, err := river.NewClient(driver, &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
		Workers: workers,
	})
	if err != nil {
		log.Fatalf("Failed to create river client: %v", err)
	}

	river.AddWorker(workers, jobs.NewCheckOverdueWorker(bookRepo, riverClient))

	if err := riverClient.Start(ctx); err != nil {
		log.Fatalf("Failed to start river client: %v", err)
	}

	log.Println("Worker started successfully")

	scheduleDailyCheck(ctx, riverClient)
}

func scheduleDailyCheck(ctx context.Context, client *river.Client[pgx.Tx]) {
	log.Println("Inserting initial overdue check job...")
	if _, err := client.Insert(ctx, jobs.CheckOverdueArgs{}, nil); err != nil {
		log.Printf("Failed to insert initial check overdue job: %v", err)
	}

	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		log.Println("Inserting daily overdue check job...")
		if _, err := client.Insert(ctx, jobs.CheckOverdueArgs{}, nil); err != nil {
			log.Printf("Failed to insert check overdue job: %v", err)
		}
	}
}
