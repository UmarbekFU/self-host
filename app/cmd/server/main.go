package main

import (
	"log"
	"net/http"
	"os"

	httpapi "newsletter/internal/http"
	"newsletter/internal/store"
	"newsletter/internal/jobs"
	"newsletter/internal/mail"
	"newsletter/internal/deliverability"

	"github.com/sirupsen/logrus"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found")
	}

	// Initialize logger
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Get configuration from environment
	dsn := getEnv("DATABASE_URL", "sqlite:///var/app/newsletter.db")
	port := getEnv("PORT", "8080")
	licenseKey := getEnv("LICENSE_KEY", "")

	if licenseKey == "" {
		logrus.Fatal("LICENSE_KEY environment variable is required")
	}

	// Initialize database
	db, err := store.Open(dsn)
	if err != nil {
		logrus.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := store.Migrate(db.DB()); err != nil {
		logrus.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize services
	queue := jobs.NewQueue(db)
	mailService := mail.NewService()
	deliverabilityService := deliverability.NewService()
	
	// Create service container
	services := &httpapi.Services{
		DB: db,
		Queue: queue,
		Mail: mailService,
		Deliverability: deliverabilityService,
		LicenseKey: licenseKey,
	}

	// Start background workers
	handlers := map[string]jobs.JobHandler{
		"send_batch": queue.SendBatchHandler,
		"process_bounce": queue.BounceProcessingHandler,
		"rotate_dkim": queue.DKIMRotationHandler,
	}
	go queue.RunWorkers(4, handlers)

	// Setup HTTP routes
	mux := httpapi.NewRouter(services)

	// Serve static files (SvelteKit build output)
	fs := http.FileServer(http.Dir("./static"))
	mux.PathPrefix("/").Handler(fs)

	// Start server
	logrus.Infof("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
