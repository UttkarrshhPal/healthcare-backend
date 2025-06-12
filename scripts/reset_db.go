package main

import (
    "log"
    "os"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL not set")
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Drop tables in reverse order due to foreign key constraints
    tables := []string{"appointments", "patients", "users"}
    
    for _, table := range tables {
        if err := db.Exec("DROP TABLE IF EXISTS " + table + " CASCADE").Error; err != nil {
            log.Printf("Failed to drop table %s: %v", table, err)
        } else {
            log.Printf("Dropped table: %s", table)
        }
    }

    log.Println("Database reset completed!")
}