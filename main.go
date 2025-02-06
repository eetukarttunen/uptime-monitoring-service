package main

import (
	"log"
	"time"
	"uptime-monitoring-service/monitoring"
)

func main() {
	db, err := monitoring.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Creating postgresql table if it doesn't exist
	err = monitoring.CreateTableIfNotExists(db)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// The place to list all the websites to be monitored (add how many you ever want)
	urls := []string{
		"https://opiskelijaruokalista.vercel.app/",
		"https://eetukarttunen.vercel.app/",
	}

	// Monitoring interval, for example every 30 seconds.
	monitoring.StartMonitoring(db, urls, 30*time.Second)
}
