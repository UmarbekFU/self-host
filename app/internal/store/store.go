package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type Store struct {
	db *sql.DB
}

type Subscriber struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	Status         string    `json:"status"`
	Attributes     json.RawMessage `json:"attributes"`
	CreatedAt      time.Time `json:"created_at"`
	UnsubscribedAt *time.Time `json:"unsubscribed_at,omitempty"`
}

type List struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Campaign struct {
	ID           int        `json:"id"`
	ListID       int        `json:"list_id"`
	Subject      string     `json:"subject"`
	HTML         string     `json:"html"`
	Text         string     `json:"text"`
	FromName     string     `json:"from_name"`
	FromEmail    string     `json:"from_email"`
	ReplyTo      string     `json:"reply_to"`
	Status       string     `json:"status"`
	ScheduledAt  *time.Time `json:"scheduled_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	SentAt       *time.Time `json:"sent_at,omitempty"`
}

type Event struct {
	ID           int             `json:"id"`
	CampaignID   int             `json:"campaign_id"`
	SubscriberID int             `json:"subscriber_id"`
	Type         string          `json:"type"`
	Meta         json.RawMessage `json:"meta"`
	At           time.Time       `json:"at"`
}

type Suppression struct {
	Email  string    `json:"email"`
	Reason string    `json:"reason"`
	At     time.Time `json:"at"`
}

type Job struct {
	ID        int             `json:"id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	RunAt     time.Time       `json:"run_at"`
	Attempts  int             `json:"attempts"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type Domain struct {
	ID              int        `json:"id"`
	Domain          string     `json:"domain"`
	DKIMSelector    string     `json:"dkim_selector"`
	DKIMPrivateKey  string     `json:"dkim_private_key"`
	DKIMPublicKey   string     `json:"dkim_public_key"`
	SPFRecord       string     `json:"spf_record"`
	DMARCRecord     string     `json:"dmarc_record"`
	PTRRecord       string     `json:"ptr_record"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type User struct {
	ID           int        `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	Role         string     `json:"role"`
	CreatedAt    time.Time  `json:"created_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

func Open(dsn string) (*Store, error) {
	// Parse DSN to extract database path
	var dbPath string
	if strings.HasPrefix(dsn, "sqlite://") {
		dbPath = strings.TrimPrefix(dsn, "sqlite://")
	} else {
		dbPath = dsn
	}

	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=1")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Set connection limits
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) DB() *sql.DB {
	return s.db
}

func Migrate(db *sql.DB) error {
	// Read migration files
	migrationDir := "migrations"
	files, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		migrationPath := filepath.Join(migrationDir, file.Name())
		content, err := ioutil.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file.Name(), err)
		}

		logrus.Infof("Running migration: %s", file.Name())
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", file.Name(), err)
		}
	}

	return nil
}

// Subscriber methods
func (s *Store) CreateSubscriber(email string, attributes json.RawMessage) (*Subscriber, error) {
	query := `INSERT INTO subscribers (email, status, attributes) VALUES (?, 'active', ?)`
	result, err := s.db.Exec(query, email, attributes)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetSubscriber(int(id))
}

