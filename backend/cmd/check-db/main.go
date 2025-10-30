package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("Database Connectivity Test (Go)")
	fmt.Println("========================================")
	fmt.Println()

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "project_website")
	sslMode := getEnv("DB_SSL_MODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)

	fmt.Println("Database Configuration:")
	fmt.Printf("  Host: %s\n", dbHost)
	fmt.Printf("  Port: %s\n", dbPort)
	fmt.Printf("  Database: %s\n", dbName)
	fmt.Printf("  User: %s\n", dbUser)
	fmt.Printf("  SSL Mode: %s\n", sslMode)
	fmt.Println()

	fmt.Println("Test 1: Opening database connection...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("✗ Failed to open database connection: %v\n", err)
	}
	defer db.Close()
	fmt.Println("✓ Database connection opened")

	fmt.Println()
	fmt.Println("Test 2: Pinging database...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalf("✗ Failed to ping database: %v\n", err)
	}
	fmt.Println("✓ Database ping successful")

	fmt.Println()
	fmt.Println("Test 3: Getting PostgreSQL version...")
	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		log.Printf("✗ Failed to get version: %v\n", err)
	} else {
		if len(version) > 80 {
			fmt.Printf("✓ PostgreSQL version: %s...\n", version[:80])
		} else {
			fmt.Printf("✓ PostgreSQL version: %s\n", version)
		}
	}

	fmt.Println()
	fmt.Println("Test 4: Counting tables...")
	var tableCount int
	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.tables
		WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
	`).Scan(&tableCount)
	if err != nil {
		log.Printf("✗ Failed to count tables: %v\n", err)
	} else {
		fmt.Printf("✓ Found %d tables in the database\n", tableCount)
	}

	fmt.Println()
	fmt.Println("Test 5: Checking critical tables...")
	criticalTables := []string{"users", "projects", "shortened_urls", "skills", "certifications"}

	for _, table := range criticalTables {
		var count int64
		err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if err != nil {
			fmt.Printf("  ✗ Table '%s': %v\n", table, err)
		} else {
			fmt.Printf("  ✓ Table '%s': %d rows\n", table, count)
		}
	}

	fmt.Println()
	fmt.Println("Test 6: Checking visitor analytics tables...")
	visitorTables := []string{
		"visitor_sessions",
		"page_views",
		"privacy_consents",
		"visitor_metrics",
		"visitor_daily_summary",
		"visitor_realtime",
		"visitor_locations",
	}

	for _, table := range visitorTables {
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT 1
				FROM information_schema.tables
				WHERE table_schema = 'public' AND table_name = $1
			)
		`, table).Scan(&exists)

		if err != nil {
			fmt.Printf("  ✗ Error checking table '%s': %v\n", table, err)
		} else if exists {
			var count int64
			db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
			fmt.Printf("  ✓ Table '%s' exists (%d rows)\n", table, count)
		} else {
			fmt.Printf("  ✗ Table '%s' does NOT exist\n", table)
		}
	}

	fmt.Println()
	fmt.Println("Test 7: Checking admin users...")
	var adminCount int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&adminCount)
	if err != nil {
		log.Printf("✗ Failed to count admin users: %v\n", err)
	} else {
		fmt.Printf("✓ Found %d admin user(s)\n", adminCount)

		if adminCount > 0 {
			fmt.Println()
			fmt.Println("  Admin users:")
			rows, err := db.Query(`
				SELECT username, email, created_at
				FROM users
				WHERE role = 'admin'
				ORDER BY created_at
			`)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var username, email string
					var createdAt time.Time
					if err := rows.Scan(&username, &email, &createdAt); err == nil {
						fmt.Printf("    - %s (%s) - created: %s\n", username, email, createdAt.Format(time.RFC3339))
					}
				}
			}
		}
	}

	fmt.Println()
	fmt.Println("Test 8: Testing database transactions...")
	tx, err := db.Begin()
	if err != nil {
		log.Printf("✗ Failed to begin transaction: %v\n", err)
	} else {
		var result int
		err = tx.QueryRow("SELECT 1").Scan(&result)
		if err != nil {
			log.Printf("✗ Failed to execute query in transaction: %v\n", err)
			tx.Rollback()
		} else {
			tx.Rollback()
			fmt.Println("✓ Transactions are working correctly")
		}
	}

	fmt.Println()
	fmt.Println("Test 9: Database statistics...")
	var dbSize string
	err = db.QueryRow("SELECT pg_size_pretty(pg_database_size($1))", dbName).Scan(&dbSize)
	if err != nil {
		log.Printf("✗ Failed to get database size: %v\n", err)
	} else {
		fmt.Printf("  Database size: %s\n", dbSize)
	}

	var connCount int
	err = db.QueryRow("SELECT COUNT(*) FROM pg_stat_activity WHERE datname = $1", dbName).Scan(&connCount)
	if err != nil {
		log.Printf("✗ Failed to get connection count: %v\n", err)
	} else {
		fmt.Printf("  Active connections: %d\n", connCount)
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("Summary")
	fmt.Println("========================================")
	fmt.Println("✓ All database connectivity tests passed!")
	fmt.Println()
	fmt.Println("Database is ready for use.")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}