// This program performs administrative tasks for the garage sale service.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/HMadhav/service/business/data/schema"
	"github.com/HMadhav/service/business/sys/database"
)

func main() {

	cfg := database.Config{
		User:         "postgres",
		Password:     "postgres",
		Host:         "localhost",
		Name:         "postgres",
		MaxIdleConns: 0,
		MaxOpenConns: 0,
		DisableTLS:   true,
	}

	if err := migrate(cfg); err != nil {
		fmt.Println("migrating database: %w", err)
		os.Exit(1)
	}

	if err := seed(cfg); err != nil {
		fmt.Println("seeding database: %w", err)
		os.Exit(1)
	}

}

// Migrate creates the schema in the database.
func migrate(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := schema.Migrate(ctx, db); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	fmt.Println("migrations complete")
	return nil
}

// Seed loads test data into the database.
func seed(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := schema.Seed(ctx, db); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	fmt.Println("seed data complete")
	return nil
}