func (s *Store) GetSubscriber(id int) (*Subscriber, error) {
	query := `SELECT id, email, status, attributes, created_at, unsubscribed_at FROM subscribers WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var sub Subscriber
	var unsubscribedAt sql.NullTime
	err := row.Scan(&sub.ID, &sub.Email, &sub.Status, &sub.Attributes, &sub.CreatedAt, &unsubscribedAt)
	if err != nil {
		return nil, err
	}

	if unsubscribedAt.Valid {
		sub.UnsubscribedAt = &unsubscribedAt.Time
	}

	return &sub, nil
}

func (s *Store) GetSubscriberByEmail(email string) (*Subscriber, error) {
	query := `SELECT id, email, status, attributes, created_at, unsubscribed_at FROM subscribers WHERE email = ?`
	row := s.db.QueryRow(query, email)

	var sub Subscriber
	var unsubscribedAt sql.NullTime
	err := row.Scan(&sub.ID, &sub.Email, &sub.Status, &sub.Attributes, &sub.CreatedAt, &unsubscribedAt)
	if err != nil {
		return nil, err
	}

	if unsubscribedAt.Valid {
		sub.UnsubscribedAt = &unsubscribedAt.Time
	}

	return &sub, nil
}

func (s *Store) UpdateSubscriberStatus(id int, status string) error {
	query := `UPDATE subscribers SET status = ?, unsubscribed_at = ? WHERE id = ?`
	now := time.Now()
	if status != "unsubscribed" {
		now = time.Time{}
	}
	_, err := s.db.Exec(query, status, now, id)
	return err
}

// List methods
func (s *Store) CreateList(name, description string) (*List, error) {
	query := `INSERT INTO lists (name, description) VALUES (?, ?)`
	result, err := s.db.Exec(query, name, description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetList(int(id))
}

func (s *Store) GetList(id int) (*List, error) {
	query := `SELECT id, name, description, created_at FROM lists WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var list List
	err := row.Scan(&list.ID, &list.Name, &list.Description, &list.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

func (s *Store) GetLists() ([]*List, error) {
	query := `SELECT id, name, description, created_at FROM lists ORDER BY created_at DESC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []*List
	for rows.Next() {
		var list List
		err := rows.Scan(&list.ID, &list.Name, &list.Description, &list.CreatedAt)
		if err != nil {
			return nil, err
		}
		lists = append(lists, &list)
	}

	return lists, nil
}

// Campaign methods
func (s *Store) CreateCampaign(campaign *Campaign) error {
	query := `INSERT INTO campaigns (list_id, subject, html, text, from_name, from_email, reply_to, status, scheduled_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := s.db.Exec(query, campaign.ListID, campaign.Subject, campaign.HTML, campaign.Text, 
		campaign.FromName, campaign.FromEmail, campaign.ReplyTo, campaign.Status, campaign.ScheduledAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	campaign.ID = int(id)
	return nil
}

func (s *Store) GetCampaign(id int) (*Campaign, error) {
	query := `SELECT id, list_id, subject, html, text, from_name, from_email, reply_to, status, scheduled_at, created_at, sent_at 
			  FROM campaigns WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var campaign Campaign
	var scheduledAt, sentAt sql.NullTime
	err := row.Scan(&campaign.ID, &campaign.ListID, &campaign.Subject, &campaign.HTML, &campaign.Text,
		&campaign.FromName, &campaign.FromEmail, &campaign.ReplyTo, &campaign.Status, &scheduledAt, &campaign.CreatedAt, &sentAt)
	if err != nil {
		return nil, err
	}

	if scheduledAt.Valid {
		campaign.ScheduledAt = &scheduledAt.Time
	}
	if sentAt.Valid {
		campaign.SentAt = &sentAt.Time
	}

	return &campaign, nil
}

func (s *Store) GetCampaigns() ([]*Campaign, error) {
	query := `SELECT id, list_id, subject, html, text, from_name, from_email, reply_to, status, scheduled_at, created_at, sent_at 
			  FROM campaigns ORDER BY created_at DESC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []*Campaign
	for rows.Next() {
		var campaign Campaign
		var scheduledAt, sentAt sql.NullTime
		err := rows.Scan(&campaign.ID, &campaign.ListID, &campaign.Subject, &campaign.HTML, &campaign.Text,
			&campaign.FromName, &campaign.FromEmail, &campaign.ReplyTo, &campaign.Status, &scheduledAt, &campaign.CreatedAt, &sentAt)
		if err != nil {
			return nil, err
		}

		if scheduledAt.Valid {
			campaign.ScheduledAt = &scheduledAt.Time
		}
		if sentAt.Valid {
			campaign.SentAt = &sentAt.Time
		}

		campaigns = append(campaigns, &campaign)
	}

	return campaigns, nil
}

func (s *Store) UpdateCampaignStatus(id int, status string) error {
	query := `UPDATE campaigns SET status = ?, sent_at = ? WHERE id = ?`
	now := time.Now()
	if status != "sent" {
		now = time.Time{}
	}
	_, err := s.db.Exec(query, status, now, id)
	return err
}

// Event methods
func (s *Store) RecordEvent(campaignID, subscriberID int, eventType string, meta json.RawMessage) error {
	query := `INSERT INTO events (campaign_id, subscriber_id, type, meta) VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(query, campaignID, subscriberID, eventType, meta)
	return err
}

func (s *Store) GetCampaignEvents(campaignID int) ([]*Event, error) {
	query := `SELECT id, campaign_id, subscriber_id, type, meta, at FROM events WHERE campaign_id = ? ORDER BY at DESC`
	rows, err := s.db.Query(query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.CampaignID, &event.SubscriberID, &event.Type, &event.Meta, &event.At)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

// Suppression methods
func (s *Store) AddSuppression(email, reason string) error {
	query := `INSERT OR REPLACE INTO suppressions (email, reason) VALUES (?, ?)`
	_, err := s.db.Exec(query, email, reason)
	return err
}

func (s *Store) IsSuppressed(email string) (bool, error) {
	query := `SELECT 1 FROM suppressions WHERE email = ? LIMIT 1`
	var exists int
	err := s.db.QueryRow(query, email).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

// Job methods
func (s *Store) EnqueueJob(jobType string, payload json.RawMessage, runAt time.Time) error {
	query := `INSERT INTO jobs (type, payload, run_at) VALUES (?, ?, ?)`
	_, err := s.db.Exec(query, jobType, payload, runAt)
	return err
}

func (s *Store) GetNextJob() (*Job, error) {
	query := `SELECT id, type, payload, run_at, attempts, status, created_at, updated_at 
			  FROM jobs WHERE status = 'queued' AND run_at <= ? ORDER BY run_at ASC LIMIT 1`
	row := s.db.QueryRow(query, time.Now())

	var job Job
	err := row.Scan(&job.ID, &job.Type, &job.Payload, &job.RunAt, &job.Attempts, &job.Status, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (s *Store) UpdateJobStatus(id int, status string) error {
	query := `UPDATE jobs SET status = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, status, time.Now(), id)
	return err
}

func (s *Store) IncrementJobAttempts(id int) error {
	query := `UPDATE jobs SET attempts = attempts + 1, updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, time.Now(), id)
	return err
}

// Domain methods
func (s *Store) CreateDomain(domain *Domain) error {
	query := `INSERT INTO domains (domain, dkim_selector, dkim_private_key, dkim_public_key, spf_record, dmarc_record, ptr_record) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := s.db.Exec(query, domain.Domain, domain.DKIMSelector, domain.DKIMPrivateKey, 
		domain.DKIMPublicKey, domain.SPFRecord, domain.DMARCRecord, domain.PTRRecord)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	domain.ID = int(id)
	return nil
}

func (s *Store) GetDomain(id int) (*Domain, error) {
	query := `SELECT id, domain, dkim_selector, dkim_private_key, dkim_public_key, spf_record, dmarc_record, ptr_record, verified_at, created_at 
			  FROM domains WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var domain Domain
	var verifiedAt sql.NullTime
	err := row.Scan(&domain.ID, &domain.Domain, &domain.DKIMSelector, &domain.DKIMPrivateKey, 
		&domain.DKIMPublicKey, &domain.SPFRecord, &domain.DMARCRecord, &domain.PTRRecord, &verifiedAt, &domain.CreatedAt)
	if err != nil {
		return nil, err
	}

	if verifiedAt.Valid {
		domain.VerifiedAt = &verifiedAt.Time
	}

	return &domain, nil
}

func (s *Store) GetDomains() ([]*Domain, error) {
	query := `SELECT id, domain, dkim_selector, dkim_private_key, dkim_public_key, spf_record, dmarc_record, ptr_record, verified_at, created_at 
			  FROM domains ORDER BY created_at DESC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []*Domain
	for rows.Next() {
		var domain Domain
		var verifiedAt sql.NullTime
		err := rows.Scan(&domain.ID, &domain.Domain, &domain.DKIMSelector, &domain.DKIMPrivateKey, 
			&domain.DKIMPublicKey, &domain.SPFRecord, &domain.DMARCRecord, &domain.PTRRecord, &verifiedAt, &domain.CreatedAt)
		if err != nil {
			return nil, err
		}

		if verifiedAt.Valid {
			domain.VerifiedAt = &verifiedAt.Time
		}

		domains = append(domains, &domain)
	}

	return domains, nil
}

func (s *Store) UpdateDomainVerification(id int, verified bool) error {
	var query string
	if verified {
		query = `UPDATE domains SET verified_at = ? WHERE id = ?`
		_, err := s.db.Exec(query, time.Now(), id)
		return err
	} else {
		query = `UPDATE domains SET verified_at = NULL WHERE id = ?`
		_, err := s.db.Exec(query, id)
		return err
	}
}
