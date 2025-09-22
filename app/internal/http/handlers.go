package http

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"newsletter/internal/store"
	"newsletter/internal/mail"
	"newsletter/internal/deliverability"
	"newsletter/internal/jobs"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Services struct {
	DB             *store.Store
	Queue          *jobs.Queue
	Mail           *mail.Service
	Deliverability *deliverability.Service
	LicenseKey     string
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewRouter(services *Services) *mux.Router {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	
	// Health check
	api.HandleFunc("/health", healthHandler).Methods("GET")
	
	// Auth routes
	api.HandleFunc("/auth/login", loginHandler(services)).Methods("POST")
	
	// Domain routes
	api.HandleFunc("/domains", createDomainHandler(services)).Methods("POST")
	api.HandleFunc("/domains", getDomainsHandler(services)).Methods("GET")
	api.HandleFunc("/domains/{id}", getDomainHandler(services)).Methods("GET")
	api.HandleFunc("/domains/{id}/status", getDomainStatusHandler(services)).Methods("GET")
	api.HandleFunc("/domains/{id}/dkim/rotate", rotateDKIMHandler(services)).Methods("POST")
	
	// List routes
	api.HandleFunc("/lists", createListHandler(services)).Methods("POST")
	api.HandleFunc("/lists", getListsHandler(services)).Methods("GET")
	api.HandleFunc("/lists/{id}", getListHandler(services)).Methods("GET")
	api.HandleFunc("/lists/{id}/import", importSubscribersHandler(services)).Methods("POST")
	api.HandleFunc("/lists/{id}/subscribers", getListSubscribersHandler(services)).Methods("GET")
	
	// Campaign routes
	api.HandleFunc("/campaigns", createCampaignHandler(services)).Methods("POST")
	api.HandleFunc("/campaigns", getCampaignsHandler(services)).Methods("GET")
	api.HandleFunc("/campaigns/{id}", getCampaignHandler(services)).Methods("GET")
	api.HandleFunc("/campaigns/{id}/test", testCampaignHandler(services)).Methods("POST")
	api.HandleFunc("/campaigns/{id}/schedule", scheduleCampaignHandler(services)).Methods("POST")
	api.HandleFunc("/campaigns/{id}/report", getCampaignReportHandler(services)).Methods("GET")
	
	// Tracking routes
	api.HandleFunc("/track/click", trackClickHandler(services)).Methods("POST")
	api.HandleFunc("/track/open", trackOpenHandler(services)).Methods("POST")
	
	// Unsubscribe route
	r.HandleFunc("/u/{subscriberId}/{token}", unsubscribeHandler(services)).Methods("GET")
	
	// Bounce webhook
	api.HandleFunc("/hooks/bounce", bounceHandler(services)).Methods("POST")

	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, APIResponse{Success: true, Data: map[string]string{"status": "ok"}})
}

func loginHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid request"}, http.StatusBadRequest)
			return
		}

		// TODO: Implement proper authentication
		// For now, just check if license key is valid
		if services.LicenseKey == "" {
			respondJSON(w, APIResponse{Success: false, Error: "License key not configured"}, http.StatusUnauthorized)
			return
		}

		respondJSON(w, APIResponse{Success: true, Data: map[string]string{"token": "dummy-token"}})
	}
}

func createDomainHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Domain string `json:"domain"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid request"}, http.StatusBadRequest)
			return
		}

		// Generate DKIM keys
		selector := "newsletter"
		privateKey, publicKey, err := mail.GenerateDKIMKeys()
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Failed to generate DKIM keys"}, http.StatusInternalServerError)
			return
		}

		// Create domain record
		domain := &store.Domain{
			Domain:          req.Domain,
			DKIMSelector:    selector,
			DKIMPrivateKey:  privateKey,
			DKIMPublicKey:   publicKey,
			SPFRecord:       fmt.Sprintf("v=spf1 a mx ip4:%s ~all", r.RemoteAddr), // TODO: Get actual server IP
			DMARCRecord:     fmt.Sprintf("v=DMARC1; p=quarantine; rua=mailto:dmarc@%s", req.Domain),
			PTRRecord:       fmt.Sprintf("mail.%s", req.Domain),
		}

		if err := services.DB.CreateDomain(domain); err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Failed to create domain"}, http.StatusInternalServerError)
			return
		}

		respondJSON(w, APIResponse{Success: true, Data: domain})
	}
}

func getDomainsHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domains, err := services.DB.GetDomains()
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Failed to get domains"}, http.StatusInternalServerError)
			return
		}

		respondJSON(w, APIResponse{Success: true, Data: domains})
	}
}

func getDomainHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid domain ID"}, http.StatusBadRequest)
			return
		}

		domain, err := services.DB.GetDomain(id)
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Domain not found"}, http.StatusNotFound)
			return
		}

		respondJSON(w, APIResponse{Success: true, Data: domain})
	}
}

func getDomainStatusHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid domain ID"}, http.StatusBadRequest)
			return
		}

		domain, err := services.DB.GetDomain(id)
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Domain not found"}, http.StatusNotFound)
			return
		}

		// Check DNS records
		status, err := services.Deliverability.CheckDomainStatus(domain)
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Failed to check domain status"}, http.StatusInternalServerError)
			return
		}

		respondJSON(w, APIResponse{Success: true, Data: status})
	}
}

func rotateDKIMHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid domain ID"}, http.StatusBadRequest)
			return
		}

		domain, err := services.DB.GetDomain(id)
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Domain not found"}, http.StatusNotFound)
			return
		}

		// Generate new DKIM keys
		selector := "newsletter"
		privateKey, publicKey, err := mail.GenerateDKIMKeys()
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Failed to generate DKIM keys"}, http.StatusInternalServerError)
			return
		}

		// Update domain with new keys
		domain.DKIMSelector = selector
		domain.DKIMPrivateKey = privateKey
		domain.DKIMPublicKey = publicKey

		// TODO: Update domain in database

		respondJSON(w, APIResponse{Success: true, Data: domain})
	}
}

func createListHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid request"}, http.StatusBadRequest)
			return
		}

		list, err := services.DB.CreateList(req.Name, req.Description)
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Failed to create list"}, http.StatusInternalServerError)
			return
		}

		respondJSON(w, APIResponse{Success: true, Data: list})
	}
}

func getListsHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lists, err := services.DB.GetLists()
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Failed to get lists"}, http.StatusInternalServerError)
			return
		}

		respondJSON(w, APIResponse{Success: true, Data: lists})
	}
}

func getListHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid list ID"}, http.StatusBadRequest)
			return
		}

		list, err := services.DB.GetList(id)
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "List not found"}, http.StatusNotFound)
			return
		}

		respondJSON(w, APIResponse{Success: true, Data: list})
	}
}

func importSubscribersHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		listID, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid list ID"}, http.StatusBadRequest)
			return
		}

		// Check if list exists
		_, err = services.DB.GetList(listID)
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "List not found"}, http.StatusNotFound)
			return
		}

		// Parse multipart form
		err = r.ParseMultipartForm(10 << 20) // 10MB max
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Failed to parse form"}, http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "No file uploaded"}, http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Parse CSV
		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Failed to parse CSV"}, http.StatusBadRequest)
			return
		}

		if len(records) < 2 {
			respondJSON(w, APIResponse{Success: false, Error: "CSV must have at least a header row and one data row"}, http.StatusBadRequest)
			return
		}

		// Process CSV
		imported := 0
		skipped := 0
		errors := []string{}

		// Get header row
		headers := records[0]
		emailIndex := -1
		for i, header := range headers {
			if strings.ToLower(strings.TrimSpace(header)) == "email" {
				emailIndex = i
				break
			}
		}

		if emailIndex == -1 {
			respondJSON(w, APIResponse{Success: false, Error: "CSV must contain an 'email' column"}, http.StatusBadRequest)
			return
		}

		// Process data rows
		for i, record := range records[1:] {
			if len(record) <= emailIndex {
				errors = append(errors, fmt.Sprintf("Row %d: insufficient columns", i+2))
				skipped++
				continue
			}

			email := strings.TrimSpace(record[emailIndex])
			if email == "" {
				errors = append(errors, fmt.Sprintf("Row %d: empty email", i+2))
				skipped++
				continue
			}

			// Validate email format
			if !isValidEmail(email) {
				errors = append(errors, fmt.Sprintf("Row %d: invalid email format: %s", i+2, email))
				skipped++
				continue
			}

			// Create subscriber attributes from other columns
			attributes := make(map[string]string)
			for j, value := range record {
				if j != emailIndex && j < len(headers) {
					attributes[headers[j]] = value
				}
			}

			attributesJSON, _ := json.Marshal(attributes)

			// Create or update subscriber
			_, err := services.DB.GetSubscriberByEmail(email)
			if err != nil {
				// Create new subscriber
				_, err = services.DB.CreateSubscriber(email, attributesJSON)
				if err != nil {
					errors = append(errors, fmt.Sprintf("Row %d: failed to create subscriber: %v", i+2, err))
					skipped++
					continue
				}
			}

			// Add to list (this would need to be implemented in the store)
			// For now, we'll just count as imported
			imported++
		}

		respondJSON(w, APIResponse{
			Success: true,
			Data: map[string]interface{}{
				"imported": imported,
				"skipped":  skipped,
				"errors":   errors,
			},
		})
	}
}

func getListSubscribersHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		_, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid list ID"}, http.StatusBadRequest)
			return
		}

		// TODO: Implement subscriber listing
		respondJSON(w, APIResponse{Success: true, Data: []interface{}{}})
	}
}

func createCampaignHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ListID    int    `json:"list_id"`
			Subject   string `json:"subject"`
			HTML      string `json:"html"`
			Text      string `json:"text"`
			FromName  string `json:"from_name"`
			FromEmail string `json:"from_email"`
			ReplyTo   string `json:"reply_to"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid request"}, http.StatusBadRequest)
			return
		}

		campaign := &store.Campaign{
			ListID:    req.ListID,
			Subject:   req.Subject,
			HTML:      req.HTML,
			Text:      req.Text,
			FromName:  req.FromName,
			FromEmail: req.FromEmail,
			ReplyTo:   req.ReplyTo,
			Status:    "draft",
		}

		if err := services.DB.CreateCampaign(campaign); err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Failed to create campaign"}, http.StatusInternalServerError)
			return
		}

		respondJSON(w, APIResponse{Success: true, Data: campaign})
	}
}

func getCampaignsHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		campaigns, err := services.DB.GetCampaigns()
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Failed to get campaigns"}, http.StatusInternalServerError)
			return
		}

		respondJSON(w, APIResponse{Success: true, Data: campaigns})
	}
}

func getCampaignHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid campaign ID"}, http.StatusBadRequest)
			return
		}

		campaign, err := services.DB.GetCampaign(id)
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Campaign not found"}, http.StatusNotFound)
			return
		}

		respondJSON(w, APIResponse{Success: true, Data: campaign})
	}
}

func testCampaignHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid campaign ID"}, http.StatusBadRequest)
			return
		}

		var req struct {
			TestEmails []string `json:"test_emails"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid request"}, http.StatusBadRequest)
			return
		}

		if len(req.TestEmails) == 0 {
			respondJSON(w, APIResponse{Success: false, Error: "No test emails provided"}, http.StatusBadRequest)
			return
		}

		campaign, err := services.DB.GetCampaign(id)
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Campaign not found"}, http.StatusNotFound)
			return
		}

		// Validate test emails
		for _, email := range req.TestEmails {
			if !isValidEmail(email) {
				respondJSON(w, APIResponse{Success: false, Error: fmt.Sprintf("Invalid email format: %s", email)}, http.StatusBadRequest)
				return
			}
		}

		// Send test emails
		results := make([]map[string]interface{}, 0, len(req.TestEmails))
		
		for _, email := range req.TestEmails {
			// Create test subscriber
			testSubscriber := &store.Subscriber{
				ID:     999999, // Special ID for test subscribers
				Email:  email,
				Status: "active",
			}

			// Create test message
			message := services.Mail.CreateCampaignMessage(campaign, testSubscriber)
			
			// Send email
			err := services.Mail.Send(message)
			
			result := map[string]interface{}{
				"email": email,
				"sent":  err == nil,
			}
			
			if err != nil {
				result["error"] = err.Error()
			}
			
			results = append(results, result)
		}

		respondJSON(w, APIResponse{
			Success: true,
			Data: map[string]interface{}{
				"message": "Test emails sent",
				"results": results,
			},
		})
	}
}

func scheduleCampaignHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid campaign ID"}, http.StatusBadRequest)
			return
		}

		var req struct {
			ScheduledAt *time.Time `json:"scheduled_at"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid request"}, http.StatusBadRequest)
			return
		}

		campaign, err := services.DB.GetCampaign(id)
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Campaign not found"}, http.StatusNotFound)
			return
		}

		// Validate campaign status
		if campaign.Status != "draft" {
			respondJSON(w, APIResponse{Success: false, Error: "Campaign must be in draft status to schedule"}, http.StatusBadRequest)
			return
		}

		// If no scheduled time provided, send immediately
		if req.ScheduledAt == nil {
			req.ScheduledAt = &time.Time{}
			*req.ScheduledAt = time.Now()
		}

		// Validate scheduled time
		if req.ScheduledAt.Before(time.Now()) {
			respondJSON(w, APIResponse{Success: false, Error: "Scheduled time cannot be in the past"}, http.StatusBadRequest)
			return
		}

		// Update campaign status
		if req.ScheduledAt.After(time.Now()) {
			// Schedule for later
			err = services.DB.UpdateCampaignStatus(id, "scheduled")
			if err != nil {
				respondJSON(w, APIResponse{Success: false, Error: "Failed to schedule campaign"}, http.StatusInternalServerError)
				return
			}

			// Enqueue campaign job
			payload := jobs.SendBatchPayload{
				CampaignID: id,
				Recipients: []string{}, // Will be populated by the job handler
			}
			
			err = services.Queue.Enqueue("send_batch", payload, *req.ScheduledAt)
			if err != nil {
				respondJSON(w, APIResponse{Success: false, Error: "Failed to enqueue campaign"}, http.StatusInternalServerError)
				return
			}

			respondJSON(w, APIResponse{
				Success: true,
				Data: map[string]interface{}{
					"message":      "Campaign scheduled successfully",
					"scheduled_at": req.ScheduledAt,
				},
			})
		} else {
			// Send immediately
			err = services.DB.UpdateCampaignStatus(id, "sending")
			if err != nil {
				respondJSON(w, APIResponse{Success: false, Error: "Failed to start campaign"}, http.StatusInternalServerError)
				return
			}

			// Enqueue campaign job for immediate execution
			payload := jobs.SendBatchPayload{
				CampaignID: id,
				Recipients: []string{}, // Will be populated by the job handler
			}
			
			err = services.Queue.Enqueue("send_batch", payload, time.Now())
			if err != nil {
				respondJSON(w, APIResponse{Success: false, Error: "Failed to enqueue campaign"}, http.StatusInternalServerError)
				return
			}

			respondJSON(w, APIResponse{
				Success: true,
				Data: map[string]interface{}{
					"message": "Campaign started successfully",
				},
			})
		}
	}
}

func getCampaignReportHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid campaign ID"}, http.StatusBadRequest)
			return
		}

		events, err := services.DB.GetCampaignEvents(id)
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Failed to get campaign events"}, http.StatusInternalServerError)
			return
		}

		respondJSON(w, APIResponse{Success: true, Data: events})
	}
}

func trackClickHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			CampaignID   int    `json:"campaign_id"`
			SubscriberID int    `json:"subscriber_id"`
			URL          string `json:"url"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid request"}, http.StatusBadRequest)
			return
		}

		meta, _ := json.Marshal(map[string]string{"url": req.URL})
		if err := services.DB.RecordEvent(req.CampaignID, req.SubscriberID, "click", meta); err != nil {
			logrus.Errorf("Failed to record click event: %v", err)
		}

		respondJSON(w, APIResponse{Success: true})
	}
}

func trackOpenHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			CampaignID   int `json:"campaign_id"`
			SubscriberID int `json:"subscriber_id"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid request"}, http.StatusBadRequest)
			return
		}

		if err := services.DB.RecordEvent(req.CampaignID, req.SubscriberID, "open", nil); err != nil {
			logrus.Errorf("Failed to record open event: %v", err)
		}

		respondJSON(w, APIResponse{Success: true})
	}
}

func unsubscribeHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		subscriberID, err := strconv.Atoi(vars["subscriberId"])
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid subscriber ID"}, http.StatusBadRequest)
			return
		}

		// TODO: Implement proper token validation
		_ = vars["token"]

		// Update subscriber status
		if err := services.DB.UpdateSubscriberStatus(subscriberID, "unsubscribed"); err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Failed to unsubscribe"}, http.StatusInternalServerError)
			return
		}

		// Add to suppressions
		subscriber, err := services.DB.GetSubscriber(subscriberID)
		if err == nil {
			services.DB.AddSuppression(subscriber.Email, "unsubscribed")
		}

		respondJSON(w, APIResponse{Success: true, Data: map[string]string{"message": "Successfully unsubscribed"}})
	}
}

func bounceHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bounceData struct {
			Email       string `json:"email"`
			Reason      string `json:"reason"`
			BounceType  string `json:"bounce_type"`
			CampaignID  int    `json:"campaign_id,omitempty"`
			SubscriberID int   `json:"subscriber_id,omitempty"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&bounceData); err != nil {
			logrus.Errorf("Failed to decode bounce data: %v", err)
			respondJSON(w, APIResponse{Success: false, Error: "Invalid bounce data"}, http.StatusBadRequest)
			return
		}

		// Validate required fields
		if bounceData.Email == "" {
			respondJSON(w, APIResponse{Success: false, Error: "Email is required"}, http.StatusBadRequest)
			return
		}

		// Determine bounce type if not provided
		if bounceData.BounceType == "" {
			if strings.Contains(strings.ToLower(bounceData.Reason), "complaint") {
				bounceData.BounceType = "complaint"
			} else if strings.Contains(strings.ToLower(bounceData.Reason), "hard") {
				bounceData.BounceType = "hard"
			} else {
				bounceData.BounceType = "soft"
			}
		}

		// Enqueue bounce processing job
		payload := jobs.BounceProcessingPayload{
			Email:      bounceData.Email,
			Reason:     bounceData.Reason,
			BounceType: bounceData.BounceType,
		}

		err := services.Queue.Enqueue("process_bounce", payload, time.Now())
		if err != nil {
			logrus.Errorf("Failed to enqueue bounce processing: %v", err)
			respondJSON(w, APIResponse{Success: false, Error: "Failed to process bounce"}, http.StatusInternalServerError)
			return
		}

		// Record bounce event if campaign/subscriber info is available
		if bounceData.CampaignID > 0 && bounceData.SubscriberID > 0 {
			meta, _ := json.Marshal(map[string]string{
				"reason":      bounceData.Reason,
				"bounce_type": bounceData.BounceType,
			})
			
			services.DB.RecordEvent(bounceData.CampaignID, bounceData.SubscriberID, "bounce", meta)
		}

		logrus.Infof("Bounce processed for %s: %s (%s)", bounceData.Email, bounceData.Reason, bounceData.BounceType)
		respondJSON(w, APIResponse{Success: true, Data: map[string]string{"message": "Bounce processed"}})
	}
}

func respondJSON(w http.ResponseWriter, response APIResponse, statusCode ...int) {
	w.Header().Set("Content-Type", "application/json")
	
	status := http.StatusOK
	if len(statusCode) > 0 {
		status = statusCode[0]
	}
	
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
