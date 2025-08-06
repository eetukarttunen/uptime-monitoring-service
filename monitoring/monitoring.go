package monitoring

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Connection to the dockerized postgresql database
func ConnectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)
	return sql.Open("postgres", connStr)
}

// Creating a table if one doesn't exist yet
func CreateTableIfNotExists(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS uptime_logs (
		id SERIAL PRIMARY KEY,
		url TEXT NOT NULL,
		status_code INT,
		response_time_ms INT,
		checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("could not create table: %v", err)
	}
	log.Println("Table 'uptime_logs' is ready.")
	return nil
}

// ValidateURL ensures the URL is well-formed and uses HTTP or HTTPS
func ValidateURL(rawURL string) (string, error) {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("unsupported URL scheme: %s", parsed.Scheme)
	}
	return rawURL, nil
}

// Checking status and saving the status to database
func CheckWebsite(rawURL string, db *sql.DB) {
	validatedURL, err := ValidateURL(rawURL)
	if err != nil {
		log.Printf("Invalid URL '%s': %v", rawURL, err)
		return
	}

	start := time.Now()
	resp, err := http.Get(validatedURL) // G107 resolved via validation
	var statusCode int

	if err != nil {
		statusCode = 0
	} else {
		statusCode = resp.StatusCode
		if cerr := resp.Body.Close(); cerr != nil { // G104 fix: handle close error
			log.Printf("Failed to close response body: %v", cerr)
		}
	}

	responseTime := time.Since(start).Milliseconds()

	_, err = db.Exec("INSERT INTO uptime_logs (url, status_code, response_time_ms) VALUES ($1, $2, $3)", validatedURL, statusCode, responseTime)
	if err != nil {
		log.Printf("DB Error: %v", err)
	}

	fmt.Printf("Checked %s - Status: %d, Response Time: %dms\n", validatedURL, statusCode, responseTime)
}

func StartMonitoring(db *sql.DB, urls []string, interval time.Duration) {
	for {
		for _, url := range urls {
			go CheckWebsite(url, db)
		}
		time.Sleep(interval)
	}
}

