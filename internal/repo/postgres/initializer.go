package postgres

import (
	_ "embed"
	"fmt"
	"log"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/repo"
	"gorm.io/gorm"
)

//go:embed sql/schema.sql
var schemaSQL string

func InitStore(dbConnector func() (*gorm.DB, error)) (repo.Store, error) {
	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = dbConnector()
		if err == nil {
			var sqlDB, dbErr = db.DB()
			if dbErr != nil {
				err = dbErr
			} else {
				err = sqlDB.Ping()
			}
		}

		if err == nil {
			break
		}

		log.Printf("Postgres not ready (attempt %d/10): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("postgres connection failed after retries: %w", err)
	}

	log.Println("Applying database schema...")
	if err := db.Exec(schemaSQL).Error; err != nil {
		return nil, fmt.Errorf("failed to execute schema.sql: %w", err)
	}

	return New(db)
}
