package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
		_, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Invalid list ID"}, http.StatusBadRequest)
			return
		}

		// TODO: Implement CSV import
		respondJSON(w, APIResponse{Success: true, Data: map[string]string{"message": "Import functionality coming soon"}})
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

		_, err = services.DB.GetCampaign(id)
		if err != nil {
			respondJSON(w, APIResponse{Success: false, Error: "Campaign not found"}, http.StatusNotFound)
			return
		}

		// TODO: Implement test send
		respondJSON(w, APIResponse{Success: true, Data: map[string]string{"message": "Test send functionality coming soon"}})
	}
}

func scheduleCampaignHandler(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		_, err := strconv.Atoi(vars["id"])
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

		// TODO: Implement campaign scheduling
		respondJSON(w, APIResponse{Success: true, Data: map[string]string{"message": "Schedule functionality coming soon"}})
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
		// TODO: Implement bounce handling
		respondJSON(w, APIResponse{Success: true})
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
