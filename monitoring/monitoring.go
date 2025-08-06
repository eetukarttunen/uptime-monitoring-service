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

var allowedHosts = map[string]bool{
	"example.com":        true,
	"status.example.org": true,
}

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

func isValidURL(input string) bool {
	parsed, err := url.ParseRequestURI(input)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	return allowedHosts[parsed.Host]
}

func CheckWebsite(rawURL string, db *sql.DB) {
	if !isValidURL(rawURL) {
		log.Printf("Skipped invalid or unauthorized URL: %s", rawURL)
		return
	}

	start := time.Now()
	resp, err := http.Get(rawURL)

	var statusCode int
	if err != nil {
		statusCode = 0
	} else {
		statusCode = resp.StatusCode
		if cerr := resp.Body.Close(); cerr != nil {
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

func StartMonitoring(db *sql.DB, urls []string, interval time.Duration) {
	for {
		for _, u := range urls {
			go CheckWebsite(u, db)
		}
		time.Sleep(interval)
	}
}
