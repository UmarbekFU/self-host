package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"newsletter/internal/store"
	"github.com/sirupsen/logrus"
)

type Queue struct {
	db *store.Store
}

type JobHandler func(ctx context.Context, payload json.RawMessage) error

func NewQueue(db *store.Store) *Queue {
	return &Queue{db: db}
}

func (q *Queue) Enqueue(jobType string, payload interface{}, runAt time.Time) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	return q.db.EnqueueJob(jobType, payloadBytes, runAt)
}

func (q *Queue) RunWorkers(numWorkers int, handlers map[string]JobHandler) {
	for i := 0; i < numWorkers; i++ {
		go q.worker(fmt.Sprintf("worker-%d", i), handlers)
	}
}

func (q *Queue) worker(name string, handlers map[string]JobHandler) {
	logrus.Infof("Starting worker: %s", name)
	
	for {
		// Get next job
		job, err := q.db.GetNextJob()
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				// No jobs available, wait a bit
				time.Sleep(5 * time.Second)
				continue
			}
			logrus.Errorf("Failed to get next job: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Mark job as running
		if err := q.db.UpdateJobStatus(job.ID, "running"); err != nil {
			logrus.Errorf("Failed to update job status to running: %v", err)
			continue
		}

		// Process job
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		err = q.processJob(ctx, job, handlers)
		cancel()

		// Update job status
		if err != nil {
			logrus.Errorf("Job %d failed: %v", job.ID, err)
			
			// Increment attempts
			if err := q.db.IncrementJobAttempts(job.ID); err != nil {
				logrus.Errorf("Failed to increment job attempts: %v", err)
			}
			
			// Check if we should retry
			if job.Attempts >= 3 {
				q.db.UpdateJobStatus(job.ID, "failed")
			} else {
				// Retry with exponential backoff
				retryAt := time.Now().Add(time.Duration(job.Attempts+1) * time.Minute)
				q.db.EnqueueJob(job.Type, job.Payload, retryAt)
				q.db.UpdateJobStatus(job.ID, "queued")
			}
		} else {
			q.db.UpdateJobStatus(job.ID, "done")
		}
	}
}

func (q *Queue) processJob(ctx context.Context, job *store.Job, handlers map[string]JobHandler) error {
	handler, exists := handlers[job.Type]
	if !exists {
		return fmt.Errorf("no handler for job type: %s", job.Type)
	}

	return handler(ctx, job.Payload)
}

// Job payloads
type SendBatchPayload struct {
	CampaignID int      `json:"campaign_id"`
	Recipients []string `json:"recipients"`
}

type BounceProcessingPayload struct {
	Email     string `json:"email"`
	Reason    string `json:"reason"`
	BounceType string `json:"bounce_type"`
}

type DKIMRotationPayload struct {
	DomainID int `json:"domain_id"`
}

// Job handlers
func (q *Queue) SendBatchHandler(ctx context.Context, payload json.RawMessage) error {
	var p SendBatchPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("failed to unmarshal send batch payload: %w", err)
	}

	// Get campaign
	campaign, err := q.db.GetCampaign(p.CampaignID)
	if err != nil {
		return fmt.Errorf("failed to get campaign: %w", err)
	}

	// Get list
	_, err = q.db.GetList(campaign.ListID)
	if err != nil {
		return fmt.Errorf("failed to get list: %w", err)
	}

	// Process each recipient
	for _, email := range p.Recipients {
		// Check if suppressed
		suppressed, err := q.db.IsSuppressed(email)
		if err != nil {
			logrus.Errorf("Failed to check suppression for %s: %v", email, err)
			continue
		}
		if suppressed {
			continue
		}

		// Get or create subscriber
		subscriber, err := q.db.GetSubscriberByEmail(email)
		if err != nil {
			// Create new subscriber
			subscriber, err = q.db.CreateSubscriber(email, json.RawMessage("{}"))
			if err != nil {
				logrus.Errorf("Failed to create subscriber %s: %v", email, err)
				continue
			}
		}

		// TODO: Send email
		// This would involve creating the email message and sending it via SMTP
		logrus.Infof("Would send email to %s for campaign %d", email, p.CampaignID)
		
		// Record delivery event
		if err := q.db.RecordEvent(p.CampaignID, subscriber.ID, "delivered", nil); err != nil {
			logrus.Errorf("Failed to record delivery event: %v", err)
		}
	}

	return nil
}

func (q *Queue) BounceProcessingHandler(ctx context.Context, payload json.RawMessage) error {
	var p BounceProcessingPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("failed to unmarshal bounce processing payload: %w", err)
	}

	// Add to suppressions
	if err := q.db.AddSuppression(p.Email, p.Reason); err != nil {
		return fmt.Errorf("failed to add suppression: %w", err)
	}

	// Update subscriber status
	subscriber, err := q.db.GetSubscriberByEmail(p.Email)
	if err != nil {
		logrus.Warnf("Subscriber not found for bounce: %s", p.Email)
		return nil
	}

	status := "bounced"
	if p.BounceType == "complaint" {
		status = "complained"
	}

	if err := q.db.UpdateSubscriberStatus(subscriber.ID, status); err != nil {
		return fmt.Errorf("failed to update subscriber status: %w", err)
	}

	logrus.Infof("Processed bounce for %s: %s", p.Email, p.Reason)
	return nil
}

func (q *Queue) DKIMRotationHandler(ctx context.Context, payload json.RawMessage) error {
	var p DKIMRotationPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("failed to unmarshal DKIM rotation payload: %w", err)
	}

	// TODO: Implement DKIM key rotation
	logrus.Infof("DKIM rotation for domain %d not implemented", p.DomainID)
	return nil
}
