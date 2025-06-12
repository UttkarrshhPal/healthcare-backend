package database

import (
    "context"
    "database/sql"
    "log"
    "os"
    "time"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Initialize() {
    // Get database URL from environment
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL is not set in .env file")
    }

    // Configure GORM logger
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             time.Second * 10,
            LogLevel:                  logger.Silent, // Set to Silent to avoid parameter errors
            IgnoreRecordNotFoundError: true,
            Colorful:                  false,
        },
    )

    // Database configuration for Neon DB
    config := &gorm.Config{
        Logger: newLogger,
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
        DisableForeignKeyConstraintWhenMigrating: true,
        SkipDefaultTransaction:                   true,
    }

    // Connect to database
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), config)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    log.Println("✓ Connected to Neon DB successfully")

    // Get underlying SQL database
    sqlDB, err := DB.DB()
    if err != nil {
        log.Fatalf("Failed to get database instance: %v", err)
    }

    // Configure connection pool
    sqlDB.SetMaxIdleConns(2)
    sqlDB.SetMaxOpenConns(5)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)

    // Verify connection
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := sqlDB.PingContext(ctx); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    log.Println("✓ Database connection verified")

    // Create tables manually
    if err := createTables(sqlDB); err != nil {
        log.Fatalf("Failed to create tables: %v", err)
    }

    log.Println("✓ Database initialization completed")
}

func createTables(sqlDB *sql.DB) error {
    log.Println("Creating database tables...")

    // Create tables with raw SQL to avoid GORM parameter issues
    tables := []struct {
        name string
        sql  string
    }{
        {
            name: "users",
            sql: `
                CREATE TABLE IF NOT EXISTS users (
                    id SERIAL PRIMARY KEY,
                    email VARCHAR(255) UNIQUE NOT NULL,
                    password VARCHAR(255) NOT NULL,
                    name VARCHAR(255) NOT NULL,
                    role VARCHAR(50) NOT NULL CHECK (role IN ('receptionist', 'doctor')),
                    is_active BOOLEAN DEFAULT true,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    deleted_at TIMESTAMP
                )`,
        },
        {
            name: "patients",
            sql: `
                CREATE TABLE IF NOT EXISTS patients (
                    id SERIAL PRIMARY KEY,
                    first_name VARCHAR(255) NOT NULL,
                    last_name VARCHAR(255) NOT NULL,
                    email VARCHAR(255) UNIQUE,
                    phone VARCHAR(50) NOT NULL,
                    date_of_birth TIMESTAMP,
                    gender VARCHAR(20),
                    address TEXT,
                    medical_history TEXT,
                    current_medication TEXT,
                    allergies TEXT,
                    emergency_contact VARCHAR(50),
                    blood_group VARCHAR(10),
                    insurance_number VARCHAR(50),
                    registered_by INTEGER REFERENCES users(id),
                    last_updated_by INTEGER REFERENCES users(id),
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    deleted_at TIMESTAMP
                )`,
        },
        {
            name: "appointments",
            sql: `
                CREATE TABLE IF NOT EXISTS appointments (
                    id SERIAL PRIMARY KEY,
                    patient_id INTEGER NOT NULL REFERENCES patients(id),
                    doctor_id INTEGER NOT NULL REFERENCES users(id),
                    date TIMESTAMP,
                    time VARCHAR(10),
                    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'completed', 'cancelled')),
                    notes TEXT,
                    created_by INTEGER REFERENCES users(id),
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    deleted_at TIMESTAMP
                )`,
        },
    }

    // Create each table
    for _, table := range tables {
        _, err := sqlDB.Exec(table.sql)
        if err != nil {
            log.Printf("Error creating %s table: %v", table.name, err)
            return err
        }
        log.Printf("  ✓ %s table ready", table.name)
    }

    // Create indexes
    indexes := []string{
        "CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
        "CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at)",
        "CREATE INDEX IF NOT EXISTS idx_patients_email ON patients(email)",
        "CREATE INDEX IF NOT EXISTS idx_patients_phone ON patients(phone)",
        "CREATE INDEX IF NOT EXISTS idx_patients_deleted_at ON patients(deleted_at)",
        "CREATE INDEX IF NOT EXISTS idx_appointments_patient_id ON appointments(patient_id)",
        "CREATE INDEX IF NOT EXISTS idx_appointments_doctor_id ON appointments(doctor_id)",
        "CREATE INDEX IF NOT EXISTS idx_appointments_date ON appointments(date)",
        "CREATE INDEX IF NOT EXISTS idx_appointments_deleted_at ON appointments(deleted_at)",
    }

    for _, idx := range indexes {
        if _, err := sqlDB.Exec(idx); err != nil {
            log.Printf("Warning: Failed to create index: %v", err)
        }
    }

    log.Println("  ✓ Indexes created")

    return nil
}

func GetDB() *gorm.DB {
    if DB == nil {
        log.Fatal("Database not initialized. Call Initialize() first")
    }
    return DB
}

// HealthCheck performs a health check on the database connection
func HealthCheck() error {
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    return sqlDB.PingContext(ctx)
}

// Close closes the database connection
func Close() error {
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}