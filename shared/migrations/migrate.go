package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Migration struct {
	ID            int
	MigrationName string
	CreatedAt     time.Time
}

func ApplyMigrations(db *sql.DB, migrationsDir string) error {
	if err := ensureMigrationsTable(db); err != nil {
		return err
	}

	log.Println("Migrations table ensured.")

	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	log.Println("Fetched applied migrations.")

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" && !strings.Contains(file.Name(), "_rollback") {
			migrationName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			if !isMigrationApplied(appliedMigrations, migrationName) {
				log.Printf("Applying migration: %s\n", migrationName)
				if err := applyMigration(db, filepath.Join(migrationsDir, file.Name()), migrationName); err != nil {
					return err
				}
				log.Printf("Applied migration: %s\n", migrationName)
			} else {
				log.Printf("Skipping already applied migration: %s\n", migrationName)
			}
		}
	}

	log.Println("Migrations applied successfully")

	return nil
}

func RollbackMigration(db *sql.DB, migrationsDir string, migrationName string) error {
	if err := ensureMigrationsTable(db); err != nil {
		return err
	}

	log.Println("Migrations table ensured.")

	rollbackFile := filepath.Join(migrationsDir, migrationName+"_rollback.sql")
	if _, err := os.Stat(rollbackFile); err == nil {
		log.Printf("Rolling back migration: %s\n", migrationName)
		if err := applyMigration(db, rollbackFile, migrationName); err != nil {
			return err
		}
		log.Printf("Rolled back migration: %s\n", migrationName)
	} else {
		return fmt.Errorf("rollback file not found for migration: %s", migrationName)
	}

	return nil
}

func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		migration_name VARCHAR(255) UNIQUE NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`)
	return err
}

func getAppliedMigrations(db *sql.DB) ([]Migration, error) {
	rows, err := db.Query(`SELECT id, migration_name, created_at FROM migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var migration Migration
		if err := rows.Scan(&migration.ID, &migration.MigrationName, &migration.CreatedAt); err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}

	return migrations, nil
}

func isMigrationApplied(appliedMigrations []Migration, migrationName string) bool {
	for _, migration := range appliedMigrations {
		if migration.MigrationName == migrationName {
			return true
		}
	}
	return false
}

func applyMigration(db *sql.DB, filePath string, migrationName string) error {
	migrationSQL, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(migrationSQL)); err != nil {
		tx.Rollback()
		return err
	}

	if !strings.Contains(filePath, "_rollback") {
		if _, err := tx.Exec(`INSERT INTO migrations (migration_name) VALUES ($1)`, migrationName); err != nil {
			tx.Rollback()
			return err
		}
		log.Printf("Logged applied migration: %s\n", migrationName)
	} else {
		if _, err := tx.Exec(`DELETE FROM migrations WHERE migration_name = $1`, migrationName); err != nil {
			tx.Rollback()
			return err
		}
		log.Printf("Logged rolled back migration: %s\n", migrationName)
	}

	return tx.Commit()
}
