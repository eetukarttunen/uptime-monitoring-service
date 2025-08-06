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

// Whitelisted hostnames
var allowedHosts = map[string]bool{
	"example.com":        true,
	"status.example.org": true,
}

// Connection to the dockerized PostgreSQL database
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

// CreateTableIfNotExists ensures the logs table exists
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

// isValidURL checks scheme and host against whitelist
func isValidURL(input string) bool {
	parsed, err := url.ParseRequestURI(input)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return false
	}
	// Enforce https/http only
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	// Check if host is allowed
	return allowedHosts[parsed.Host]
}

// CheckWebsite checks status and logs it to the database
func CheckWebsite(rawURL string, db *sql.DB) {
	if !isValidURL(rawURL) {
		log.Printf("Skipped invalid or unauthorized URL: %s", rawURL)
		return
	}

	start := time.Now()
	resp, err := http.Get(rawURL) // G107 safe — input validated

	var statusCode int
	if err != nil {
		statusCode = 0
	} else {
		statusCode = resp.StatusCode
		if cerr := resp.Body.Close(); cerr != nil { // G104: handle Close error
			log.Printf("Error closing response body for %s: %v", rawURL, cerr)
		}
	}

	responseTime := time.Since(start).Milliseconds()

	_, err = db.Exec(`INSERT INTO uptime_logs (url, status_code, response_time_ms) VALUES ($1, $2, $3)`,
		rawURL, statusCode, responseTime)
	if err != nil {
		log.Printf("DB insert error: %v", err)
	}

	log.Printf("Checked %s - Status: %d, Response Time: %dms", rawURL, statusCode, responseTime)
}

// StartMonitoring begins polling the list of URLs
func StartMonitoring(db *sql.DB, urls []string, interval time.Duration) {
	for {
		for _, u := range urls {
			go CheckWebsite(u, db)
		}
		time.Sleep(interval)
	}
}
