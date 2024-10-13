package db

import (
    "database/sql"
    "log"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/sqlite3"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// ConnectSQLite initializes the SQLite connection and enables WAL mode
func ConnectSQLite(dbPath string) error {
    var err error
    DB, err = sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatalf("Failed to open SQLite database: %v", err)
        return err
    }

    // Enable Write-Ahead Logging (WAL) mode to improve concurrency
    _, err = DB.Exec("PRAGMA journal_mode = WAL;")
    if err != nil {
        log.Fatalf("Failed to enable WAL mode: %v", err)
        return err
    }
    log.Println("WAL mode enabled for SQLite.")

    if err := DB.Ping(); err != nil {
        log.Fatalf("Failed to connect to SQLite database: %v", err)
        return err
    }
    log.Println("Successfully connected to SQLite database.")
    return nil
}

// Migrate applies database migrations from the specified migration path
func Migrate(migrationPath string) {
    driver, err := sqlite3.WithInstance(DB, &sqlite3.Config{})
    if err != nil {
        log.Fatalf("Failed to create SQLite driver: %v", err)
    }
    m, err := migrate.NewWithDatabaseInstance(
        "file://"+migrationPath,
        "sqlite3", driver)

    if err != nil {
        log.Fatalf("Migration failed: %v", err)
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Migration failed: %v", err)
    } else {
        log.Println("Migrations applied successfully.")
    }
}

// CloseSQLite closes the SQLite database connection
func CloseSQLite() {
    if DB != nil {
        DB.Close()
        log.Println("SQLite database connection closed.")
    }
}
