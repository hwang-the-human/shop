package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(connectionString string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Printf("Error opening database: %v\n", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Printf("Error pinging database: %v\n", err)
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return &PostgresRepository{
		DB: db,
	}, nil
}

func (r *PostgresRepository) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	log.Printf("Executing query: %s with args: %v\n", query, args)
	result, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
	}
	return result, err
}

func (r *PostgresRepository) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	log.Printf("Running query: %s with args: %v\n", query, args)
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error running query: %v\n", err)
	}
	return rows, err
}

func (r *PostgresRepository) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	log.Printf("Querying row: %s with args: %v\n", query, args)
	return r.DB.QueryRowContext(ctx, query, args...)
}

func (r *PostgresRepository) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error beginning transaction: %v\n", err)
		return err
	}

	if err := fn(tx); err != nil {
		log.Printf("Error in transaction function: %v\n", err)
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("Transaction rollback error: %v, original error: %v\n", rbErr, err)
			return fmt.Errorf("transaction rollback error: %v, original error: %v", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		return err
	}

	log.Println("Transaction committed successfully")
	return nil
}

func (r *PostgresRepository) Close() error {
	log.Println("Closing database connection")
	return r.DB.Close()
}
