package billing

import (
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

// HandleStripeWebhook processes asynchronous HTTP receipts from Stripe mapping
// them locally to Postgres monetization subscriptions.
func HandleStripeWebhook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const MaxBodyBytes = int64(65536)
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Failed to read Stripe webhook body", "err", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
		signatureHeader := r.Header.Get("Stripe-Signature")

		event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
		if err != nil {
			slog.Warn("Invalid Stripe Signature Received", "err", err)
			w.WriteHeader(http.StatusBadRequest) // Return 400
			return
		}

		switch event.Type {
		case "customer.subscription.created", "customer.subscription.updated":
			var sub stripe.Subscription
			err := json.Unmarshal(event.Data.Raw, &sub)
			if err != nil {
				slog.Error("Failed parsing subscription event", "err", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			slog.Info("Stripe Subscription Created/Updated", "id", sub.ID, "status", sub.Status)
			
			// Map strict subscription bounds and calendar limits explicitly
			_, err = db.ExecContext(r.Context(), `
				UPDATE subscriptions
				SET status = $1, current_period_end = $2
				WHERE external_subscription_id = $3
			`, string(sub.Status), time.Unix(sub.CurrentPeriodEnd, 0), sub.ID)
			
			if err != nil {
				slog.Error("Failed syncing Stripe payload to database", "err", err)
			}
			
		case "customer.subscription.deleted":
			var sub stripe.Subscription
			_ = json.Unmarshal(event.Data.Raw, &sub)
			slog.Info("Stripe Subscription Deleted", "id", sub.ID)
			
			_, _ = db.ExecContext(r.Context(), `
				UPDATE subscriptions
				SET status = 'canceled'
				WHERE external_subscription_id = $1
			`, sub.ID)

		case "invoice.payment_failed":
			var invoice stripe.Invoice
			_ = json.Unmarshal(event.Data.Raw, &invoice)
			slog.Info("Stripe Invoice Payment Failed", "sub_id", invoice.Subscription.ID)
			
			// Instantly revoke strict active status preventing module access
			if invoice.Subscription != nil {
				_, _ = db.ExecContext(r.Context(), `
					UPDATE subscriptions
					SET status = 'past_due'
					WHERE external_subscription_id = $1
				`, invoice.Subscription.ID)
			}

		case "invoice.payment_succeeded":
			var invoice stripe.Invoice
			_ = json.Unmarshal(event.Data.Raw, &invoice)
			slog.Info("Stripe Invoice Payment Success", "sub_id", invoice.Subscription.ID)

		default:
			slog.Debug("Unhandled Stripe Event Type", "type", event.Type)
		}

		w.WriteHeader(http.StatusOK)
	}
}
