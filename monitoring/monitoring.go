package monitoring

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Dockerized postgresql connection
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

// Creating table if doesn't exist yet
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

// Checking status and saving it to database
func CheckWebsite(url string, db *sql.DB) {
	start := time.Now()
	resp, err := http.Get(url)

	var statusCode int
	if err != nil {
		statusCode = 0
	} else {
		statusCode = resp.StatusCode
		resp.Body.Close()
	}

	responseTime := time.Since(start).Milliseconds()

	_, err = db.Exec("INSERT INTO uptime_logs (url, status_code, response_time_ms) VALUES ($1, $2, $3)", url, statusCode, responseTime)
	if err != nil {
		fmt.Println("DB Error:", err)
	}

	fmt.Printf("Checked %s - Status: %d, Response Time: %dms\n", url, statusCode, responseTime)
}

func StartMonitoring(db *sql.DB, urls []string, interval time.Duration) {
	for {
		for _, url := range urls {
			go CheckWebsite(url, db)
		}
		time.Sleep(interval)
	}
}
